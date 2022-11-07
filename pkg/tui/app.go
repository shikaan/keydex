package tui

import (
	"github.com/gdamore/tcell/v2"
	"github.com/gdamore/tcell/v2/views"
	"github.com/shikaan/kpcli/pkg/errors"
	"github.com/shikaan/kpcli/pkg/kdbx"
)

var App = &views.Application{}

func runView[K interface{}](newView ViewCreator[K], arguments K) error {
	if screen, err := tcell.NewScreen(); err == nil {
		App.SetRootWidget(newView(screen, arguments))
		App.SetScreen(screen)

		return App.Run()
	}

	return errors.MakeError("Unable to start screen", "t")
}

func RunEditView(entry kdbx.Entry) error {
	return runView(NewEditView, entry)
}
