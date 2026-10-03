package field

import (
	"github.com/gdamore/tcell/v2"
	"github.com/gdamore/tcell/v2/views"
	"github.com/shikaan/keydex/tui/components"
)

type Field struct {
	label *views.SimpleStyledText

	components.Focusable
	views.BoxLayout
}

type InputField struct {
	Input *Input
	Field
}

type CheckboxField struct {
	Checkbox *Checkbox
	Field
}

func NewInputField(label string, options *InputOptions) *InputField {
	field := &InputField{}
	field.SetOrientation(views.Horizontal)

	input := NewInput(options)
	labelView := views.NewSimpleStyledText()
	labelView.SetStyle(tcell.StyleDefault.Attributes(tcell.AttrBold))
	labelView.SetText(label + ": ")

	field.AddWidget(labelView, 0)
	field.AddWidget(input, 1)

	field.Input = input
	field.label = labelView
	field.Focusable = input

	return field
}

func NewCheckboxField(label string, options *CheckboxOptions) *CheckboxField {
	field := &CheckboxField{}
	field.SetOrientation(views.Horizontal)

	checkbox := NewCheckbox(options)
	labelView := views.NewSimpleStyledText()
	labelView.SetStyle(tcell.StyleDefault.Attributes(tcell.AttrBold))
	labelView.SetText(label + ": ")

	field.AddWidget(labelView, 0)
	field.AddWidget(checkbox, 1)

	field.Checkbox = checkbox
	field.label = labelView
	field.Focusable = checkbox

	return field
}
