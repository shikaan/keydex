package pages

import (
	"fmt"
	"syscall"

	"golang.org/x/term"

	"github.com/rivo/tview"
	"github.com/shikaan/kpcli/pages/explore"
	"github.com/shikaan/kpcli/pkg/kdbx"
)

func Run(database, keyPath string) error {
	fmt.Println("Insert password: ")

	pwd, err := term.ReadPassword(int(syscall.Stdin))
	handleError(err)

	kdbx, err := kdbx.New(database)
	handleError(err)

	err = kdbx.Unlock(string(pwd))
	handleError(err)

	app := tview.NewApplication()
	router := tview.NewPages()
	modal := tview.NewModal()

	openDialog := func(msg string) { openDialog(app, router, modal, msg) }
	setFocus := func(p tview.Primitive) { app.SetFocus(p) }

	router.
		AddPage("explore", explore.Render(*kdbx, openDialog, setFocus), true, true).
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

func handleError(e error) {
	if e != nil {
		panic(e)
	}
}
