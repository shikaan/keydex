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

func makeDatabase(_ string, groups ...gokeepasslib.Group) *Database {
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

func TestDatabase_NewEntry(t *testing.T) {
	db := makeDatabase("test.kdbx")

	entry := db.NewEntry()

	if entry == nil {
		t.Fatal("Database.NewEntry() returned nil")
	}

	// Check that title field exists and has expected value
	title := entry.GetTitle()
	if title != "New" {
		t.Errorf("Database.NewEntry() title = %v, want %v", title, "New")
	}

	// Check that username field exists and has expected value
	username := entry.Get(USERNAME_KEY)
	if username == nil {
		t.Fatal("Database.NewEntry() username field not found")
	}
	if username.Value.Content != "user" {
		t.Errorf("Database.NewEntry() username = %v, want %v", username.Value.Content, "user")
	}

	// Check that password field exists
	password := entry.Get(PASSWORD_KEY)
	if password == nil {
		t.Fatal("Database.NewEntry() password field not found")
	}

	// Check that password is protected
	if !password.Value.Protected.Bool {
		t.Error("Database.NewEntry() password is not protected")
	}

	// Check that password is not empty
	if password.Value.Content == "" {
		t.Error("Database.NewEntry() password is empty")
	}

	// Check that password is not the fallback value
	if password.Value.Content == "change-me" {
		t.Error("Database.NewEntry() password generation failed, got fallback value")
	}
}

func TestDatabase_GetGroupForEntry(t *testing.T) {
	entry1 := makeEntry("entry1")
	entry2 := makeEntry("entry2")
	entry3 := makeEntry("entry3")

	group1 := makeGroup("Group1", entry1)
	group2 := makeGroup("Group2", entry2)
	nestedGroup := makeGroup("NestedGroup", entry3)
	group1.Groups = append(group1.Groups, nestedGroup)

	db := makeDatabase("test.kdbx", group1, group2)

	tests := []struct {
		name      string
		entry     *Entry
		wantGroup string
	}{
		{"finds group for entry in top level group", &entry1, "Group1"},
		{"finds group for entry in different top level group", &entry2, "Group2"},
		{"finds group for entry in nested group", &entry3, "NestedGroup"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotGroup := db.GetGroupForEntry(tt.entry)
			if gotGroup == nil {
				t.Fatalf("Database.GetGroupForEntry() returned nil")
			}
			if gotGroup.Name != tt.wantGroup {
				t.Errorf("Database.GetGroupForEntry() group name = %v, want %v", gotGroup.Name, tt.wantGroup)
			}
		})
	}

	t.Run("returns nil for non-existent entry", func(t *testing.T) {
		nonExistentEntry := makeEntry("nonexistent")
		gotGroup := db.GetGroupForEntry(&nonExistentEntry)
		if gotGroup != nil {
			t.Errorf("Database.GetGroupForEntry() expected nil for non-existent entry, got %v", gotGroup)
		}
	})
}

func TestDatabase_GetEntryPath(t *testing.T) {
	entry1 := makeEntry("Entry1")
	entry2 := makeEntry("Entry2")
	entry3 := makeEntry("Entry3")
	entry4 := makeEntry("Entry4")

	group1 := makeGroup("Group1", entry1)
	group2 := makeGroup("Group2", entry2)
	nestedGroup := makeGroup("NestedGroup", entry3)
	group1.Groups = append(group1.Groups, nestedGroup)

	db := makeDatabase("test.kdbx", group1, group2)

	tests := []struct {
		name     string
		group    *Group
		entry    *Entry
		wantPath string
	}{
		{"gets path for entry in top level group", &group1, &entry1, "/Group1/Entry1"},
		{"gets path for entry in different top level group", &group2, &entry2, "/Group2/Entry2"},
		{"gets path for entry in nested group", &nestedGroup, &entry3, "/Group1/NestedGroup/Entry3"},
		{"gets path for entry in no groups", &nestedGroup, &entry4, "/Group1/NestedGroup/Entry4"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotPath, _ := db.MakeEntryPath(tt.entry, tt.group)
			if gotPath != tt.wantPath {
				t.Errorf("Database.GetEntryPath() = %v, want %v", gotPath, tt.wantPath)
			}
		})
	}

	t.Run("returns UNKNOWN for non-existent group", func(t *testing.T) {
		nonExistentGroup := makeGroup("NonExistent")
		gotPath, err := db.MakeEntryPath(&entry1, &nonExistentGroup)
		if err == nil {
			t.Errorf("Database.GetEntryPath() expected err for non-existent group, got %v", gotPath)
		}
	})
}

func TestDatabase_RemoveEntry(t *testing.T) {
	t.Run("removes entry from top level group", func(t *testing.T) {
		entry1 := makeEntry("entry1")
		entry2 := makeEntry("entry2")
		group := makeGroup("Group1", entry1, entry2)
		db := makeDatabase("test.kdbx", group)

		// Verify entry exists before removal
		if len(db.Content.Root.Groups[0].Entries) != 2 {
			t.Fatalf("Expected 2 entries before removal, got %d", len(db.Content.Root.Groups[0].Entries))
		}

		// Remove entry1
		err := db.RemoveEntry(entry1.UUID)
		if err != nil {
			t.Fatalf("RemoveEntry() error = %v", err)
		}

		// Verify entry was actually removed
		if len(db.Content.Root.Groups[0].Entries) != 1 {
			t.Errorf("Expected 1 entry after removal, got %d", len(db.Content.Root.Groups[0].Entries))
		}

		// Verify the correct entry was removed
		if db.Content.Root.Groups[0].Entries[0].UUID.Compare(entry1.UUID) {
			t.Error("Removed wrong entry - entry1 still exists")
		}
		if !db.Content.Root.Groups[0].Entries[0].UUID.Compare(entry2.UUID) {
			t.Error("entry2 should still exist")
		}
	})

	t.Run("removes entry from nested group", func(t *testing.T) {
		entry1 := makeEntry("entry1")
		entry2 := makeEntry("entry2")
		nestedGroup := makeGroup("NestedGroup", entry1, entry2)
		topGroup := makeGroup("TopGroup")
		topGroup.Groups = append(topGroup.Groups, nestedGroup)
		db := makeDatabase("test.kdbx", topGroup)

		// Remove entry from nested group
		err := db.RemoveEntry(entry1.UUID)
		if err != nil {
			t.Fatalf("RemoveEntry() error = %v", err)
		}

		// Verify entry was removed from nested group
		if len(db.Content.Root.Groups[0].Groups[0].Entries) != 1 {
			t.Errorf("Expected 1 entry in nested group after removal, got %d", len(db.Content.Root.Groups[0].Groups[0].Entries))
		}
	})

	t.Run("returns error for non-existent entry", func(t *testing.T) {
		entry := makeEntry("entry1")
		group := makeGroup("Group1", entry)
		db := makeDatabase("test.kdbx", group)

		// Try to remove non-existent entry
		nonExistentUUID := gokeepasslib.NewUUID()
		err := db.RemoveEntry(nonExistentUUID)
		if err == nil {
			t.Error("Expected error when removing non-existent entry, got nil")
		}
	})
}
