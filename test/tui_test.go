package test

import (
	"os"
	"strings"
	"testing"
	"time"

	"github.com/gdamore/tcell/v2"
	"github.com/shikaan/keydex/pkg/kdbx"
	"github.com/shikaan/keydex/tui"
	"github.com/tobischo/gokeepasslib/v3"
	"github.com/tobischo/gokeepasslib/v3/wrappers"
)

const e2eTimeout = 3 * time.Second
const e2ePassword = "test-password"
const ghUser = "ghuser"
const ghPassword = "ghpass123"
const glUser = "gluser"
const glPassword = "glpass123"

// makeTestKdbxFile creates a temp .kdbx file with a known structure:
//
//	TestDB (root group)
//	  └── Coding (subgroup)
//	        ├── GitHub  (user: ghuser, pass: ghpass123)
//	        └── GitLab  (user: gluser, pass: glpass456)
func makeTestKdbxFile(t *testing.T) (filePath string, password string) {
	t.Helper()

	tmpFile, err := os.CreateTemp(t.TempDir(), "keydex-e2e-*.kdbx")
	if err != nil {
		t.Fatalf("failed to create temp file: %v", err)
	}

	db := gokeepasslib.NewDatabase()
	db.Credentials = gokeepasslib.NewPasswordCredentials(e2ePassword)

	rootGroup := gokeepasslib.NewGroup()
	rootGroup.Name = "TestDB"
	rootGroup.Entries = make([]gokeepasslib.Entry, 0)

	codingGroup := gokeepasslib.NewGroup()
	codingGroup.Name = "Coding"
	codingGroup.Entries = make([]gokeepasslib.Entry, 0)
	codingGroup.Groups = make([]gokeepasslib.Group, 0)

	github := gokeepasslib.NewEntry()
	github.Values = append(github.Values,
		gokeepasslib.ValueData{Key: "Title", Value: gokeepasslib.V{Content: "GitHub"}},
		gokeepasslib.ValueData{Key: "UserName", Value: gokeepasslib.V{Content: ghUser}},
		gokeepasslib.ValueData{Key: "Password", Value: gokeepasslib.V{Content: ghPassword, Protected: wrappers.NewBoolWrapper(true)}},
	)

	gitlab := gokeepasslib.NewEntry()
	gitlab.Values = append(gitlab.Values,
		gokeepasslib.ValueData{Key: "Title", Value: gokeepasslib.V{Content: "GitLab"}},
		gokeepasslib.ValueData{Key: "UserName", Value: gokeepasslib.V{Content: glUser}},
		gokeepasslib.ValueData{Key: "Password", Value: gokeepasslib.V{Content: glPassword, Protected: wrappers.NewBoolWrapper(true)}},
	)

	codingGroup.Entries = append(codingGroup.Entries, github, gitlab)
	rootGroup.Groups = append(rootGroup.Groups, codingGroup)
	db.Content.Root.Groups = []gokeepasslib.Group{rootGroup}

	if err := db.LockProtectedEntries(); err != nil {
		t.Fatalf("failed to lock entries: %v", err)
	}
	if err := gokeepasslib.NewEncoder(tmpFile).Encode(db); err != nil {
		t.Fatalf("failed to encode kdbx: %v", err)
	}
	tmpFile.Close()

	return tmpFile.Name(), e2ePassword
}

func openTestDatabase(t *testing.T, filePath, password string) *kdbx.Database {
	t.Helper()
	db, err := kdbx.New(filePath, password, "")
	if err != nil {
		t.Fatalf("failed to open test database: %v", err)
	}
	return db
}

func startApp(t *testing.T, state tui.State, readOnly bool) tcell.SimulationScreen {
	t.Helper()

	screen := tcell.NewSimulationScreen("UTF-8")
	if err := screen.Init(); err != nil {
		t.Fatalf("failed to init simulation screen: %v", err)
	}
	screen.SetSize(80, 24)

	// Reset the global App
	tui.App = &tui.Application{}
	tui.Setup(screen, state, readOnly)
	go tui.App.Run()

	waitFor(t, screen, "Help", e2eTimeout)
	return screen
}

