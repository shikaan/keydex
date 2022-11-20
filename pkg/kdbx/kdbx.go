package kdbx

import (
	"os"
	"strings"

	"github.com/shikaan/kpcli/pkg/errors"
	"github.com/tobischo/gokeepasslib/v3"
)

type Database struct {
	file os.File

	gokeepasslib.Database
}

type Entry = gokeepasslib.Entry
type EntryField = gokeepasslib.ValueData
type UUID = gokeepasslib.UUID

// A string like "/Database/Group/EntryName"
type EntryPath = string

const PATH_SEPARATOR = "/"

func New(filepath string) (*Database, error) {
	file, err := os.Open(filepath)

	if err != nil {
		return nil, errors.MakeError(err.Error(), "kdbx")
	}

	db := gokeepasslib.NewDatabase()
	return &Database{*file, *db}, nil
}

func NewUnlocked(filepath, password string) (*Database, error) {
	kdbx, err := New(filepath)

	if err != nil {
		return nil, err
	}

	if err := kdbx.Unlock(password); err != nil {
		return nil, err
	}

	return kdbx, nil
}

func (d *Database) Unlock(password string) error {
	d.Credentials = gokeepasslib.NewPasswordCredentials(password)

	if err := gokeepasslib.NewDecoder(&d.file).Decode(&d.Database); err != nil {
		return errors.MakeError(err.Error(), "kdbx")
	}

	d.UnlockProtectedEntries()
	return nil
}

func (d *Database) GetEntryPaths() []EntryPath {
	result := []EntryPath{}

	for _, uEP := range d.getEntryPathsAndUUIDs() {
		result = append(result, uEP.path)
	}

	return result
}

// Returns the first entry matching the entry path provided.
// Please note: the path might not be unique! Use the UUID method
func (d *Database) GetFirstEntryByPath(p EntryPath) *Entry {
	for _, uEP := range d.getEntryPathsAndUUIDs() {
		if uEP.path == p {
			return d.GetEntry(uEP.uuid)
		}
	}

	return nil
}

func (d *Database) GetEntry(uuid gokeepasslib.UUID) *Entry {
	for _, g := range d.Content.Root.Groups {
		if e := getEntryByUUID(g, uuid); e != nil {
			return e
		}
	}

	return nil
}

func (d *Database) Save() error {
	if err := d.LockProtectedEntries(); err != nil {
		return err
	}

	d.file.Close()
	file, _ := os.Create(d.file.Name())

	if err := gokeepasslib.NewEncoder(file).Encode(&d.Database); err != nil {
		return err
	}

	d.file = *file

	return nil
}

// Private

type uniqueEntryPath struct { 
  path EntryPath 
  uuid UUID 
}

func (d *Database) getEntryPathsAndUUIDs() []uniqueEntryPath {
	result := []uniqueEntryPath{}

	for _, g := range d.Content.Root.Groups {
		result = append(result, getEntryPathsFromGroup(g, PATH_SEPARATOR)...)
	}

	return result
}

func getEntryPathsFromGroup(g gokeepasslib.Group, prefix string) []uniqueEntryPath {
	groupPrefix := prefix + g.Name + PATH_SEPARATOR
	entries := []uniqueEntryPath{}

	for _, subGroup := range g.Groups {
		subEntries := getEntryPathsFromGroup(subGroup, groupPrefix)
		entries = append(entries, subEntries...)
	}

	for _, entry := range g.Entries {
		key := groupPrefix + sanitizePathPortion(entry.GetTitle())
    entries = append(entries, uniqueEntryPath{path: key, uuid: entry.UUID})
	}

	return entries
}

func getEntryFromGroup(g gokeepasslib.Group, entryPathPortions []string) *Entry {
	isLeaf := len(entryPathPortions) == 1
	current := entryPathPortions[0]

  println("searching in", g.Name)

	if isLeaf {
		for _, e := range g.Entries {
			if e.GetTitle() == current {
        println("found in", g.Name)
				return &e
			}
		}

		return nil
	}

	for _, gs := range g.Groups {
		if gs.Name == current {
			return getEntryFromGroup(gs, entryPathPortions[1:])
		}
	}

	return nil
}

func getEntryByUUID(g gokeepasslib.Group, uuid gokeepasslib.UUID) *Entry {
	for _, e := range g.Entries {
		if e.UUID.Compare(uuid) {
			return &e
		}
	}
	
  for _, gs := range g.Groups {
		if e := getEntryByUUID(gs, uuid); e != nil {
			return e
		}
	}

	return nil
}

func sanitizePathPortion(s string) string {
	return strings.ReplaceAll(s, PATH_SEPARATOR, "")
}
