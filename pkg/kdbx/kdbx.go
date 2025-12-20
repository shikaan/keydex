package kdbx

import (
	"crypto/rand"
	"encoding/base64"
	"os"
	"strings"

	"github.com/shikaan/keydex/pkg/errors"
	"github.com/tobischo/gokeepasslib/v3"
	"github.com/tobischo/gokeepasslib/v3/wrappers"
)

type Database struct {
	file os.File

	gokeepasslib.Database
}

type Entry = gokeepasslib.Entry
type Group = gokeepasslib.Group
type EntryField = gokeepasslib.ValueData
type UUID = gokeepasslib.UUID

// A string like "/Database/Group/EntryName"
type EntryPath = string

const PATH_SEPARATOR = "/"

const TITLE_KEY = "Title"
const PASSWORD_KEY = "Password"
const USERNAME_KEY = "UserName"

func New(filepath, password, keypath string) (*Database, error) {
	file, err := os.Open(filepath)

	if err != nil {
		return nil, errors.MakeError(err.Error(), "kdbx")
	}

	kdbx := &Database{*file, *gokeepasslib.NewDatabase()}

	if err := kdbx.UnlockWithPasswordAndKey(password, keypath); err != nil {
		return nil, err
	}

	return kdbx, nil
}

func (d *Database) UnlockWithPasswordAndKey(password, keypath string) error {
	if keypath == "" {
		d.Credentials = gokeepasslib.NewPasswordCredentials(password)
	} else {
		credentials, err := gokeepasslib.NewPasswordAndKeyCredentials(password, keypath)

		if err != nil {
			return errors.MakeError(err.Error(), "kdbx")
		}

		d.Credentials = credentials
	}

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

// Returns an entry by its UUID
func (d *Database) GetEntry(uuid gokeepasslib.UUID) *Entry {
	for _, g := range d.Content.Root.Groups {
		if e := getEntryByUUID(g, uuid); e != nil {
			return e
		}
	}

	return nil
}

// Builds the full path for an entry within the specified group.
// Returns an error if the group is not found in the database.
func (d *Database) MakeEntryPath(entry *Entry, group *Group) (EntryPath, error) {
	for _, path := range d.getGroupPaths() {
		if path.uuid.Compare(group.UUID) {
			return path.path + entry.GetTitle(), nil
		}
	}

	return "", errors.MakeError("cannot find group "+group.Name, "kdbx")
}

// Return a new Entry with default title, user, and a random password set
func (d *Database) NewEntry() *Entry {
	entry := gokeepasslib.NewEntry()
	entry.Values = append(entry.Values, gokeepasslib.ValueData{
		Key:   TITLE_KEY,
		Value: gokeepasslib.V{Content: "New"},
	})
	entry.Values = append(entry.Values, gokeepasslib.ValueData{
		Key:   USERNAME_KEY,
		Value: gokeepasslib.V{Content: "user"},
	})
	entry.Values = append(entry.Values, gokeepasslib.ValueData{
		Key: PASSWORD_KEY,
		Value: gokeepasslib.V{
			Content:   generateRandomString(16),
			Protected: wrappers.NewBoolWrapper(true),
		},
	})
	return &entry
}

func (d *Database) getGroupForEntry(entry *Entry, group *Group) *Group {
	for _, e := range group.Entries {
		if e.UUID.Compare(entry.UUID) {
			return group
		}
	}

	for _, g := range group.Groups {
		if result := d.getGroupForEntry(entry, &g); result != nil {
			return result
		}
	}

	return nil
}

func (d *Database) GetGroupForEntry(entry *Entry) *Group {
	for _, g := range d.Content.Root.Groups {
		if result := d.getGroupForEntry(entry, &g); result != nil {
			return result
		}
	}
	return nil
}

func (d *Database) GetRootGroup() *Group {
	if len(d.Content.Root.Groups) > 0 {
		return &d.Content.Root.Groups[0]
	}
	return nil
}

func (d *Database) Save() error {
	if err := d.Database.LockProtectedEntries(); err != nil {
		return errors.MakeError(err.Error(), "kdbx")
	}

	if err := d.file.Close(); err != nil {
		return errors.MakeError(err.Error(), "kdbx")
	}

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

func (d *Database) getGroupPaths() []uniqueEntryPath {
	result := []uniqueEntryPath{}

	for _, g := range d.Content.Root.Groups {
		result = append(result, getGroupPathsFromGroup(g, PATH_SEPARATOR)...)
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

	d.Database.UnlockProtectedEntries()
	return nil
}

func getEntryPathsFromGroup(g Group, prefix string) []uniqueEntryPath {
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

func getGroupPathsFromGroup(g Group, prefix string) []uniqueEntryPath {
	groupPrefix := prefix + g.Name + PATH_SEPARATOR
	paths := []uniqueEntryPath{}

	paths = append(paths, uniqueEntryPath{path: groupPrefix, uuid: g.UUID})

	for _, subGroup := range g.Groups {
		subPaths := getGroupPathsFromGroup(subGroup, groupPrefix)
		paths = append(paths, subPaths...)
	}

	return paths
}

func getEntryByUUID(g Group, uuid gokeepasslib.UUID) *Entry {
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

func generateRandomString(length uint) string {
	b := make([]byte, length)
	_, err := rand.Read(b)
	if err != nil {
		return "change-me"
	}
	return base64.RawStdEncoding.EncodeToString(b)
}