func startAppWithRef(t *testing.T, db *kdbx.Database, readOnly bool) tcell.SimulationScreen {
	t.Helper()

	screen := tcell.NewSimulationScreen("UTF-8")
	if err := screen.Init(); err != nil {
		t.Fatalf("failed to init simulation screen: %v", err)
	}
	screen.SetSize(80, 24)

	// Reset the global App
	tui.App = &tui.Application{}

	group := &db.GetRootGroup().Groups[0]
	entry := db.GetEntry(group.Entries[0].UUID)
	var ref string
	for _, path := range db.GetEntryPaths() {
		if strings.Contains(path, entry.GetTitle()) {
			ref = path
			break
		}
	}
	state := tui.State{
		Database:  db,
		Group:     group,
		Entry:     entry,
		Reference: ref,
	}

	tui.Setup(screen, state, readOnly)
	go tui.App.Run()

	waitFor(t, screen, entry.GetTitle(), e2eTimeout)
	return screen
}

func readScreen(screen tcell.SimulationScreen) string {
	cells, width, height := screen.GetContents()
	var b strings.Builder
	for y := range height {
		for x := range width {
			idx := y*width + x
			if idx >= len(cells) {
				b.WriteRune(' ')
				continue
			}
			cell := cells[idx]
			if len(cell.Runes) > 0 && cell.Runes[0] != 0 {
				b.WriteRune(cell.Runes[0])
			} else {
				b.WriteRune(' ')
			}
		}
		b.WriteRune('\n')
	}
	return b.String()
}

func readField(t *testing.T, screen tcell.SimulationScreen, field string) string {
	t.Helper()
	s := readScreen(screen)
	for line := range strings.Lines(s) {
		i := strings.Index(line, field)
		if i > -1 {
			i = i + len(field) + 1 // account for the ":" char
			return strings.TrimSpace(line[i:])
		}
	}
	return ""
}

func waitFor(t *testing.T, screen tcell.SimulationScreen, text string, timeout time.Duration) {
	t.Helper()
	deadline := time.Now().Add(timeout)
	for time.Now().Before(deadline) {
		if strings.Contains(readScreen(screen), text) {
			return
		}
		time.Sleep(20 * time.Millisecond)
	}
	t.Fatalf("timed out waiting for %q on screen.\nScreen content:\n%s", text, readScreen(screen))
}

func waitForAbsent(t *testing.T, screen tcell.SimulationScreen, text string, timeout time.Duration) {
	t.Helper()
	deadline := time.Now().Add(timeout)
	for time.Now().Before(deadline) {
		if !strings.Contains(readScreen(screen), text) {
			return
		}
	}
	t.Fatalf("text %q still present on screen after timeout.\nScreen content:\n%s", text, readScreen(screen))
}

func typeText(screen tcell.SimulationScreen, text string) {
	for _, r := range text {
		screen.InjectKey(tcell.KeyRune, r, 0)
	}
}

func navigateToEntryList(t *testing.T, screen tcell.SimulationScreen) {
	t.Helper()
	screen.InjectKey(tcell.KeyCtrlP, 0, tcell.ModCtrl)
	waitFor(t, screen, "Search", e2eTimeout)
}

func selectEntry(t *testing.T, screen tcell.SimulationScreen, searchText string) {
	t.Helper()
	typeText(screen, searchText)
	screen.InjectKey(tcell.KeyEnter, 0, 0)
}

func TestCreateEntryAndSave(t *testing.T) {
	filePath, password := makeTestKdbxFile(t)
	db := openTestDatabase(t, filePath, password)
	screen := startApp(t, tui.State{Database: db}, false)

	// Create new entry (^N)
	screen.InjectKey(tcell.KeyCtrlN, 0, tcell.ModCtrl)
	waitFor(t, screen, "New", e2eTimeout)
	waitFor(t, screen, "[MODIFIED]", e2eTimeout)

	// Navigate to Password field (Title -> UserName -> Password)
	screen.InjectKey(tcell.KeyDown, 0, 0)
	screen.InjectKey(tcell.KeyDown, 0, 0)
	// Reveal password (^R)
	screen.InjectKey(tcell.KeyCtrlR, 0, tcell.ModCtrl)
	waitForAbsent(t, screen, "********", e2eTimeout)
	// Password does not include "=="
	waitForAbsent(t, screen, "==", e2eTimeout)

	// Save (^O) -> Confirm (Y)
	screen.InjectKey(tcell.KeyCtrlO, 0, tcell.ModCtrl)
	waitFor(t, screen, "Create", e2eTimeout)
	screen.InjectKey(tcell.KeyRune, 'Y', 0)
	waitFor(t, screen, "created successfully", e2eTimeout)
	waitForAbsent(t, screen, "[MODIFIED]", e2eTimeout)

	// Open entry list (^P) and verify
	screen.InjectKey(tcell.KeyCtrlP, 0, tcell.ModCtrl)
	waitFor(t, screen, tui.App.State.Reference, e2eTimeout)
}

