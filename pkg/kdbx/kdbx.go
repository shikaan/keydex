package kdbx

import (
	"crypto/rand"
	"encoding/base64"
	"os"
	"slices"
	"strings"
	"time"

	"github.com/shikaan/keydex/pkg/errors"
	"github.com/tobischo/gokeepasslib/v3"
	"github.com/tobischo/gokeepasslib/v3/wrappers"
)

type Database struct {
	file os.File

	gokeepasslib.Database
}

type Entry struct{ gokeepasslib.Entry }
type Group = gokeepasslib.Group
type EntryField = gokeepasslib.ValueData
type UUID = gokeepasslib.UUID

// A string like "/Database/Group/EntryName"
type EntityPath = string

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

func (d *Database) GetEntryPaths() []EntityPath {
	result := []EntityPath{}

	for _, uEP := range d.getEntryPathsAndUUIDs() {
		result = append(result, uEP.path)
	}

	return result
}

func (d *Database) GetGroupPaths() []EntityPath {
	result := []EntityPath{}

	for _, uEP := range d.getGroupPathsAndUUIDs() {
		result = append(result, uEP.path)
	}

	return result
}

// Returns the first entry matching the entry path provided.
// Note that the path might not be unique. Use the UUID method
// when identifying a precise entry is necessary
func (d *Database) GetFirstEntryByPath(p EntityPath) *Entry {
	for _, uEP := range d.getEntryPathsAndUUIDs() {
		if uEP.path == p {
			return d.GetEntry(uEP.uuid)
		}
	}

	return nil
}

// Returns the first group matching the entry path provided.
// Note that the path might not be unique. Use the UUID method
// when identifying a precise group is necessary
func (d *Database) GetFirstGroupByPath(p EntityPath) *Group {
	for _, uEP := range d.getGroupPathsAndUUIDs() {
		if uEP.path == p {
			return d.GetGroup(uEP.uuid)
		}
	}

	return nil
}

func (d *Database) GetGroup(uuid gokeepasslib.UUID) *Group {
	for i := range d.Content.Root.Groups {
		if d.Content.Root.Groups[i].UUID.Compare(uuid) {
			return &d.Content.Root.Groups[i]
		}

		if g, _ := getNestedGroupByUUID(&d.Content.Root.Groups[i], uuid); g != nil {
			return g
		}
	}

	return nil
}

// Returns an entry by its UUID
func (d *Database) GetEntry(uuid gokeepasslib.UUID) *Entry {
	for i := range d.Content.Root.Groups {
		if e, _ := getEntryByUUID(&d.Content.Root.Groups[i], uuid); e != nil {
			return e
		}
	}

	return nil
}

func (d *Database) RemoveEntry(uuid gokeepasslib.UUID) error {
	for i := range d.Content.Root.Groups {
		if entry, subGroup := getEntryByUUID(&d.Content.Root.Groups[i], uuid); subGroup != nil {
			subGroup.Entries = slices.DeleteFunc(subGroup.Entries, func(e gokeepasslib.Entry) bool {
				return e.UUID.Compare(entry.UUID)
			})
			return nil
		}
	}

	return errors.MakeError("entry not found", "kdbx")
}

func (d *Database) RemoveGroup(uuid gokeepasslib.UUID) error {
	for i := range d.Content.Root.Groups {
		if group, parent := getNestedGroupByUUID(&d.Content.Root.Groups[i], uuid); group != nil && parent != nil {
			parent.Groups = slices.DeleteFunc(parent.Groups,
				func(g gokeepasslib.Group) bool {
					return g.UUID.Compare(group.UUID)
				})
			return nil
		}
	}

	return errors.MakeError("group not found", "kdbx")
}

func (d *Database) MoveEntryToGroup(entry *Entry, group *Group) {
	entryGroup := d.GetGroupForEntry(entry)

	// Group is nil when this is a new entry, no need to move
	if entryGroup == nil {
		group.Entries = append(group.Entries, entry.Entry)
		return
	}

	// If source and destination are the same, do nothing
	if entryGroup.UUID.Compare(group.UUID) {
		return
	}

	group.Entries = append(group.Entries, entry.Entry)
	entryGroup.Entries = slices.DeleteFunc(entryGroup.Entries, func(e gokeepasslib.Entry) bool {
		return e.UUID.Compare(entry.UUID)
	})
}

