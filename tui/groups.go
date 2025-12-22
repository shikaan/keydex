package tui

import (
	"fmt"

	"github.com/gdamore/tcell/v2"
	"github.com/gdamore/tcell/v2/views"

	"github.com/shikaan/keydex/pkg/log"
	"github.com/shikaan/keydex/tui/components"
	"github.com/shikaan/keydex/tui/components/autocomplete"
)

type GroupsView struct {
	autoComplete *autocomplete.AutoComplete
	components.Container
}

func (gv *GroupsView) HandleEvent(ev tcell.Event) bool {
	switch ev := ev.(type) {
	case *tcell.EventKey:
		if ev.Name() == "Ctrl+D" {
			if App.State.IsReadOnly {
				msg := "Could not delete. Archive in read-only mode."
				App.Notify(msg)
				log.Info(msg)
				return true
			}

			group := App.State.Database.GetFirstGroupByPath(gv.autoComplete.CurrentEntry)
			if group == nil {
				msg := "Could not delete. Group cannot be found."
				App.Notify(msg)
				log.Error(msg, nil)
			}

			name := group.Name

			App.Confirm(
				"Are you sure you want to delete \""+name+"\"?",
				func() {
					err := App.State.Database.RemoveGroup(group.UUID)
					if err != nil {
						msg := "Could not delete. Group cannot be found."
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

					msg := fmt.Sprintf("Group \"%s\" deleted successfully.", name)
					App.Notify(msg)
					log.Info(msg)

					App.NavigateTo(NewGroupListView)
				}, func() {
					msg := "Operation cancelled. Group was not deleted."
					App.Notify(msg)
					log.Info(msg)
				},
			)
		}
	}

	return gv.autoComplete.HandleEvent(ev)
}

func NewGroupListView(screen tcell.Screen) views.Widget {
	App.SetTitle("Select group for \"" + App.State.Entry.GetTitle() + "\"")
	view := &GroupsView{}
	view.Container = *components.NewContainer(screen)
	paths := App.State.Database.GetGroupPaths()
	count := len(paths)
	maxX, maxY := getBoundaries(screen)

	autoCompleteOptions := autocomplete.AutoCompleteOptions{
		Screen:     screen,
		Entries:    paths,
		TotalCount: count,
		MaxX:       maxX,
		MaxY:       maxY,
		OnSelect: func(groupRef string) bool {
			App.State.Group = App.State.Database.GetFirstGroupByPath(groupRef)
			App.State.HasUnsavedChanges = true
			App.NavigateTo(NewEntryView)
			return true
		},
		OnEmpty: func(input string) bool {
			group := App.State.Database.NewGroup(input)
			root := App.State.Database.GetRootGroup()
			root.Groups = append(root.Groups, *group)

			if e := App.State.Database.SaveAndUnlockEntries(); e != nil {
				App.State.IsReadOnly = true
				msg := "Could not save. Switching to read-only to not corrupt data."
				App.Notify(msg)
				log.Error(msg, e)
				return true
			}

			App.State.Group = group
			App.State.HasUnsavedChanges = true
			App.NavigateTo(NewEntryView)

			App.Notify(fmt.Sprintf("Group \"%s\" created successfully.", input))
			return true
		},
		FormatEmptyMessage: func(input string) string {
			return fmt.Sprintf("Create group \"%s\"", input)
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