func TestCreateEntryAndDismissSave(t *testing.T) {
	filePath, password := makeTestKdbxFile(t)
	db := openTestDatabase(t, filePath, password)
	screen := startApp(t, tui.State{Database: db}, false)

	screen.InjectKey(tcell.KeyCtrlN, 0, tcell.ModCtrl)
	waitFor(t, screen, "New", e2eTimeout)

	// Save -> Dismiss (N)
	screen.InjectKey(tcell.KeyCtrlO, 0, tcell.ModCtrl)
	waitFor(t, screen, "Create", e2eTimeout)
	screen.InjectKey(tcell.KeyRune, 'N', 0)
	waitFor(t, screen, "not created", e2eTimeout)

	// Open entry list (^P) and verify
	screen.InjectKey(tcell.KeyCtrlP, 0, tcell.ModCtrl)
	waitForAbsent(t, screen, tui.App.State.Reference, e2eTimeout)
}

func TestCreateEntryModifyThenCancel(t *testing.T) {
	filePath, password := makeTestKdbxFile(t)
	db := openTestDatabase(t, filePath, password)
	screen := startApp(t, tui.State{Database: db}, false)

	screen.InjectKey(tcell.KeyCtrlN, 0, tcell.ModCtrl)
	waitFor(t, screen, "New", e2eTimeout)

	// Type in a field
	typeText(screen, "TestEntry")
	waitFor(t, screen, "[MODIFIED]", e2eTimeout)

	// Cancel (ESC) — [MODIFIED] stays because it's a new entry
	screen.InjectKey(tcell.KeyEsc, 0, 0)
	waitFor(t, screen, "[MODIFIED]", e2eTimeout)
	waitFor(t, screen, "Help", e2eTimeout)
}

func TestViewEntryAndRevealPassword(t *testing.T) {
	filePath, password := makeTestKdbxFile(t)
	db := openTestDatabase(t, filePath, password)
	screen := startApp(t, tui.State{Database: db}, false)

	navigateToEntryList(t, screen)
	selectEntry(t, screen, "GitHub")
	waitFor(t, screen, "GitHub", e2eTimeout)
	waitFor(t, screen, ghUser, e2eTimeout)

	// Password should be masked initially
	waitFor(t, screen, "********", e2eTimeout)

	// Navigate to Password field (Title → UserName → Password)
	screen.InjectKey(tcell.KeyDown, 0, 0)
	screen.InjectKey(tcell.KeyDown, 0, 0)

	// Reveal password (^R)
	screen.InjectKey(tcell.KeyCtrlR, 0, tcell.ModCtrl)
	waitFor(t, screen, ghPassword, e2eTimeout)
}

func TestViewEntryModifyThenCancel(t *testing.T) {
	filePath, password := makeTestKdbxFile(t)
	db := openTestDatabase(t, filePath, password)
	screen := startApp(t, tui.State{Database: db}, false)

	navigateToEntryList(t, screen)
	selectEntry(t, screen, "GitHub")
	waitFor(t, screen, "GitHub", e2eTimeout)
	waitFor(t, screen, ghUser, e2eTimeout)

	// Select User field -> Delete content -> Type New Content
	screen.InjectKey(tcell.KeyDown, 0, 0)
	for _ = range len(ghUser) {
		screen.InjectKey(tcell.KeyDelete, 0, 0)
	}
	typeText(screen, "Modified")
	waitFor(t, screen, "[MODIFIED]", e2eTimeout)
	waitForAbsent(t, screen, ghUser, e2eTimeout)

	// Cancel (ESC) — entry returns to previous state
	screen.InjectKey(tcell.KeyEsc, 0, 0)
	waitFor(t, screen, "GitHub", e2eTimeout)
	waitForAbsent(t, screen, "[MODIFIED]", e2eTimeout)
}

