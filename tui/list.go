package tui

import (
	"github.com/gdamore/tcell/v2"
	"github.com/gdamore/tcell/v2/views"

	"github.com/shikaan/kpcli/pkg/utils"
	"github.com/shikaan/kpcli/tui/components"
)

type ListView struct {
	views.Panel
}

func NewListView(screen tcell.Screen, state State) views.Widget {
	container := components.NewContainer(screen)
	view := &ListView{}
	paths := state.Database.GetEntryPaths()
	maxX, maxY := getBoundaries(screen)

	autoCompleteOptions := components.AutoCompleteOptions{
		Screen:     screen,
		Entries:    paths,
		TotalCount: len(paths),
		MaxX:       maxX,
		MaxY:       maxY,
		OnSelect: func(entry string) bool {
			App.State.Reference = entry
			App.State.Entry = App.State.Database.GetFirstEntryByPath(entry)
			App.NavigateTo(NewHomeView)
			return true
		},
	}

	autoComplete := components.NewAutoComplete(autoCompleteOptions)
	autoComplete.OnFocus(func() bool {
			App.LastFocused = autoComplete
			return true
		})
  autoComplete.SetFocus(true)
  view.AddWidget(autoComplete, 1)

	container.SetContent(view)
	return container
}

func getBoundaries(screen tcell.Screen) (int, int) {
	x, y := (screen).Size()

	// one third of the screen width
	// all the height - title, status, search, counter, and notification
	return utils.Max(x/3, components.MIN_WIDTH), y - 6
}
