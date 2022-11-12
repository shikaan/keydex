package kdbx

import (
	"os"

	"github.com/shikaan/kpcli/pkg/errors"
	"github.com/shikaan/kpcli/pkg/utils"
	"github.com/tobischo/gokeepasslib/v3"
)

type Database struct {
	file os.File

	gokeepasslib.Database
}

// Reexporting the types

type Entry = gokeepasslib.Entry
type EntryValue = gokeepasslib.ValueData

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

func (kdbx *Database) Unlock(password string) error {
	kdbx.Credentials = gokeepasslib.NewPasswordCredentials(password)

	err := gokeepasslib.NewDecoder(&kdbx.file).Decode(&kdbx.Database)

	if err != nil {
		return errors.MakeError(err.Error(), "kdbx")
	}

	kdbx.UnlockProtectedEntries()
	return nil
}

func (kdbx *Database) GetEntries() Entries {
  result := make(Entries, 1)

  for _, g := range kdbx.Content.Root.Groups {
    result = utils.Merge(result, getEntriesFromGroup(g, "/"))
  }

  return result
}

func getEntriesFromGroup(g gokeepasslib.Group, prefix string) Entries {
	groupPrefix := prefix + g.Name + "/"
	entries := make(Entries, 1)

	for _, subGroup := range g.Groups {
		subEntries := getEntriesFromGroup(subGroup, groupPrefix)
		entries = utils.Merge(entries, subEntries)
	}

	for _, entry := range g.Entries {
		key := prefix + entry.GetTitle()

		entries[key] = &entry
	}
	
  return entries
}
