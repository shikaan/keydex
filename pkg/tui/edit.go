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

	title := components.NewTitle(e.Name)
	status := components.NewStatus()
	content := newMain(e, status, screen)

	view.SetTitle(title)
	view.SetContent(content)
	view.SetStatus(status)

	return view
}

func newMain(e kdbx.Entry, status *components.Status, screen tcell.Screen) views.Widget {
	form := components.NewForm()

	for _, f := range e.Fields {
		if field := newEntryField(f, status); field != nil {
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

func newEntryField(entryField kdbx.Field, status *components.Status) *components.Field {
	key := entryField[0]
	value := entryField[1]

	// Do not print empty fields
	if value == "" {
		return nil
	}

	inputType := components.InputTypeText
	if key == kdbx.PASSWORD_KEY {
		inputType = components.InputTypePassword
	}

	fieldOptions := &components.FieldOptions{Label: key, InitialValue: value, InputType: inputType}
	field := components.NewField(fieldOptions)

	field.OnKeyPress(func(ev *tcell.EventKey) bool {
		if ev.Key() == tcell.KeyEnter {
			clipboard.Write(value)
			status.Notify(fmt.Sprintf("Copied \"%s\" to the clipboard", key))
			return true
		}
		return false
	})

	return field
}
