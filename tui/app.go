package tui

import (
	"github.com/gdamore/tcell/v2"
	"github.com/gdamore/tcell/v2/views"
	"github.com/shikaan/keydex/pkg/errors"
	"github.com/shikaan/keydex/pkg/kdbx"
	"github.com/shikaan/keydex/pkg/log"
	"github.com/shikaan/keydex/tui/components"
)

type Application struct {
	layout *Layout
	screen tcell.Screen

	LastFocused components.Focusable
	State       State

	lastWidget views.Widget
	isDirty    bool
	isReadOnly bool

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

func (a *Application) SetDirty(value bool) {
	a.isDirty = value
	a.layout.Title.SetDirty(value)
}

func (a *Application) IsDirty() bool {
	return a.isDirty
}

func (a *Application) LockCurrentDatabase(e error) {
	a.isReadOnly = true
	msg := "Could not save. Switching to read-only to preserve database integrity."
	a.Notify(msg)
	log.Error(msg, e)
}

func (a *Application) IsReadOnly() bool {
	return a.isReadOnly
}

func (a *Application) Quit() {
	if !a.IsDirty() {
		a.Application.Quit()
		return
	}

	a.Confirm(
		"Are you sure you want to quit and lose unsaved changes?",
		func() { a.Application.Quit() },
		nil,
	)
}

func (a *Application) CreateEmptyEntry() error {
	entry := a.State.Database.NewEntry()
	a.State.Entry = entry

	if a.State.Group == nil {
		a.State.Group = a.State.Database.GetRootGroup()
	}

	path, err := a.State.Database.MakeEntryEntityPath(entry, a.State.Group)
	if err != nil {
		return err
	}

	a.State.Reference = path
	a.SetDirty(true)
	return nil
}

var App = &Application{}

type State struct {
	Entry     *kdbx.Entry
	Group     *kdbx.Group
	Database  *kdbx.Database
	Reference string
}

func (a *Application) SetScreen(screen tcell.Screen) {
	a.Application.SetScreen(screen)
	a.screen = screen
}

// This is exported only for test purposes
func Setup(screen tcell.Screen, state State, readOnly bool) {
	App.SetScreen(screen)
	App.State = state
	App.layout = NewLayout(screen)
	App.SetRootWidget(App.layout)
	App.isReadOnly = readOnly

	if state.Reference == "" {
		App.NavigateTo(NewHelpView)
	} else {
		App.NavigateTo(NewEntryView)
	}
}

func Run(state State, readOnly bool) error {
	screen, err := tcell.NewScreen()
	if err != nil {
		return errors.MakeError("Unable to start screen", "tui")
	}
	Setup(screen, state, readOnly)
	return App.Run()
}
