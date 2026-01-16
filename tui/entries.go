package tui

import (
	"fmt"

	"github.com/gdamore/tcell/v2"
	"github.com/gdamore/tcell/v2/views"

	"github.com/shikaan/keydex/pkg/log"
	"github.com/shikaan/keydex/tui/components"
	"github.com/shikaan/keydex/tui/components/autocomplete"
)

type EntriesView struct {
	autoComplete *autocomplete.AutoComplete
	components.Container
}

func (lv *EntriesView) HandleEvent(ev tcell.Event) bool {
	switch ev := ev.(type) {
	case *tcell.EventKey:
		if ev.Name() == "Ctrl+D" {
			if App.State.IsReadOnly {
				msg := "Could not delete. Archive in read-only mode."
				App.Notify(msg)
				log.Info(msg)
				return true
			}

			entry := App.State.Database.GetFirstEntryByPath(lv.autoComplete.CurrentEntry)
			if entry == nil {
				msg := "Could not delete. Entry cannot be found."
				App.Notify(msg)
				log.Error(msg, nil)
			}

			title := entry.GetTitle()

			App.Confirm(
				"Delete \""+title+"\"? This cannot be undone.",
				func() {
					err := App.State.Database.RemoveEntry(entry.UUID)
					if err != nil {
						msg := "Could not delete. Entry cannot be found."
						App.Notify(msg)
						log.Error(msg, err)
						return
					}

					if e := App.State.Database.SaveAndUnlockEntries(); e != nil {
						App.State.IsReadOnly = true
						msg := "Could not save. Switching to read-only to not corrupt data."
						App.Notify(msg)
						log.Error(msg, e)
						return
					}

					msg := fmt.Sprintf("Entry \"%s\" deleted successfully.", title)
					App.Notify(msg)
					log.Info(msg)

					App.NavigateTo(NewEntryListView)
				}, func() {
					msg := "Operation cancelled. Entry was not deleted."
					App.Notify(msg)
					log.Info(msg)
				},
			)
		}
	}

	return lv.autoComplete.HandleEvent(ev)
}

func NewEntryListView(screen tcell.Screen) views.Widget {
	App.SetTitle("Search")
	view := &EntriesView{}
	view.Container = components.Container{}
	paths := App.State.Database.GetEntryPaths()
	count := len(paths)
	maxX, maxY := getBoundaries(screen)

	autoCompleteOptions := autocomplete.AutoCompleteOptions{
		Screen:     screen,
		Entries:    paths,
		TotalCount: count,
		MaxX:       maxX,
		MaxY:       maxY,
		OnSelect: func(ref string) bool {
			App.State.Reference = ref
			App.State.Entry = App.State.Database.GetFirstEntryByPath(ref)
			App.State.Group = App.State.Database.GetGroupForEntry(App.State.Entry)
			App.NavigateTo(NewEntryView)
			return true
		},
	}

	autoComplete := autocomplete.NewAutoComplete(autoCompleteOptions)
	autoComplete.OnFocus(func() bool {
		App.LastFocused = autoComplete
		return true
	})
	autoComplete.SetFocus(true)
	view.SetContent(autoComplete)
	view.autoComplete = autoComplete

	return view
}

func getBoundaries(screen tcell.Screen) (int, int) {
	x, y := screen.Size()

	// one third of the screen width
	// all the height - title, status, search, counter, notification,
	// and 4 more lines of buffer just in case
	return max(x/3, components.MIN_WIDTH), y - 10
}
