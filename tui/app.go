package tui

import (
	"github.com/gdamore/tcell/v2"
	"github.com/gdamore/tcell/v2/views"
	"github.com/shikaan/kpcli/pkg/errors"
	"github.com/shikaan/kpcli/pkg/kdbx"
	"github.com/shikaan/kpcli/tui/components"
)

type Application struct {
	layout *Layout
	screen tcell.Screen
	State  State

	views.Application
}

func (a *Application) NavigateTo(newView func(tcell.Screen, State) views.Widget) {
	a.layout.SetContent(newView(a.screen, a.State))
}

func (a *Application) Notify(msg string) {
	a.layout.Status.Notify(msg)
}

func (a *Application) SetTitle(title string) {
	a.layout.SetTitle(components.NewTitle(title))
}

var App = &Application{}

type State struct {
	Entry     *kdbx.Entry
	Database  *kdbx.Database
	Reference string
}

func Run(state State) error {
	if screen, err := tcell.NewScreen(); err == nil {
		App.SetScreen(screen)
    App.screen = screen
    App.State = state
		App.layout = NewLayout(screen)
		App.SetRootWidget(App.layout)

		App.NavigateTo(NewEditView)

		return App.Run()
	}

	return errors.MakeError("Unable to start screen", "tui")
}
