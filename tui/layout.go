package tui

import (
	"github.com/gdamore/tcell/v2"
	"github.com/gdamore/tcell/v2/views"
	"github.com/shikaan/keydex/tui/components"
	"github.com/shikaan/keydex/tui/components/status"
)

type Layout struct {
	Status *status.Status
	Title  *components.Title

	Screen tcell.Screen

	views.Panel
}

func (l *Layout) SetContent(w views.Widget) {
	l.Panel.SetContent(w)
	// Make sure the bottom panel is _always_ shown writing it last
	l.Panel.SetStatus(l.Status)
	l.Panel.SetTitle(l.Title)
	l.Resize()
}

func (v *Layout) HandleEvent(ev tcell.Event) bool {
	// If there is a pending confirmation, delegate to panel to handle Y/N/Cancel
	if v.Status.IsConfirming() {
		return v.Panel.HandleEvent(ev)
	}

	switch ev := ev.(type) {
	case *tcell.EventKey:
		if ev.Name() == "Ctrl+X" {
			App.Quit()
			return true
		}
		if ev.Name() == "Ctrl+P" {
			App.NavigateTo(NewEntryListView)
			return true
		}
		if ev.Name() == "Ctrl+G" {
			App.NavigateTo(NewHelpView)
			return true
		}
		if ev.Name() == "Ctrl+N" {
			if App.IsReadOnly() {
				App.Notify("Could not create. Archive in read-only mode.")
				return true
			}
			err := App.CreateEmptyEntry()

			if err != nil {
				App.Notify("Could not create. Check logs for details.")
				return true
			}

			App.NavigateTo(NewEntryView)
			return true
		}
		if ev.Name() == "Ctrl+C" {
			handled := v.Panel.HandleEvent(ev)

			if !handled {
				App.Notify("No field selected for copy. Use ^X to close.")
			}

			return true
		}
		if ev.Key() == tcell.KeyEsc {
			if App.State.Entry == nil {
				App.Notify("No entry selected yet.")
				return true
			}

			if App.IsDirty() {
				App.Notify("Operation cancelled. Updates were not saved.")
				App.SetDirty(false)
			}

			// Group for entry is nil when the entry to be edited has just been created.
			// In that case, we will use the root group.
			group := App.State.Database.GetGroupForEntry(App.State.Entry)
			isNewEntry := group == nil

			if isNewEntry {
				group = App.State.Database.GetRootGroup()

				// If the entry is new, saving the database will incur a write on disk
				// This means the state is dirty, and we make sure it's reflected here.
				App.SetDirty(true)
			}
			App.State.Group = group

			// Needed to reset group selection on cancelled operations
			App.NavigateTo(NewEntryView)
			return true
		}
	}
	return v.Panel.HandleEvent(ev)
}

// Returns a Layout component responsible for the shell of the application
// and of the routing in between pages
func NewLayout(screen tcell.Screen) *Layout {
	l := &Layout{}
	title := components.NewTitle(App.State.Database.Content.Meta.DatabaseName)
	s := status.NewStatus()

	t := views.NewText()
	t.SetText(" ")
	l.SetMenu(t)

	l.SetStatus(s)
	l.Status = s

	l.SetTitle(title)
	l.Title = title

	l.Screen = screen

	return l
}
