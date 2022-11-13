package tui

import (
	"github.com/gdamore/tcell/v2"
	"github.com/gdamore/tcell/v2/views"
	"github.com/shikaan/kpcli/pkg/errors"
)

var App = &views.Application{}

func runView[K interface{}](newView ViewCreator[K], arguments K) error {
	if screen, err := tcell.NewScreen(); err == nil {
		App.SetRootWidget(newView(screen, arguments))
		App.SetScreen(screen)

		return App.Run()
	}

	return errors.MakeError("Unable to start screen", "tui")
}

func RunEditView(props EditViewProps) error {
	return runView(NewEditView, props)
}