func TestViewEntryModifyHiddenFieldThenCancel(t *testing.T) {
	filePath, password := makeTestKdbxFile(t)
	db := openTestDatabase(t, filePath, password)
	screen := startApp(t, tui.State{Database: db}, false)

	navigateToEntryList(t, screen)
	selectEntry(t, screen, "GitHub")
	waitFor(t, screen, "GitHub", e2eTimeout)

	screen.InjectKey(tcell.KeyDown, 0, 0)
	screen.InjectKey(tcell.KeyDown, 0, 0)
	typeText(screen, "Modified")
	// Notify that needs reveal
	waitFor(t, screen, "Reveal", e2eTimeout)
	// Reveal
	screen.InjectKey(tcell.KeyCtrlR, 0, tcell.ModCtrl)
	// Password did not change
	waitFor(t, screen, ghPassword, e2eTimeout)
	typeText(screen, "Modified")
	// Now password has changed
	waitFor(t, screen, "Modified", e2eTimeout)

	// Cancel (ESC) — entry returns to previous state
	screen.InjectKey(tcell.KeyEsc, 0, 0)
	waitForAbsent(t, screen, "[MODIFIED]", e2eTimeout)
	waitFor(t, screen, "********", e2eTimeout)
	screen.InjectKey(tcell.KeyDown, 0, 0)
	screen.InjectKey(tcell.KeyDown, 0, 0)
	screen.InjectKey(tcell.KeyCtrlR, 0, tcell.ModCtrl)
	// Password was restored
	waitFor(t, screen, ghPassword, e2eTimeout)
}

func TestViewEntryUpdateTitle(t *testing.T) {
	filePath, password := makeTestKdbxFile(t)
	db := openTestDatabase(t, filePath, password)
	screen := startApp(t, tui.State{Database: db}, false)

	navigateToEntryList(t, screen)
	selectEntry(t, screen, "GitHub")
	waitFor(t, screen, "GitHub", e2eTimeout)

	// Update title
	typeText(screen, "asdf1234")
	waitFor(t, screen, "[MODIFIED]", e2eTimeout)

	// Save -> Confirm
	screen.InjectKey(tcell.KeyCtrlO, 0, tcell.ModCtrl)
	waitFor(t, screen, "Save changes", e2eTimeout)
	screen.InjectKey(tcell.KeyRune, 'Y', 0)
	waitFor(t, screen, "saved successfully", e2eTimeout)
	waitForAbsent(t, screen, "[MODIFIED]", e2eTimeout)

	navigateToEntryList(t, screen)
	waitFor(t, screen, "asdf1234", e2eTimeout)
}

func TestViewEntryModifyThenSave(t *testing.T) {
	filePath, password := makeTestKdbxFile(t)
	db := openTestDatabase(t, filePath, password)
	screen := startApp(t, tui.State{Database: db}, false)

	navigateToEntryList(t, screen)
	selectEntry(t, screen, "GitHub")
	waitFor(t, screen, "GitHub", e2eTimeout)

	// Navigate to UserName field
	screen.InjectKey(tcell.KeyDown, 0, 0)
	typeText(screen, "2")
	waitFor(t, screen, "[MODIFIED]", e2eTimeout)

	// Save -> Confirm
	screen.InjectKey(tcell.KeyCtrlO, 0, tcell.ModCtrl)
	waitFor(t, screen, "Save changes", e2eTimeout)
	screen.InjectKey(tcell.KeyRune, 'Y', 0)
	waitFor(t, screen, "saved successfully", e2eTimeout)
	waitForAbsent(t, screen, "[MODIFIED]", e2eTimeout)

	waitFor(t, screen, "saved successfully", e2eTimeout)
}

func TestViewEntryUpdateGroupAndSave(t *testing.T) {
	filePath, password := makeTestKdbxFile(t)
	db := openTestDatabase(t, filePath, password)
	screen := startApp(t, tui.State{Database: db}, false)

	navigateToEntryList(t, screen)
	selectEntry(t, screen, "GitHub")
	waitFor(t, screen, "GitHub", e2eTimeout)

	// Change group (^K)
	screen.InjectKey(tcell.KeyCtrlK, 0, tcell.ModCtrl)
	waitFor(t, screen, "Select group", e2eTimeout)

	// Select the root "TestDB" group
	typeText(screen, "TestDB")
	time.Sleep(100 * time.Millisecond)
	screen.InjectKey(tcell.KeyEnter, 0, 0)
	waitFor(t, screen, "[MODIFIED]", e2eTimeout)

	// Save → Confirm
	screen.InjectKey(tcell.KeyCtrlO, 0, tcell.ModCtrl)
	waitFor(t, screen, "Save changes", e2eTimeout)
	screen.InjectKey(tcell.KeyRune, 'Y', 0)
	waitFor(t, screen, "saved successfully", e2eTimeout)

	screen.InjectKey(tcell.KeyCtrlP, 0, tcell.ModCtrl)
	waitFor(t, screen, "TestDB/GitHub", e2eTimeout)
}

