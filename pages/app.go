package app

import (
	"fmt"
	"syscall"

	"golang.org/x/term"

	"github.com/rivo/tview"
	"github.com/shikaan/kpcli/pages/explore"
	"github.com/shikaan/kpcli/pkg/kdbx"
	"github.com/shikaan/kpcli/pkg/logger"
)

func Run(database, keyPath string, logger *logger.Logger) error {
	fmt.Println("Insert password: ")

	pwd, err := term.ReadPassword(int(syscall.Stdin))
	if err != nil {
		return err
	}

	kdbx, err := kdbx.New(database)
	if err != nil {
		return err
	}

	err = kdbx.Unlock(string(pwd))
	if err != nil {
		return err
	}

	app := tview.NewApplication()
	router := tview.NewPages()
	modal := tview.NewModal()

	openDialog := func(msg string) { openDialog(app, router, modal, msg) }
	setFocus := func(p tview.Primitive) { app.SetFocus(p) }

	router.
		AddPage("explore", explore.Render(*kdbx, openDialog, setFocus, logger), true, true).
		AddPage("ui", modal, true, false)

	return app.SetRoot(router, true).Run()
}

func openDialog(app *tview.Application, router *tview.Pages, modal *tview.Modal, message string) {
	modal.SetText(message)
	modal.ClearButtons().AddButtons([]string{"ok"})

	// What if we have nested modals?
	lastFocused := app.GetFocus()
	router.ShowPage("ui")

	modal.SetDoneFunc(func(buttonIndex int, buttonLabel string) {
		router.HidePage("ui")
		app.SetFocus(lastFocused)
	})
}
