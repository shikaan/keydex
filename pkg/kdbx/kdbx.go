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

// A string like "/Database/Group/EntryName"
type EntryPath = string
type Entries = map[EntryPath]*gokeepasslib.Entry

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

	err = kdbx.Unlock(password)
	if err != nil {
		return nil, err
	}

	return kdbx, nil
}

func (d *Database) Unlock(password string) error {
	d.Credentials = gokeepasslib.NewPasswordCredentials(password)

	err := gokeepasslib.NewDecoder(&d.file).Decode(&d.Database)

	if err != nil {
		return errors.MakeError(err.Error(), "kdbx")
	}

	d.UnlockProtectedEntries()
	return nil
}

func (d *Database) GetEntryPaths() []EntryPath {
	result := []EntryPath{}

	for _, g := range d.Content.Root.Groups {
		result = append(result, getEntryPathsFromGroup(g, "/")...)
	}

	return result
}

func (d *Database) GetEntry(p EntryPath) *Entry {
	// Skip the first /
	portions := strings.Split(p[1:], "/")

	for _, g := range d.Content.Root.Groups {
		if e := getEntryFromGroup(g, portions); e != nil {
			return e
		}
	}

	return nil
}

func (d *Database) SetEntry(p EntryPath, e Entry) error {
  // Skip the first /
	portions := strings.Split(p[1:], "/")

	for _, g := range d.Content.Root.Groups {
    if err := setEntryFromGroup(g, portions, e); err != nil {
      return err
    }
	}

  return nil
}

func (d *Database) Save() error {
	if err := d.LockProtectedEntries(); err != nil {
		return err
	}

	file, _ := os.Create(d.file.Name())
	encoder := gokeepasslib.NewEncoder(file)

	if err := encoder.Encode(&d.Database); err != nil {
		return err
	}

	return nil
}

// Private

func getEntryPathsFromGroup(g gokeepasslib.Group, prefix string) []EntryPath {
	groupPrefix := prefix + g.Name + "/"
	entries := []EntryPath{}

	for _, subGroup := range g.Groups {
		subEntries := getEntryPathsFromGroup(subGroup, groupPrefix)
		entries = append(entries, subEntries...)
	}

	for _, entry := range g.Entries {
		key := prefix + entry.GetTitle()
		entries = append(entries, key)
	}

	return entries
}

func getEntryFromGroup(g gokeepasslib.Group, entryPathPortions []string) *Entry {
	isLeaf := len(entryPathPortions) == 1
	current := entryPathPortions[0]

	if isLeaf {
		for _, e := range g.Entries {
			if e.GetTitle() == current {
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

func setEntryFromGroup(g gokeepasslib.Group, entryPathPortions []string, newEntry Entry) error {
	isLeaf := len(entryPathPortions) == 1
	current := entryPathPortions[0]

  if isLeaf {
		for i, e := range g.Entries {
			if e.GetTitle() == current {
        g.Entries[i] = newEntry
				return nil
			}
		}

		return errors.MakeError("Could not find the entry", "kdbx")
	}

	for _, gs := range g.Groups {
		if gs.Name == current {
			return setEntryFromGroup(gs, entryPathPortions[1:], newEntry)
		}
	}

	return errors.MakeError("Could not find the entry", "kdbx")
}