func TestViewEntryUpdateGroupAndCancel(t *testing.T) {
	filePath, password := makeTestKdbxFile(t)
	db := openTestDatabase(t, filePath, password)
	screen := startApp(t, tui.State{Database: db}, false)

	navigateToEntryList(t, screen)
	selectEntry(t, screen, "GitHub")
	waitFor(t, screen, "GitHub", e2eTimeout)
	waitFor(t, screen, "Coding", e2eTimeout)

	// Change group (^K)
	screen.InjectKey(tcell.KeyCtrlK, 0, tcell.ModCtrl)
	waitFor(t, screen, "Select group", e2eTimeout)

	typeText(screen, "TestDB")
	time.Sleep(100 * time.Millisecond)
	screen.InjectKey(tcell.KeyEnter, 0, 0)
	waitFor(t, screen, "[MODIFIED]", e2eTimeout)
	waitFor(t, screen, "TestDB", e2eTimeout)
	waitForAbsent(t, screen, "Coding", e2eTimeout)

	// Cancel
	screen.InjectKey(tcell.KeyEsc, 0, 0)

	// Should return to entry view with original group
	waitFor(t, screen, "Coding", e2eTimeout)
}

func TestViewEntryDismissDeletion(t *testing.T) {
	filePath, password := makeTestKdbxFile(t)
	db := openTestDatabase(t, filePath, password)
	screen := startApp(t, tui.State{Database: db}, false)

	navigateToEntryList(t, screen)
	selectEntry(t, screen, "GitHub")
	waitFor(t, screen, "GitHub", e2eTimeout)

	// Delete -> Say No
	screen.InjectKey(tcell.KeyCtrlD, 0, tcell.ModCtrl)
	waitFor(t, screen, "Delete", e2eTimeout)
	screen.InjectKey(tcell.KeyRune, 'N', 0)
	waitFor(t, screen, "not deleted", e2eTimeout)

	// Verify entry still in list
	navigateToEntryList(t, screen)
	waitFor(t, screen, "GitHub", e2eTimeout)
}

func TestListEntriesDismissDeletion(t *testing.T) {
	filePath, password := makeTestKdbxFile(t)
	db := openTestDatabase(t, filePath, password)
	screen := startApp(t, tui.State{Database: db}, false)

	navigateToEntryList(t, screen)
	// First entry is auto-focused; delete from list view
	screen.InjectKey(tcell.KeyCtrlD, 0, tcell.ModCtrl)
	waitFor(t, screen, "Delete", e2eTimeout)
	screen.InjectKey(tcell.KeyRune, 'N', 0)
	waitFor(t, screen, "not deleted", e2eTimeout)

	waitFor(t, screen, "GitHub", e2eTimeout)
}

func TestViewEntryAndDelete(t *testing.T) {
	filePath, password := makeTestKdbxFile(t)
	db := openTestDatabase(t, filePath, password)
	screen := startApp(t, tui.State{Database: db}, false)

	navigateToEntryList(t, screen)
	selectEntry(t, screen, "GitHub")
	waitFor(t, screen, "GitHub", e2eTimeout)

	// Delete → Confirm
	screen.InjectKey(tcell.KeyCtrlD, 0, tcell.ModCtrl)
	waitFor(t, screen, "Delete", e2eTimeout)
	screen.InjectKey(tcell.KeyRune, 'Y', 0)
	waitFor(t, screen, "deleted successfully", e2eTimeout)

	// Should navigate to entry list
	waitFor(t, screen, "Search", e2eTimeout)
	waitForAbsent(t, screen, "Coding/GitHub", e2eTimeout)
}

func TestListEntriesAndDelete(t *testing.T) {
	filePath, password := makeTestKdbxFile(t)
	db := openTestDatabase(t, filePath, password)
	screen := startApp(t, tui.State{Database: db}, false)

	navigateToEntryList(t, screen)
	typeText(screen, "GitHub")
	screen.InjectKey(tcell.KeyCtrlD, 0, tcell.ModCtrl)
	waitFor(t, screen, "Delete", e2eTimeout)
	screen.InjectKey(tcell.KeyRune, 'Y', 0)
	waitFor(t, screen, "deleted successfully", e2eTimeout)
	waitForAbsent(t, screen, "Coding/GitHub", e2eTimeout)
}

