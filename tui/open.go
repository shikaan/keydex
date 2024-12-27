package tui

import (
	"fmt"

	"github.com/gdamore/tcell/v2"
	"github.com/gdamore/tcell/v2/views"

	"github.com/shikaan/keydex/pkg/clipboard"
	"github.com/shikaan/keydex/pkg/kdbx"
	"github.com/shikaan/keydex/pkg/log"
	"github.com/shikaan/keydex/tui/components"
	"github.com/shikaan/keydex/tui/components/field"
)

type fieldKey = string
type fieldMap = map[fieldKey]*field.Field

type HomeView struct {
	fieldByKey fieldMap
	form       *components.Form
	components.Container
}

var IsReadOnly = false

func (v *HomeView) HandleEvent(ev tcell.Event) bool {
	switch ev := ev.(type) {
	case *tcell.EventKey:
		if ev.Name() == "Ctrl+O" {
			if IsReadOnly {
				App.Notify("Could not save: archive in read-only mode.")
				log.Info("Could not save: archive in read-only mode.")
				return true
			}

			uuid := App.State.Entry.UUID
			entry := App.State.Database.GetEntry(uuid)

			if entry == nil {
				App.Notify("Could not find entry at " + App.State.Reference + ".")
				log.Info("Could not find entry at " + App.State.Reference + ".")
				return true
			}

			App.Confirm(
				"Are you sure?",
				func() {
					for i, vd := range entry.Values {
						if field, ok := v.fieldByKey[vd.Key]; ok {
							entry.Values[i].Value.Content = string(field.GetContent())
						}
					}

					if e := App.State.Database.Save(); e != nil {
						log.Error("Could not save. See logs for error.", e)
						App.Notify("Could not save. See logs for error.")
						return
					}

					// Unlocking again to allow further modifications
					if e := App.State.Database.UnlockProtectedEntries(); e != nil {
						IsReadOnly = true
						App.Notify("Could not save. Switching to read-only to not corrupt data.")
						log.Error("Could not save. Switching to read-only to not corrupt data.", e)
						return
					}

					App.Notify(fmt.Sprintf("Entry \"%s\" saved succesfully.", entry.GetTitle()))
					App.State.HasUnsavedChanges = false
				}, func() {
					App.Notify("Operation cancelled. Entry was not saved.")
					log.Info("Operation cancelled. Entry was not saved.")
				},
			)
		}
	}

	return v.Container.HandleEvent(ev)
}

func NewHomeView(screen tcell.Screen) views.Widget {
	App.SetTitle("\"" + App.State.Entry.GetTitle() + "\"")
	view := &HomeView{}
	view.Container = *components.NewContainer(screen)

	form, fieldMap := view.newForm(screen, App.State)
	view.fieldByKey = fieldMap

	view.SetContent(form)
	view.form = form

	return view
}

func (view *HomeView) newForm(_ tcell.Screen, props State) (*components.Form, fieldMap) {
	form := components.NewForm()
	fields := fieldMap{}

	for _, f := range props.Entry.Values {
		if field := view.newEntryField(f.Key, f.Value.Content, f.Value.Protected.Bool); field != nil {
			form.AddWidget(field, 0)
			// Using f.Value as binding key (for example, is we just used props.reference)
			// would cause the title field to be unmodifiable, because the reference
			// which is based on the title would change
			fields[f.Key] = field
		}
	}

	fs := form.Focusables()
	if len(fs) > 0 {
		fs[0].SetFocus(true)
	}

	return form, fields
}

func (view *HomeView) newEntryField(label, initialValue string, isProtected bool) *field.Field {
	// Do not print empty fields, unless they are the title
	if initialValue == "" && label != kdbx.TITLE_KEY {
		return nil
	}

	inputType := field.InputTypeText
	if isProtected {
		inputType = field.InputTypePassword
	}

	fieldOptions := &field.FieldOptions{Label: label, InitialValue: initialValue, InputType: inputType}
	f := field.NewField(fieldOptions)

	f.OnFocus(func() bool {
		App.LastFocused = f
		return true
	})

	f.OnKeyPress(func(ev *tcell.EventKey) bool {
		App.State.HasUnsavedChanges = true

		if ev.Name() == "Ctrl+C" {
			clipboard.Write(string(f.GetContent()))
			App.Notify(fmt.Sprintf("Copied \"%s\" to the clipboard.", label))
			return true
		}

		if ev.Name() == "Ctrl+R" {
			if isProtected {
				if f.GetInputType() == field.InputTypePassword {
					f.SetInputType(field.InputTypeText)
				} else {
					f.SetInputType(field.InputTypePassword)
				}
			}

			return true
		}

		if ev.Key() == tcell.KeyRune {
			if isProtected && f.GetInputType() == field.InputTypePassword {
				App.Notify("Reveal (^R) the field to edit.")
			}
		}

		return false
	})

	return f
}
