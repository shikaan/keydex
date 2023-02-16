package tui

import (
	"github.com/gdamore/tcell/v2"
	"github.com/gdamore/tcell/v2/views"
	"github.com/shikaan/keydex/pkg/errors"
	"github.com/shikaan/keydex/pkg/kdbx"
	"github.com/shikaan/keydex/tui/components"
)

type Application struct {
	layout *Layout
	screen tcell.Screen

	LastFocused components.Focusable
	State       State

	views.Application
}

func (a *Application) NavigateTo(newView func(tcell.Screen) views.Widget) {
	a.layout.SetContent(newView(a.screen))
}

func (a *Application) Notify(msg string) {
	a.layout.Status.Notify(msg)
}

func (a *Application) Confirm(msg string, onAccept func(), onReject func()) {
	if a.LastFocused != nil {
		a.LastFocused.SetFocus(false)
	}

	a.layout.Status.Confirm(
		msg,
		func() {
			if onAccept != nil {
				onAccept()
			}

			if a.LastFocused != nil {
				a.LastFocused.SetFocus(true)
			}
		},
		func() {
			if onReject != nil {
				onReject()
			}

			if a.LastFocused != nil {
				a.LastFocused.SetFocus(true)
			}
		},
	)
}

func (a *Application) SetTitle(title string) {
	a.layout.Title.SetTitle(title)
}

func (a *Application) Quit() {
	if !a.State.HasUnsavedChanges {
		a.Application.Quit()
		return
	}

	a.Confirm("Are you sure you want to quit and lose unsaved changes?", func() { a.Application.Quit() }, nil)
}

var App = &Application{}

type State struct {
	Entry             *kdbx.Entry
	Database          *kdbx.Database
	Reference         string
	HasUnsavedChanges bool
}

func Run(state State) error {
	if screen, err := tcell.NewScreen(); err == nil {
		App.SetScreen(screen)
		App.screen = screen
		App.State = state
		App.layout = NewLayout(screen)
		App.SetRootWidget(App.layout)

		if state.Reference == "" {
			App.NavigateTo(NewHelpView)
		} else {
			App.NavigateTo(NewHomeView)
		}

		return App.Run()
	}

	return errors.MakeError("Unable to start screen", "tui")
}