func TestGroupCancelsUpdates(t *testing.T) {
	filePath, password := makeTestKdbxFile(t)
	db := openTestDatabase(t, filePath, password)
	screen := startApp(t, tui.State{Database: db}, false)

	navigateToEntryList(t, screen)
	selectEntry(t, screen, "GitHub")
	waitFor(t, screen, "GitHub", e2eTimeout)

	// Make the entry dirty
	typeText(screen, "x")
	waitFor(t, screen, "[MODIFIED]", e2eTimeout)

	// Try to change group — should show confirmation
	screen.InjectKey(tcell.KeyCtrlK, 0, tcell.ModCtrl)
	waitFor(t, screen, "Navigate away", e2eTimeout)
}

func TestReadOnly(t *testing.T) {
	filePath, password := makeTestKdbxFile(t)
	db := openTestDatabase(t, filePath, password)
	screen := startApp(t, tui.State{Database: db}, true)

	navigateToEntryList(t, screen)
	selectEntry(t, screen, "GitHub")

	waitFor(t, screen, "READ ONLY", e2eTimeout)
	// Try ^K (change group)
	screen.InjectKey(tcell.KeyCtrlK, 0, tcell.ModCtrl)
	waitFor(t, screen, "read-only", e2eTimeout)

	// Try ^O (save)
	screen.InjectKey(tcell.KeyCtrlO, 0, tcell.ModCtrl)
	waitFor(t, screen, "read-only", e2eTimeout)

	// Try ^D (delete)
	screen.InjectKey(tcell.KeyCtrlD, 0, tcell.ModCtrl)
	waitFor(t, screen, "read-only", e2eTimeout)
}

func TestReadOnlyWithRef(t *testing.T) {
	filePath, password := makeTestKdbxFile(t)
	db := openTestDatabase(t, filePath, password)
	screen := startAppWithRef(t, db, true)

	waitFor(t, screen, "READ ONLY", e2eTimeout)
	// Try ^K (change group)
	screen.InjectKey(tcell.KeyCtrlK, 0, tcell.ModCtrl)
	waitFor(t, screen, "read-only", e2eTimeout*2)

	// Try ^O (save)
	screen.InjectKey(tcell.KeyCtrlO, 0, tcell.ModCtrl)
	waitFor(t, screen, "read-only", e2eTimeout)

	// Try ^D (delete)
	screen.InjectKey(tcell.KeyCtrlD, 0, tcell.ModCtrl)
	waitFor(t, screen, "read-only", e2eTimeout)
}

func TestViewAndRevealPasswordWithRef(t *testing.T) {
	filePath, password := makeTestKdbxFile(t)
	db := openTestDatabase(t, filePath, password)
	screen := startAppWithRef(t, db, true)

	waitFor(t, screen, "********", e2eTimeout)

	// Navigate to Password field
	screen.InjectKey(tcell.KeyDown, 0, 0)
	screen.InjectKey(tcell.KeyDown, 0, 0)

	screen.InjectKey(tcell.KeyCtrlR, 0, tcell.ModCtrl)
	waitFor(t, screen, "ghpass123", e2eTimeout)
}

func TestViewModifyThenCancelWithRef(t *testing.T) {
	filePath, password := makeTestKdbxFile(t)
	db := openTestDatabase(t, filePath, password)
	screen := startAppWithRef(t, db, false)

	screen.InjectKey(tcell.KeyDown, 0, 0)
	for _ = range len(ghUser) {
		screen.InjectKey(tcell.KeyDelete, 0, 0)
	}
	typeText(screen, "Modified")
	waitFor(t, screen, "[MODIFIED]", e2eTimeout)
	waitForAbsent(t, screen, ghUser, e2eTimeout)

	// Cancel (ESC) — entry returns to previous state
	screen.InjectKey(tcell.KeyEsc, 0, 0)
	waitFor(t, screen, "GitHub", e2eTimeout)
	waitForAbsent(t, screen, "[MODIFIED]", e2eTimeout)
}

