package kdbx

import (
	"os"
	"reflect"
	"testing"

	"github.com/tobischo/gokeepasslib/v3"
	"github.com/tobischo/gokeepasslib/v3/wrappers"
)

func count(item string, items []string) int {
	count := 0
	for _, v := range items {
		if item == v {
			count = count + 1
		}
	}

	return count
}

func makeEntry(title string) Entry {
	entry := gokeepasslib.NewEntry()
	entry.Values = append(entry.Values, gokeepasslib.ValueData{Key: "Title", Value: gokeepasslib.V{Content: title, Protected: wrappers.NewBoolWrapper(false)}})

	return entry
}

func makeGroup(name string, entries ...Entry) gokeepasslib.Group {
	group := gokeepasslib.NewGroup()
	group.Name = name

	group.Entries = append(group.Entries, entries...)

	return group
}

func makeDatabase(filename string, groups ...gokeepasslib.Group) *Database {
	gdb := gokeepasslib.NewDatabase()
	db := &Database{os.File{}, *gdb}
	gdb.Content.Root.Groups = groups

	return db
}

func TestDatabase_GetEntryPaths(t *testing.T) {
	topLevelEntry := makeEntry("TopLevelEntry")
	topLevelGroup := makeGroup("TopLevelGroup", topLevelEntry)

	nestedEntry := makeEntry("NestedEntry")
	entryWithNoName := makeEntry("")
	entryWithInvalidChars := makeEntry("Not/Split")
	entryWithDuplicateName := makeEntry("NestedEntry")
	nestedGroup := makeGroup("NestedGroup", nestedEntry, entryWithNoName, entryWithInvalidChars, entryWithDuplicateName)
	topLevelGroup.Groups = append(topLevelGroup.Groups, nestedGroup)

	db := makeDatabase("test.kdbx", topLevelGroup)
	entryPaths := db.GetEntryPaths()

	tests := []struct {
		name      string
		path      string
		wantCount int
	}{
		{"nested entries, with duplicates", "/TopLevelGroup/NestedGroup/NestedEntry", 2},
		{"top level entries", "/TopLevelGroup/TopLevelEntry", 1},
		{"entries without title", "/TopLevelGroup/NestedGroup/(UNKNOWN)", 1},
		{"entries with invalid characters", "/TopLevelGroup/NestedGroup/NotSplit", 1},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotCount := count(tt.path, entryPaths)
			if gotCount != tt.wantCount {
				t.Errorf("Database.GetEntryPaths(). Unexected '%v', wantedCount %v got %v.", tt.path, tt.wantCount, gotCount)
			}
		})
	}
}

func TestDatabase_GetFirstEntryByPath(t *testing.T) {
	entry1 := makeEntry("e")
	entry2 := makeEntry("e")

	tests := []struct {
		name      string
		path      string
		wantEntry *Entry
		db        *Database
	}{
		{"finds the first entry", "/g/e", &entry1, makeDatabase("d", makeGroup("g", entry1))},
		{"finds the second entry", "/g/e", &entry2, makeDatabase("d", makeGroup("g", entry2, entry1))},
		{"does not find the entry", "/no/no", nil, makeDatabase("d", makeGroup("g", entry1))},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotEntry := tt.db.GetFirstEntryByPath(tt.path)
			if !reflect.DeepEqual(tt.wantEntry, gotEntry) {
				t.Errorf("Database.GetFirstEntryByPath() gotEntry = %v, want %v", gotEntry, tt.wantEntry)
			}
		})
	}
}

func TestDatabase_GetEntry(t *testing.T) {
	entry1 := makeEntry("e")
	entry2 := makeEntry("e")

	tests := []struct {
		name      string
		uuid      gokeepasslib.UUID
		wantEntry *Entry
		db        *Database
	}{
		{"finds the first entry", entry1.UUID, &entry1, makeDatabase("d", makeGroup("g", entry1))},
		{"finds the second entry", entry2.UUID, &entry2, makeDatabase("d", makeGroup("g", entry2, entry1))},
		{"does not find the entry", gokeepasslib.NewUUID(), nil, makeDatabase("d", makeGroup("g", entry1))},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotEntry := tt.db.GetEntry(tt.uuid)
			if !reflect.DeepEqual(tt.wantEntry, gotEntry) {
				t.Errorf("Database.GetEntry() gotEntry = %v, want %v", gotEntry, tt.wantEntry)
			}
		})
	}
}

func TestDatabase_Save(t *testing.T) {
	files := []*os.File{}

	makeFile := func() *os.File {
		file, _ := os.CreateTemp(os.TempDir(), "test*.kdbx")
		files = append(files, file)
		return file
	}

	defer func() {
		for _, file := range files {
			os.Remove(file.Name())
		}
	}()

	type fields struct {
		file     os.File
		Database gokeepasslib.Database
	}
	tests := []struct {
		name    string
		fields  fields
		wantErr bool
	}{
		{"does not save with missing file", fields{os.File{}, *gokeepasslib.NewDatabase()}, true},
		{"saves the database to file", fields{*makeFile(), *gokeepasslib.NewDatabase()}, false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			d := &Database{
				file:     tt.fields.file,
				Database: tt.fields.Database,
			}
			if err := d.Save(); (err != nil) != tt.wantErr {
				t.Errorf("Database.Save() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