// Builds the full path for an entry within the specified group.
// Returns an error if the group is not found in the database.
func (d *Database) MakeEntryEntityPath(entry *Entry, group *Group) (EntityPath, error) {
	for _, path := range d.getGroupPathsAndUUIDs() {
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
	return &Entry{entry}
}

func (d *Database) NewGroup(name string) *Group {
	group := gokeepasslib.NewGroup()
	group.Name = name
	group.Entries = make([]gokeepasslib.Entry, 0)
	group.Groups = make([]gokeepasslib.Group, 0)

	return &group
}

func (d *Database) getGroupForEntry(entry *Entry, group *Group) *Group {
	for _, e := range group.Entries {
		if e.UUID.Compare(entry.UUID) {
			return group
		}
	}

	for i := range group.Groups {
		if g := d.getGroupForEntry(entry, &group.Groups[i]); g != nil {
			return g
		}
	}

	return nil
}

func (d *Database) GetGroupForEntry(entry *Entry) *Group {
	for i := range d.Content.Root.Groups {
		if g := d.getGroupForEntry(entry, &d.Content.Root.Groups[i]); g != nil {
			return g
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

func (d *Database) SaveAndUnlockEntries() error {
	err := d.Save()
	if err != nil {
		return err
	}

	err = d.UnlockProtectedEntries()
	if err != nil {
		return err
	}

	return nil
}

type uniqueEntityPath struct {
	path EntityPath
	uuid UUID
}

func (d *Database) getEntryPathsAndUUIDs() []uniqueEntityPath {
	result := []uniqueEntityPath{}

	for _, g := range d.Content.Root.Groups {
		result = append(result, getEntityPathsFromGroup(g, PATH_SEPARATOR)...)
	}

	return result
}

func (d *Database) getGroupPathsAndUUIDs() []uniqueEntityPath {
	result := []uniqueEntityPath{}

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

func getEntityPathsFromGroup(g Group, prefix string) []uniqueEntityPath {
	groupPrefix := prefix + g.Name + PATH_SEPARATOR
	entries := []uniqueEntityPath{}

	for _, subGroup := range g.Groups {
		subEntries := getEntityPathsFromGroup(subGroup, groupPrefix)
		entries = append(entries, subEntries...)
	}

	for _, entry := range g.Entries {
		title := entry.GetTitle()

		if title == "" {
			title = "(UNKNOWN)"
		}

		key := groupPrefix + sanitizePathPortion(title)
		entries = append(entries, uniqueEntityPath{path: key, uuid: entry.UUID})
	}

	return entries
}

func getGroupPathsFromGroup(g Group, prefix string) []uniqueEntityPath {
	groupPrefix := prefix + g.Name + PATH_SEPARATOR
	paths := []uniqueEntityPath{}

	paths = append(paths, uniqueEntityPath{path: groupPrefix, uuid: g.UUID})

	for _, subGroup := range g.Groups {
		subPaths := getGroupPathsFromGroup(subGroup, groupPrefix)
		paths = append(paths, subPaths...)
	}

	return paths
}

func getEntryByUUID(g *Group, uuid gokeepasslib.UUID) (*Entry, *Group) {
	for i := range g.Entries {
		if g.Entries[i].UUID.Compare(uuid) {
			return &Entry{g.Entries[i]}, g
		}
	}

	for i := range g.Groups {
		if e, foundGroup := getEntryByUUID(&g.Groups[i], uuid); e != nil {
			return e, foundGroup
		}
	}

	return nil, nil
}

func getNestedGroupByUUID(root *Group, uuid gokeepasslib.UUID) (*Group, *Group) {
	for i := range root.Groups {
		if root.Groups[i].UUID.Compare(uuid) {
			return &root.Groups[i], root
		}

		if group, parent := getNestedGroupByUUID(&root.Groups[i], uuid); group != nil {
			return group, parent
		}
	}

	return nil, nil
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

func (e *Entry) SetValue(key string, value string) {
	v := e.Get(key)
	v.Value.Content = value
}

func (e *Entry) SetLastUpdated() {
	e.Times.LastModificationTime = &wrappers.TimeWrapper{Time: time.Now()}
}
