package tui

import (
	"testing"

	"github.com/gdamore/tcell/v2"
	"github.com/shikaan/keydex/pkg/kdbx"
	"github.com/tobischo/gokeepasslib/v3"
)

// Helper function to create a test database with basic structure
func createTestDatabase() *kdbx.Database {
	return &kdbx.Database{
		Database: gokeepasslib.Database{
			Content: &gokeepasslib.DBContent{
				Meta: &gokeepasslib.MetaData{
					DatabaseName: "TestDB",
				},
				Root: &gokeepasslib.RootData{
					Groups: []gokeepasslib.Group{
						{
							Name:    "RootGroup",
							Entries: []gokeepasslib.Entry{},
							Groups:  []gokeepasslib.Group{},
						},
					},
				},
			},
		},
	}
}

func TestLayout_HandleEvent_Esc_AlwaysSetsGroup(t *testing.T) {
	tests := []struct {
		name              string
		setupDatabase     func() *kdbx.Database
		setupEntry        func(db *kdbx.Database) *kdbx.Entry
		expectGroupFromDB bool // true if we expect the group from GetGroupForEntry, false for root group
	}{
		{
			name: "sets group from database when entry has a group",
			setupDatabase: func() *kdbx.Database {
				db := createTestDatabase()
				// Add a subgroup to the root group
				db.Content.Root.Groups[0].Groups = []gokeepasslib.Group{
					{
						Name:    "SubGroup",
						Entries: []gokeepasslib.Entry{},
						Groups:  []gokeepasslib.Group{},
					},
				}
				return db
			},
			setupEntry: func(db *kdbx.Database) *kdbx.Entry {
				// Create an entry and add it to SubGroup
				entry := db.NewEntry()
				subGroup := &db.Content.Root.Groups[0].Groups[0]
				subGroup.Entries = append(subGroup.Entries, entry.Entry)
				return entry
			},
			expectGroupFromDB: true,
		},
		{
			name: "sets root group when entry has no group (newly created entry)",
			setupDatabase: func() *kdbx.Database {
				return createTestDatabase()
			},
			setupEntry: func(db *kdbx.Database) *kdbx.Entry {
				// Create an entry but don't add it to any group
				entry := db.NewEntry()
				return entry
			},
			expectGroupFromDB: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Setup test database
			db := tt.setupDatabase()
			entry := tt.setupEntry(db)

			App.State = State{
				Database: db,
				Entry:    entry,
				Group:    nil, // Start with nil group
			}

			// Create layout and set up App
			screen := tcell.NewSimulationScreen("UTF-8")
			screen.Init()
			defer screen.Fini()

			layout := NewLayout(screen)
			App.layout = layout
			App.screen = screen

			// Create ESC key event
			escEvent := tcell.NewEventKey(tcell.KeyEsc, 0, tcell.ModNone)

			// Handle the event
			layout.HandleEvent(escEvent)

			// Verify App.State.Group is set
			if App.State.Group == nil {
				t.Errorf("App.State.Group should not be nil after handling ESC event")
				return
			}

			var expectedGroup *kdbx.Group
			if tt.expectGroupFromDB {
				expectedGroup = db.GetGroupForEntry(entry)
			} else {
				expectedGroup = db.GetRootGroup()
			}

			if expectedGroup == nil {
				t.Fatal("Expected to find root group, but got nil")
			}
			if !App.State.Group.UUID.Compare(expectedGroup.UUID) {
				t.Errorf("Expected root group (name: %s), got group (name: %s)",
					expectedGroup.Name, App.State.Group.Name)
			}
		})
	}
}

func TestLayout_HandleEvent_Esc_RequiresExistingEntry(t *testing.T) {
	// Setup test database
	db := createTestDatabase()

	App.State = State{
		Database: db,
		Entry:    nil, // No entry selected
		Group:    nil,
	}

	// Create layout and set up App
	screen := tcell.NewSimulationScreen("UTF-8")
	screen.Init()
	defer screen.Fini()

	layout := NewLayout(screen)
	App.layout = layout
	App.screen = screen

	// Create ESC key event
	escEvent := tcell.NewEventKey(tcell.KeyEsc, 0, tcell.ModNone)

	// Handle the event
	handled := layout.HandleEvent(escEvent)

	// Verify the event was handled (returns true)
	if !handled {
		t.Error("Expected ESC event to be handled even without entry")
	}

	// Verify App.State.Group is still nil (since we exit early when Entry is nil)
	if App.State.Group != nil {
		t.Error("App.State.Group should remain nil when Entry is nil")
	}
}

