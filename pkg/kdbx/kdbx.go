package kdbx

import (
	"os"

	"github.com/shikaan/kpcli/pkg/errors"
	"github.com/shikaan/kpcli/pkg/utils"
	"github.com/tobischo/gokeepasslib/v3"
)

type Database struct {
	db   *gokeepasslib.Database
	file *os.File

	Name        string
	Description string

	Groups []Group
  Entries map[EntryPath]Entry
}

type Group struct {
	Name        string
	Description string

	Groups  []Group
	Entries []Entry
}

type Entry = gokeepasslib.Entry

// A string like "/Database/Group/EntryName"
type EntryPath = string 

type Credentials = string // This will need to support files and so forth

func New(filepath string) (*Database, error) {
	file, err := os.Open(filepath)

	if err != nil {
		return nil, errors.MakeError(err.Error(), "kdbx")
	}

	db := gokeepasslib.NewDatabase()

	return &Database{db: db, file: file}, nil
}

func NewUnlocked(filepath string, password Credentials) (*Database, error) {
  kdbx, err := New(filepath)
	if err != nil {
		return nil, err
	}

	err = kdbx.Unlock(password)
	if err != nil {
		return nil, err
	}

  return kdbx, nil
}

func (kdbx *Database) Unlock(password Credentials) error {
	kdbx.db.Credentials = gokeepasslib.NewPasswordCredentials(password)

	err := gokeepasslib.NewDecoder(kdbx.file).Decode(kdbx.db)

	if err != nil {
		return errors.MakeError(err.Error(), "kdbx")
	}

  kdbx.db.UnlockProtectedEntries()

	kdbx.parseUnlockedDatabase(*kdbx.db)
	return nil
}

func (kdbx *Database) parseUnlockedDatabase(db gokeepasslib.Database) {
	kdbx.Name = db.Content.Meta.DatabaseName
	kdbx.Description = db.Content.Meta.DatabaseDescription

	kdbx.Groups, kdbx.Entries = parseGroups(db.Content.Root.Groups, "")
}

func parseGroups(root []gokeepasslib.Group, prefix string) ([]Group, map[EntryPath]Entry) {
  groups := []Group{}
  entries := map[EntryPath]Entry{}

	for _, g := range root {
    groupPrefix := prefix + "/" + g.Name
    subGroups, subEntries := parseGroups(g.Groups, groupPrefix)
    
    for _, e := range g.Entries {
      subEntries[prefix + "/" + e.GetTitle()] = e
    }

    group := Group{Name: g.Name, Entries: g.Entries, Groups: subGroups}	
    groups = append(groups, group) 
    
    entries = utils.Merge(entries, subEntries)
	}

	return groups, entries
}

