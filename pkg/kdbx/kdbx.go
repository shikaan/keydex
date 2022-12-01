package kdbx

import (
	"os"
	"strings"

	"github.com/shikaan/kpcli/pkg/errors"
	"github.com/tobischo/gokeepasslib/v3"
)

type Database struct {
	file os.File
  unlocked bool

	gokeepasslib.Database
}

type Entry = gokeepasslib.Entry
type EntryField = gokeepasslib.ValueData
type UUID = gokeepasslib.UUID

// A string like "/Database/Group/EntryName"
type EntryPath = string

const PATH_SEPARATOR = "/"
const TITLE_KEY = "Title"

func New(filepath string) (*Database, error) {
	file, err := os.Open(filepath)

	if err != nil {
		return nil, errors.MakeError(err.Error(), "kdbx")
	}

	db := gokeepasslib.NewDatabase()
	return &Database{*file, false, *db}, nil
}

func NewUnlocked(filepath, password string) (*Database, error) {
	kdbx, err := New(filepath)

	if err != nil {
		return nil, err
	}

	if err := kdbx.UnlockWithPassword(password); err != nil {
		return nil, err
	}

	return kdbx, nil
}

func (d *Database) UnlockWithPassword(password string) error {
	d.Credentials = gokeepasslib.NewPasswordCredentials(password)
  
  return d.unlock()
}

func (d *Database) GetEntryPaths() []EntryPath {
	result := []EntryPath{}

	for _, uEP := range d.getEntryPathsAndUUIDs() {
		result = append(result, uEP.path)
	}

	return result
}

// Returns the first entry matching the entry path provided.
// Note that the path might not be unique. Use the UUID method 
// when identifying a precise entry is necessary 
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

func (d *Database) LockProtectedEntries() error {
  if !d.unlocked {
    return errors.MakeError("Cannot lock a locked database", "kdbx")
  }

  if e := d.Database.LockProtectedEntries(); e != nil {
    return errors.MakeError(e.Error(), "kdbx")
  }
  d.unlocked = false
  return nil
}

func (d *Database) UnlockProtectedEntries() error {
  if d.unlocked {
    return errors.MakeError("Cannot unlock an unlocked database", "kdbx")
  }

  if e := d.Database.UnlockProtectedEntries(); e != nil {
    return errors.MakeError(e.Error(), "kdbx")
  }
  d.unlocked = true
  return nil
}

func (d *Database) Save() error {	
  if err := d.LockProtectedEntries(); err != nil {
		return errors.MakeError(err.Error(), "kdbx")
	}
	
  d.file.Close()
	file, err := os.Create(d.file.Name())
  if err != nil {
    return errors.MakeError(err.Error(), "kdbx")
  }

	if err := gokeepasslib.NewEncoder(file).Encode(&d.Database); err != nil {
		return errors.MakeError(err.Error(), "kdbx")
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

// Decodes and unlocks a database whose credentials are known.
// Use UnlockWith* methods to store credentials
func (d *Database) unlock() error {
  // TODO: it's probably possible to unlock without password,
  // since you can create credentials-less archives
	if d.Credentials == nil {
    return errors.MakeError("Cannot unlock without credentials", "kdbx")
  }

  if err := gokeepasslib.NewDecoder(&d.file).Decode(&d.Database); err != nil {
		return errors.MakeError(err.Error(), "kdbx")
	}

	d.UnlockProtectedEntries()
	return nil
}


func getEntryPathsFromGroup(g gokeepasslib.Group, prefix string) []uniqueEntryPath {
	groupPrefix := prefix + g.Name + PATH_SEPARATOR
	entries := []uniqueEntryPath{}

	for _, subGroup := range g.Groups {
		subEntries := getEntryPathsFromGroup(subGroup, groupPrefix)
		entries = append(entries, subEntries...)
	}

	for _, entry := range g.Entries {
    title := entry.GetTitle()

    if title == "" {
      title = "(UNKNOWN)"
    }
    
    key := groupPrefix + sanitizePathPortion(title)
    entries = append(entries, uniqueEntryPath{path: key, uuid: entry.UUID})
	}

	return entries
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
