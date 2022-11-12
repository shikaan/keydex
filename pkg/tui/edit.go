package tui

import (
	"fmt"

	"github.com/gdamore/tcell/v2"
	"github.com/gdamore/tcell/v2/views"
	"github.com/shikaan/kpcli/pkg/clipboard"
	"github.com/shikaan/kpcli/pkg/kdbx"
	"github.com/shikaan/kpcli/pkg/tui/components"
)

func NewEditView(screen tcell.Screen, e kdbx.Entry) views.Widget {
	view := NewView()

	title := components.NewTitle(e.GetTitle())
	status := components.NewStatus()
	content := newMain(e, status, screen)

	view.SetTitle(title)
	view.SetContent(content)
	view.SetStatus(status)

	return view
}

func newMain(e kdbx.Entry, status *components.Status, screen tcell.Screen) views.Widget {
	form := components.NewForm()

	for i, f := range e.Values {
    isPassword := i == e.GetPasswordIndex()

		if field := newEntryField(f, isPassword, status); field != nil {
			form.AddWidget(field, 0)
		}
	}

	fs := form.Focusables()
	if len(fs) > 0 {
		fs[1].SetFocus(true)
	}

	flex := components.NewContainer(screen)
	flex.SetContent(form)

	return flex
}

func newEntryField(entryValue kdbx.EntryValue, isPassword bool, status *components.Status) *components.Field {
	// Do not print empty fields
	if entryValue.Value.Content == "" {
		return nil
	}

	inputType := components.InputTypeText
	if isPassword {
		inputType = components.InputTypePassword
	}

	fieldOptions := &components.FieldOptions{Label: entryValue.Key, InitialValue: entryValue.Value.Content, InputType: inputType}
	field := components.NewField(fieldOptions)

	field.OnKeyPress(func(ev *tcell.EventKey) bool {
		if ev.Key() == tcell.KeyEnter {
			clipboard.Write(entryValue.Value.Content)
			status.Notify(fmt.Sprintf("Copied \"%s\" to the clipboard", entryValue.Key))
			return true
		}
		return false
	})

	return field
}