func TestLayout_HandleEvent_Esc_ClearsUnsavedChanges(t *testing.T) {
	db := createTestDatabase()
	entry := db.NewEntry()

	App.State = State{
		Database: db,
		Entry:    entry,
		Group:    nil,
	}

	App.isDirty = true

	screen := tcell.NewSimulationScreen("UTF-8")
	screen.Init()
	defer screen.Fini()

	layout := NewLayout(screen)
	App.layout = layout
	App.screen = screen

	escEvent := tcell.NewEventKey(tcell.KeyEsc, 0, tcell.ModNone)

	layout.HandleEvent(escEvent)

	if App.IsDirty() {
		t.Error("Should not be dirty after handling ESC event")
	}

	if App.State.Group == nil {
		t.Error("App.State.Group should not be nil after handling ESC event")
	}
}

func TestLayout_HandleEvent_DelegatesWhenConfirming(t *testing.T) {
	db := createTestDatabase()
	entry := db.NewEntry()

	App.State = State{
		Database: db,
		Entry:    entry,
		Group:    nil,
	}

	screen := tcell.NewSimulationScreen("UTF-8")
	screen.Init()
	defer screen.Fini()

	layout := NewLayout(screen)
	App.layout = layout
	App.screen = screen

	confirmRejectCalled := false
	layout.Status.Confirm("Test confirmation?",
		func() { /* accept handler */ },
		func() { confirmRejectCalled = true },
	)

	if !layout.Status.IsConfirming() {
		t.Fatal("Expected confirmation to be active")
	}

	escEvent := tcell.NewEventKey(tcell.KeyEsc, 0, tcell.ModNone)
	handled := layout.HandleEvent(escEvent)

	if !handled {
		t.Error("Expected event to be handled")
	}

	if !confirmRejectCalled {
		t.Error("Expected confirmation reject callback to be called when ESC is pressed during confirmation")
	}

	if layout.Status.IsConfirming() {
		t.Error("Expected confirmation to be dismissed after ESC")
	}

	if App.State.Group != nil {
		t.Error("App.State.Group should remain nil when ESC is handled during confirmation")
	}
}

func TestLayout_HandleEvent_Esc_ModifiedBadgeBehavior(t *testing.T) {
	tests := []struct {
		name               string
		setupEntry         func(db *kdbx.Database) *kdbx.Entry
		initialDirty       bool
		expectedDirtyAfter bool
		description        string
	}{
		{
			name: "keeps [MODIFIED] badge for new entries",
			setupEntry: func(db *kdbx.Database) *kdbx.Entry {
				// Create entry but don't add it to any group (simulates new entry)
				return db.NewEntry()
			},
			initialDirty:       true,
			expectedDirtyAfter: true,
			description:        "New entries should keep dirty state after ESC",
		},
		{
			name: "clears [MODIFIED] badge for updated entries",
			setupEntry: func(db *kdbx.Database) *kdbx.Entry {
				// Create entry and add it to root group (simulates existing entry)
				entry := db.NewEntry()
				db.Content.Root.Groups[0].Entries = append(db.Content.Root.Groups[0].Entries, entry.Entry)
				return entry
			},
			initialDirty:       true,
			expectedDirtyAfter: false,
			description:        "Updated entries should clear dirty state after ESC",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			db := createTestDatabase()
			entry := tt.setupEntry(db)

			App.State = State{
				Database: db,
				Entry:    entry,
				Group:    nil,
			}
			App.isDirty = tt.initialDirty

			screen := tcell.NewSimulationScreen("UTF-8")
			screen.Init()
			defer screen.Fini()

			layout := NewLayout(screen)
			App.layout = layout
			App.screen = screen

			escEvent := tcell.NewEventKey(tcell.KeyEsc, 0, tcell.ModNone)
			layout.HandleEvent(escEvent)

			if App.IsDirty() != tt.expectedDirtyAfter {
				t.Errorf("%s: expected IsDirty() = %v, got %v",
					tt.description, tt.expectedDirtyAfter, App.IsDirty())
			}
		})
	}
}
