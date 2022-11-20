package tui

import (
	"fmt"

	"github.com/gdamore/tcell/v2"
	"github.com/gdamore/tcell/v2/views"

	"github.com/shikaan/kpcli/pkg/clipboard"
	"github.com/shikaan/kpcli/pkg/tui/components"
)

type fieldKey = string
type fieldMap = map[fieldKey]components.Field

type EditView struct {
	fieldByKey fieldMap
	form views.Widget
	components.Container
}

func (v *EditView) HandleEvent(ev tcell.Event) bool {
	switch ev := ev.(type) {
	case *tcell.EventKey:
		if ev.Name() == "Ctrl+O" {
      uuid := App.State.Entry.UUID
			entry := App.State.Database.GetEntry(uuid)

      if entry == nil {
        App.Notify("Could not find entry at " + App.State.Reference)
        return false
      }

			for i, vd := range entry.Values {
				if field, ok := v.fieldByKey[vd.Key]; ok {
					entry.Values[i].Value.Content = field.GetContent()
				}
			}

			if e := App.State.Database.Save(); e != nil {
				// TODO: logging
				App.Notify("Could not save. See logs for error.")
			  return false
      }

			App.Notify(fmt.Sprintf("Entry \"%s\" saved succesfully", entry.GetTitle()))
			return true
		}
	}

	return v.Container.HandleEvent(ev)
}

func NewEditView(screen tcell.Screen, state State) views.Widget {
	view := &EditView{}
	view.Container = *components.NewContainer(screen)

	form, fieldMap := view.newForm(screen, state)
	view.fieldByKey = fieldMap

	view.SetContent(form)
	view.form = form

	return view
}

func (view *EditView) newForm(screen tcell.Screen, props State) (views.Widget, fieldMap) {
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

	return form, fields
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
			App.Notify(fmt.Sprintf("Copied \"%s\" to the clipboard", label))
			return true
		}

		if ev.Name() == "Ctrl+R" {
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
