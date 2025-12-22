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
	components.Container
}

func NewGroupsView(screen tcell.Screen) views.Widget {
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
			App.NavigateTo(NewHomeView)
			return true
		},
		OnEmpty: func(input string) bool {
			group := App.State.Database.NewGroup(input)
			root := App.State.Database.GetRootGroup()
			root.Groups = append(root.Groups, *group)

			if e := App.State.Database.Save(); e != nil {
				msg := "Could not save. See logs for error."
				App.Notify(msg)
				log.Error(msg, e)
				return true
			}

			// Unlocking again to allow further modifications
			if e := App.State.Database.UnlockProtectedEntries(); e != nil {
				App.State.IsReadOnly = true
				msg := "Could not save. Switching to read-only to not corrupt data."
				App.Notify(msg)
				log.Error(msg, e)
				return true
			}

			App.State.Group = group
			App.State.HasUnsavedChanges = true
			App.NavigateTo(NewHomeView)

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

	return view
}