func TestViewUpdateTitlelWithRef(t *testing.T) {
	filePath, password := makeTestKdbxFile(t)
	db := openTestDatabase(t, filePath, password)
	screen := startAppWithRef(t, db, false)

	// Update title
	typeText(screen, "asdf1234")
	waitFor(t, screen, "[MODIFIED]", e2eTimeout)

	// Save -> Confirm
	screen.InjectKey(tcell.KeyCtrlO, 0, tcell.ModCtrl)
	waitFor(t, screen, "Save changes", e2eTimeout)
	screen.InjectKey(tcell.KeyRune, 'Y', 0)
	waitFor(t, screen, "saved successfully", e2eTimeout)
	waitForAbsent(t, screen, "[MODIFIED]", e2eTimeout)

	navigateToEntryList(t, screen)
	waitFor(t, screen, "asdf1234", e2eTimeout)
}

func TestViewModifyThenSaveWithRef(t *testing.T) {
	filePath, password := makeTestKdbxFile(t)
	db := openTestDatabase(t, filePath, password)
	screen := startAppWithRef(t, db, false)

	// Navigate to UserName field
	screen.InjectKey(tcell.KeyDown, 0, 0)
	typeText(screen, "2")
	waitFor(t, screen, "[MODIFIED]", e2eTimeout)

	// Save -> Confirm
	screen.InjectKey(tcell.KeyCtrlO, 0, tcell.ModCtrl)
	waitFor(t, screen, "Save changes", e2eTimeout)
	screen.InjectKey(tcell.KeyRune, 'Y', 0)
	waitFor(t, screen, "saved successfully", e2eTimeout)
	waitForAbsent(t, screen, "[MODIFIED]", e2eTimeout)

	waitFor(t, screen, "saved successfully", e2eTimeout)
}

func TestUpdateGroupAndSaveWithRef(t *testing.T) {
	filePath, password := makeTestKdbxFile(t)
	db := openTestDatabase(t, filePath, password)
	screen := startAppWithRef(t, db, false)

	// Change group (^K)
	screen.InjectKey(tcell.KeyCtrlK, 0, tcell.ModCtrl)
	waitFor(t, screen, "Select group", e2eTimeout)

	// Select the root "TestDB" group
	typeText(screen, "TestDB")
	time.Sleep(100 * time.Millisecond)
	screen.InjectKey(tcell.KeyEnter, 0, 0)
	waitFor(t, screen, "[MODIFIED]", e2eTimeout)

	// Save → Confirm
	screen.InjectKey(tcell.KeyCtrlO, 0, tcell.ModCtrl)
	waitFor(t, screen, "Save changes", e2eTimeout)
	screen.InjectKey(tcell.KeyRune, 'Y', 0)
	waitFor(t, screen, "saved successfully", e2eTimeout)

	screen.InjectKey(tcell.KeyCtrlP, 0, tcell.ModCtrl)
	waitFor(t, screen, "TestDB/GitHub", e2eTimeout)
}

func TestViewUpdateGroupAndCancelWithRef(t *testing.T) {
	filePath, password := makeTestKdbxFile(t)
	db := openTestDatabase(t, filePath, password)
	screen := startAppWithRef(t, db, false)

	// Change group (^K)
	screen.InjectKey(tcell.KeyCtrlK, 0, tcell.ModCtrl)
	waitFor(t, screen, "Select group", e2eTimeout)

	typeText(screen, "TestDB")
	time.Sleep(100 * time.Millisecond)
	screen.InjectKey(tcell.KeyEnter, 0, 0)
	waitFor(t, screen, "[MODIFIED]", e2eTimeout)
	waitFor(t, screen, "TestDB", e2eTimeout)
	waitForAbsent(t, screen, "Coding", e2eTimeout)

	// Cancel
	screen.InjectKey(tcell.KeyEsc, 0, 0)

	// Should return to entry view with original group
	waitFor(t, screen, "Coding", e2eTimeout)
}

func TestViewDismissDeletionWithRef(t *testing.T) {
	filePath, password := makeTestKdbxFile(t)
	db := openTestDatabase(t, filePath, password)
	screen := startAppWithRef(t, db, false)

	navigateToEntryList(t, screen)
	selectEntry(t, screen, "GitHub")
	waitFor(t, screen, "GitHub", e2eTimeout)

	// Delete -> Say No
	screen.InjectKey(tcell.KeyCtrlD, 0, tcell.ModCtrl)
	waitFor(t, screen, "Delete", e2eTimeout)
	screen.InjectKey(tcell.KeyRune, 'N', 0)
	waitFor(t, screen, "not deleted", e2eTimeout)

	// Verify entry still in list
	navigateToEntryList(t, screen)
	waitFor(t, screen, "GitHub", e2eTimeout)
}

