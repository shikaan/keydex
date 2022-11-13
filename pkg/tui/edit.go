package tui

import (
	"fmt"

	"github.com/gdamore/tcell/v2"
	"github.com/gdamore/tcell/v2/views"
	"github.com/shikaan/kpcli/pkg/clipboard"
	"github.com/shikaan/kpcli/pkg/kdbx"
	"github.com/shikaan/kpcli/pkg/tui/components"
)

type EditView struct {
	entry    kdbx.Entry
	database kdbx.Database
  ref string

  fields map[string]components.Field
	status *components.Status

	View
}

func (v *EditView) HandleEvent(ev tcell.Event) bool {
	switch ev := ev.(type) {
	case *tcell.EventKey:
		if ev.Name() == "Ctrl+O" {
      entry := v.database.GetEntry(v.ref)
		
      for i, vd := range entry.Values {
        if field, ok := v.fields[vd.Key]; ok {
          entry.Values[i].Value.Content = field.GetContent()
        }
      }

			e := v.database.Save()
      if e != nil {
        panic(e.Error())
			}

			v.status.Notify("Entry saved!")
			return true
		}
	}

	return v.View.HandleEvent(ev)
}

type EditViewProps struct {
	Entry    kdbx.Entry
	Database kdbx.Database
  Reference string
}

func NewEditView(screen tcell.Screen, props EditViewProps) views.Widget {
	title := components.NewTitle(props.Entry.GetTitle())
	status := components.NewStatus()
	content, fields := newMain(props.Entry, props.Database, props.Reference, status, screen)

	view := &EditView{props.Entry, props.Database, props.Reference, fields, status, *NewView()}

	view.SetTitle(title)
	view.SetContent(content)
	view.SetStatus(status)

	return view
}

func newMain(e kdbx.Entry, db kdbx.Database, ref string, status *components.Status, screen tcell.Screen) (views.Widget, map[string]components.Field) {
	form := components.NewForm()
  fields := map[string]components.Field{}

	for i, f := range e.Values {
		isPassword := i == e.GetPasswordIndex()

		if field := newEntryField(f, db, e, ref, isPassword, status); field != nil {
			form.AddWidget(field, 0)
      // Use key as binding value instead of value, so Title can change
      fields[f.Key] = *field
		}
	}

	fs := form.Focusables()
	if len(fs) > 0 {
		fs[0].SetFocus(true)
	}

	flex := components.NewContainer(screen)
	flex.SetContent(form)

	return flex, fields
}

func newEntryField(entryValue kdbx.EntryValue, d kdbx.Database, e kdbx.Entry, ref string, isPassword bool, status *components.Status) *components.Field {
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
		if ev.Name() == "Ctrl+C" {
			clipboard.Write(entryValue.Value.Content)
			status.Notify(fmt.Sprintf("Copied \"%s\" to the clipboard", entryValue.Key))
			return true
		}
		return false
	})

	return field
}
