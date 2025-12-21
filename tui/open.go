package tui

import (
	"fmt"
	"time"

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
	fieldByKey   fieldMap
	initialGroup *kdbx.Group
	form         *components.Form
	components.Container
}

func (v *HomeView) updateEntry(entry *kdbx.Entry) {
	for key, field := range v.fieldByKey {
		entry.SetValue(key, field.GetContent())
	}

	entry.SetLastUpdated()
	App.State.Database.AddEntryToGroup(entry, App.State.Group)
}

func (v *HomeView) writeDatabase() {
	if e := App.State.Database.Save(); e != nil {
		msg := "Could not save. See logs for error."
		App.Notify(msg)
		log.Error(msg, e)
		return
	}

	// Unlocking again to allow further modifications
	if e := App.State.Database.UnlockProtectedEntries(); e != nil {
		App.State.IsReadOnly = true
		msg := "Could not save. Switching to read-only to not corrupt data."
		App.Notify(msg)
		log.Error(msg, e)
		return
	}
}

func (v *HomeView) HandleEvent(ev tcell.Event) bool {
	switch ev := ev.(type) {
	case *tcell.EventKey:
		if ev.Name() == "Ctrl+K" {
			App.NavigateTo(NewGroupsView)
			return true
		}

		if ev.Name() == "Ctrl+O" {
			if App.State.IsReadOnly {
				msg := "Could not save. Archive in read-only mode."
				App.Notify(msg)
				log.Info(msg)
				return true
			}

			uuid := App.State.Entry.UUID
			existingEntry := App.State.Database.GetEntry(uuid)

			if existingEntry == nil {
				App.Confirm(
					"Do you want to create \""+App.State.Entry.GetTitle()+"\"?",
					func() {
						v.updateEntry(App.State.Entry)
						v.writeDatabase()
						msg := fmt.Sprintf("Entry \"%s\" created successfully.", App.State.Entry.GetTitle())
						App.Notify(msg)
						log.Info(msg)
						App.State.HasUnsavedChanges = false
					}, func() {
						msg := "Operation cancelled. Entry was not created."
						App.Notify(msg)
						log.Info(msg)
					})
				return true
			}

			App.Confirm(
				"Are you sure?",
				func() {
					v.updateEntry(existingEntry)
					v.writeDatabase()
					msg := fmt.Sprintf("Entry \"%s\" saved successfully.", App.State.Entry.GetTitle())
					App.Notify(msg)
					log.Info(msg)
					App.State.HasUnsavedChanges = false
				}, func() {
					msg := "Operation cancelled. Entry was not saved."
					App.Notify(msg)
					log.Info(msg)
				},
			)
		}

		if ev.Name() == "Ctrl+D" {
			if App.State.IsReadOnly {
				msg := "Could not delete. Archive in read-only mode."
				App.Notify(msg)
				log.Info(msg)
				return true
			}

			uuid := App.State.Entry.UUID
			existingEntry := App.State.Database.GetEntry(uuid)

			if existingEntry == nil {
				msg := "Could not find entry at " + App.State.Reference + "."
				App.Notify(msg)
				log.Info(msg)
				return true
			}

			App.Confirm(
				"Are you sure you want to delete \""+App.State.Entry.GetTitle()+"\"?",
				func() {
					title := App.State.Entry.GetTitle()

					err := App.State.Database.RemoveEntry(App.State.Entry.UUID)
					if err != nil {
						msg := "Could not delete. Entry cannot be found."
						App.Notify(msg)
						log.Error(msg, err)
						return
					}

					v.writeDatabase()

					msg := fmt.Sprintf("Entry \"%s\" deleted successfully.", title)
					App.Notify(msg)
					log.Info(msg)
					App.State.HasUnsavedChanges = false

					App.NavigateTo(NewListView)
				}, func() {
					msg := "Operation cancelled. Entry was not deleted."
					App.Notify(msg)
					log.Info(msg)
				},
			)
		}
	}

	return v.Container.HandleEvent(ev)
}

func NewHomeView(screen tcell.Screen) views.Widget {
	title := "\"" + App.State.Entry.GetTitle() + "\""
	if App.State.IsReadOnly {
		title += " [READ ONLY]"
	}
	App.SetTitle(title)
	view := &HomeView{}
	view.Container = *components.NewContainer(screen)

	form, fieldMap := view.newForm(screen, App.State.Entry, App.State.Group)
	view.fieldByKey = fieldMap

	view.SetContent(form)
	view.form = form

	return view
}

func (view *HomeView) newForm(_ tcell.Screen, entry *kdbx.Entry, group *kdbx.Group) (*components.Form, fieldMap) {
	form := components.NewForm()
	fields := fieldMap{}

	for _, f := range entry.Values {
		if field := view.newEntryField(f.Key, f.Value.Content, f.Value.Protected.Bool); field != nil {
			form.AddWidget(field, 0)
			// Using f.Value as binding key (for example, is we just used props.reference)
			// would cause the title field to be unmodifiable, because the reference
			// which is based on the title would change
			fields[f.Key] = field
		}
	}

	spacer := &views.Spacer{}
	form.AddWidget(spacer, 1)

	// The space is just for alignment
	groupField := view.newMetaField("Group", "  "+group.Name)
	form.AddWidget(groupField, 0)

	createdAt := entry.Times.CreationTime.Time.Format(time.DateTime)
	created := view.newMetaField("Created", createdAt)
	form.AddWidget(created, 0)

	updatedAt := entry.Times.LastModificationTime.Time.Format(time.DateTime)
	updated := view.newMetaField("Updated", updatedAt)
	form.AddWidget(updated, 0)

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

	fieldOptions := &field.FieldOptions{Label: label, InitialValue: initialValue, InputType: inputType, Disabled: App.State.IsReadOnly}
	f := field.NewField(fieldOptions)

	f.OnFocus(func() bool {
		App.LastFocused = f
		return true
	})

	f.OnChange(func(ev tcell.Event) bool {
		App.State.HasUnsavedChanges = true
		return false
	})

	f.OnKeyPress(func(ev *tcell.EventKey) bool {
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

func (view *HomeView) newMetaField(label, value string) *views.Text {
	field := views.NewText()
	field.SetStyle(tcell.StyleDefault.Normal().Dim(true))
	field.SetText(label + ": " + value)
	return field
}
