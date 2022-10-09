package kdbx

import (
	"os"

	"github.com/shikaan/kpcli/pkg/errors"
	"github.com/tobischo/gokeepasslib/v3"
)

type Database struct {
	db   *gokeepasslib.Database
	file *os.File

	Name        string
	Description string

	Groups []Group
}

type Group struct {
	Name        string
	Description string

	Groups  []Group
	Entries []Entry
}

type Entry = gokeepasslib.Entry

type Credentials = string // This will need to support files and so forth

func New(filepath string) (*Database, error) {
	file, err := os.Open(filepath)

	if err != nil {
		return nil, errors.MakeError(err.Error(), "kdbx")
	}

	db := gokeepasslib.NewDatabase()

	return &Database{db: db, file: file}, nil
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

	kdbx.Groups = parseGroups(db.Content.Root.Groups)
}

func parseGroups(root []gokeepasslib.Group) []Group {
	var result []Group

	for _, g := range root {
		group := Group{Name: g.Name, Entries: g.Entries, Groups: parseGroups(g.Groups)}
		result = append(result, group)
	}

	return result
}
