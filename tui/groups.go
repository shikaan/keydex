package tui

import (
	"github.com/gdamore/tcell/v2"
	"github.com/gdamore/tcell/v2/views"

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
