package tui

import (
	"fmt"

	"github.com/gdamore/tcell/v2"
	"github.com/gdamore/tcell/v2/views"
	"github.com/shikaan/keydex/pkg/kdbx"
	"github.com/shikaan/keydex/pkg/log"
	"github.com/shikaan/keydex/tui/components"
	"github.com/shikaan/keydex/tui/components/field"
)

type FieldView struct {
	label     *field.InputField
	protected *field.CheckboxField
	form      *components.Form
	components.Container
}

func (v *FieldView) HandleEvent(ev tcell.Event) bool {
	switch ev := ev.(type) {
	case *tcell.EventKey:
		if ev.Name() == "Ctrl+O" {
			if App.IsReadOnly() {
				msg := "Cannot save. Archive in read-only mode."
				App.Notify(msg)
				log.Info(msg)
				return true
			}

			label := v.label.Input.GetContent()
			if kdbx.IsStandardField(label) {
				msg := fmt.Sprintf("Cannot save. Label \"%s\" is reserved for standard fields.", label)
				App.Notify(msg)
				log.Info(msg)
				return true
			}

			entry, entryField := App.State.Entry, App.State.EntryField
			if existing := entry.Get(label); existing != nil && existing != entryField {
				msg := fmt.Sprintf("Cannot save. Label \"%s\" is already in use.", label)
				App.Notify(msg)
				log.Info(msg)
				return true
			}

			App.Confirm(
				"Save changes? This will overwrite the existing file.",
				func() {
					entryField.Key = label
					entryField.Value.Protected.Bool = v.protected.Checkbox.GetContent()

					if !entry.HasField(entryField) {
						entry.Values = append(entry.Values, *entryField)
					}

					entry.SetLastUpdated()
					App.SaveEntry()

					if e := App.State.Database.SaveAndUnlockEntries(); e != nil {
						App.LockCurrentDatabase(e)
						App.NavigateTo(NewEntryView)
						return
					}

					msg := fmt.Sprintf("Field \"%s\" saved successfully.", label)
					App.Notify(msg)
					log.Info(msg)
					App.SetDirty(false)
					App.NavigateTo(NewEntryView)
				}, func() {
					msg := "Operation cancelled. Field was not saved."
					App.Notify(msg)
					log.Info(msg)
					App.RefreshCurrentView()
				},
			)
			return true
		}
	}

	return v.Container.HandleEvent(ev)
}

func NewFieldView(screen tcell.Screen) views.Widget {
	if App.State.Entry == nil {
		panic("missing entry")
	}

	if App.State.Group == nil {
		panic("missing group")
	}

	if App.State.EntryField == nil {
		panic("missing field")
	}

	if !App.State.Entry.HasField(App.State.EntryField) {
		App.SetDirty(true)
	}

	App.SetTitle(App.State.EntryField.Key)
	view := &FieldView{}
	view.Container = components.Container{}

	view.form = components.NewForm()

	labelOpts := &field.InputOptions{
		InitialValue: App.State.EntryField.Key,
		Disabled:     false,
	}
	view.label = field.NewInputField("Label", labelOpts)

	view.label.Input.OnChange(func(ev tcell.Event) bool {
		App.SetDirty(true)
		return false
	})

	protectedOpts := &field.CheckboxOptions{
		InitialValue: App.State.EntryField.Value.Protected.Bool,
		Disabled:     false,
	}
	view.protected = field.NewCheckboxField("Protected", protectedOpts)

	view.protected.Checkbox.OnChange(func(ev tcell.Event) bool {
		App.SetDirty(true)
		return false
	})

	view.form.AddWidget(view.label, 1)
	view.form.AddWidget(view.protected, 1)

	fs := view.form.Focusables()
	if len(fs) > 0 {
		fs[0].SetFocus(true)
	}
	view.SetContent(view.form)

	return view
}