func TestEntryListShowsAllEntries(t *testing.T) {
	filePath, password := makeTestKdbxFile(t)
	db := openTestDatabase(t, filePath, password)
	screen := startApp(t, tui.State{Database: db}, false)

	navigateToEntryList(t, screen)

	// Both entries visible
	waitFor(t, screen, "Coding/GitHub", e2eTimeout)
	waitFor(t, screen, "Coding/GitLab", e2eTimeout)
	// Counter shows total
	waitFor(t, screen, "2/2", e2eTimeout)
}

func TestEntryListFuzzySearch(t *testing.T) {
	filePath, password := makeTestKdbxFile(t)
	db := openTestDatabase(t, filePath, password)
	screen := startApp(t, tui.State{Database: db}, false)

	navigateToEntryList(t, screen)

	// Type a query that matches only GitHub
	typeText(screen, "Hub")
	waitFor(t, screen, "Coding/GitHub", e2eTimeout)
	waitForAbsent(t, screen, "Coding/GitLab", e2eTimeout)
	waitFor(t, screen, "1/2", e2eTimeout)
}

func TestEntryListEmptyState(t *testing.T) {
	filePath, password := makeTestKdbxFile(t)
	db := openTestDatabase(t, filePath, password)
	screen := startApp(t, tui.State{Database: db}, false)

	navigateToEntryList(t, screen)

	// Type a query that matches nothing
	typeText(screen, "zzzzz")
	waitFor(t, screen, "No Results", e2eTimeout)
	waitFor(t, screen, "0/2", e2eTimeout)
}

func TestEntryListArrowNavigation(t *testing.T) {
	filePath, password := makeTestKdbxFile(t)
	db := openTestDatabase(t, filePath, password)
	screen := startApp(t, tui.State{Database: db}, false)

	navigateToEntryList(t, screen)

	// First entry is focused by default (has "> " prefix)
	waitFor(t, screen, "> ", e2eTimeout)

	// Move down — focus should shift to second entry
	screen.InjectKey(tcell.KeyDown, 0, 0)
	time.Sleep(50 * time.Millisecond)

	s := readScreen(screen)
	lines := strings.Split(s, "\n")
	firstFocused, secondFocused := false, false
	for _, line := range lines {
		if strings.Contains(line, "> ") && strings.Contains(line, "GitHub") {
			firstFocused = true
		}
		if strings.Contains(line, "> ") && strings.Contains(line, "GitLab") {
			secondFocused = true
		}
	}

	if firstFocused {
		t.Error("first entry should not be focused after pressing Down")
	}
	if !secondFocused {
		t.Error("second entry should be focused after pressing Down")
	}

	// Move up — focus should return to first entry
	screen.InjectKey(tcell.KeyUp, 0, 0)
	time.Sleep(50 * time.Millisecond)

	s = readScreen(screen)
	lines = strings.Split(s, "\n")
	firstFocused, secondFocused = false, false
	for _, line := range lines {
		if strings.Contains(line, "> ") && strings.Contains(line, "GitHub") {
			firstFocused = true
		}
		if strings.Contains(line, "> ") && strings.Contains(line, "GitLab") {
			secondFocused = true
		}
	}

	if !firstFocused {
		t.Error("first entry should be focused after pressing Up")
	}
	if secondFocused {
		t.Error("second entry should not be focused after pressing Up")
	}
}

func TestViewAndDeleteWithRef(t *testing.T) {
	filePath, password := makeTestKdbxFile(t)
	db := openTestDatabase(t, filePath, password)
	screen := startAppWithRef(t, db, false)

	navigateToEntryList(t, screen)
	selectEntry(t, screen, "GitHub")
	waitFor(t, screen, "GitHub", e2eTimeout)

	// Delete → Confirm
	screen.InjectKey(tcell.KeyCtrlD, 0, tcell.ModCtrl)
	waitFor(t, screen, "Delete", e2eTimeout)
	screen.InjectKey(tcell.KeyRune, 'Y', 0)
	waitFor(t, screen, "deleted successfully", e2eTimeout)

	// Should navigate to entry list
	waitFor(t, screen, "Search", e2eTimeout)
	waitForAbsent(t, screen, "Coding/GitHub", e2eTimeout)
}
