package tui

import (
	"fmt"

	"github.com/gdamore/tcell/v2"
	"github.com/gdamore/tcell/v2/views"

	"github.com/shikaan/kpcli/pkg/clipboard"
	"github.com/shikaan/kpcli/pkg/kdbx"
	"github.com/shikaan/kpcli/pkg/tui/components"
)

type fieldKey = string
type fieldMap = map[fieldKey]components.Field

type EditView struct {
	// Model
	entry    kdbx.Entry
	database kdbx.Database
	reference      string

	// View
	fieldByKey fieldMap

	status *components.Status
	form   views.Widget
	title  views.Widget

	View
}

func (v *EditView) HandleEvent(ev tcell.Event) bool {
	switch ev := ev.(type) {
	case *tcell.EventKey:
		if ev.Name() == "Ctrl+O" {
			entry := v.database.GetEntry(v.reference)

			for i, vd := range entry.Values {
				if field, ok := v.fieldByKey[vd.Key]; ok {
					entry.Values[i].Value.Content = field.GetContent()
        }
			}

			if e := v.database.Save(); e != nil {
				// TODO: logging
				v.status.Notify("Could not save. See logs for error.")
			}

			v.status.Notify(fmt.Sprintf("Entry \"%s\" saved succesfully", entry.GetTitle()))
			return true
		}
	}

	return v.View.HandleEvent(ev)
}

type EditViewProps struct {
	Entry     kdbx.Entry
	Database  kdbx.Database
	Reference string
}

func NewEditView(screen tcell.Screen, props EditViewProps) views.Widget {
	view := &EditView{}
	view.View = *NewView()
	view.entry = props.Entry
	view.database = props.Database
	view.reference = props.Reference

	title := components.NewTitle(props.Entry.GetTitle())
	status := components.NewStatus()
	form, fieldMap := view.newForm(screen, props)

	view.fieldByKey = fieldMap

	view.SetStatus(status)
	view.status = status

	view.SetTitle(title)
  view.title = title

	view.SetContent(form)
  view.form = form

	return view
}

func (view *EditView) newForm(screen tcell.Screen, props EditViewProps) (views.Widget, fieldMap) {
	form := components.NewForm()
	fields := fieldMap{}

	for _, f := range props.Entry.Values {
		if field := view.newEntryField(f.Key, f.Value.Content, f.Value.Protected.Bool); field != nil {
			form.AddWidget(field, 0)
			// Using f.Value as binding key (for example, is we just used props.reference)
			// would cause the title field to be unmodifiable, because the reference
			// which is based on the title would change
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

func (view *EditView) newEntryField(label, initialValue string, isProtected bool) *components.Field {
	// Do not print empty fields
	if initialValue == "" {
		return nil
	}

  inputType := components.InputTypeText
	if isProtected {
			inputType = components.InputTypePassword
	}

	fieldOptions := &components.FieldOptions{Label: label, InitialValue: initialValue, InputType: inputType}
	field := components.NewField(fieldOptions)

	field.OnKeyPress(func(ev *tcell.EventKey) bool {
		if ev.Name() == "Ctrl+C" {
			clipboard.Write(field.GetContent())
			view.status.Notify(fmt.Sprintf("Copied \"%s\" to the clipboard", label))
			return true
		}

    if ev.Name() == "Ctrl+H" || ev.Name() == "Ctrl+J" {
      if isProtected {
        if field.GetInputType() == components.InputTypePassword {
          field.SetInputType(components.InputTypeText)
        } else {
          field.SetInputType(components.InputTypePassword)
        }
      }

      return true
    }

		return false
	})

	return field
}
