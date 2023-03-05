package tui

import (
	"math"

	"github.com/gdamore/tcell/v2"
	"github.com/gdamore/tcell/v2/views"

	"github.com/shikaan/keydex/tui/components"
)

type ListView struct {
	components.Container
}

func NewListView(screen tcell.Screen) views.Widget {
	App.SetTitle("Search")
	view := &ListView{}
	view.Container = *components.NewContainer(screen)
	paths := App.State.Database.GetEntryPaths()
	count := len(paths)
	maxX, maxY := getBoundaries(screen)

	autoCompleteOptions := components.AutoCompleteOptions{
		Screen:     screen,
		Entries:    paths,
		TotalCount: count,
		MaxX:       maxX,
		MaxY:       maxY,
		OnSelect: func(ref string) bool {
			App.State.Reference = ref
			App.State.Entry = App.State.Database.GetFirstEntryByPath(ref)
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
	view.SetContent(autoComplete)

	return view
}

func getBoundaries(screen tcell.Screen) (int, int) {
	x, y := screen.Size()

	// one third of the screen width
	// all the height - title, status, search, counter, notification,
	// and 4 more lines of buffer just in case
	return int(math.Max(float64(x/3), components.MIN_WIDTH)), y - 10
}
