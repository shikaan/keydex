package field

import (
	"github.com/gdamore/tcell/v2"
	"github.com/gdamore/tcell/v2/views"
	"github.com/shikaan/keydex/tui/components"
)

type Field struct {
	input *Input
	label *views.SimpleStyledText

	components.Focusable
	views.BoxLayout
}

type FieldOptions struct {
	Label        string
	InitialValue string
	InputType    InputType
	Disabled     bool
}

func (f *Field) HasFocus() bool {
	return f.input.HasFocus()
}

func (f *Field) SetFocus(on bool) {
	f.input.SetFocus(on)
}

func (f *Field) OnFocus(cb func() bool) func() {
	return f.input.OnFocus(cb)
}

func (f *Field) HandleEvent(ev tcell.Event) bool {
	if !f.HasFocus() {
		return false
	}

	return f.input.HandleEvent(ev)
}

func (f *Field) OnKeyPress(cb func(ev *tcell.EventKey) bool) func() {
	return f.input.OnKeyPress(cb)
}

func (f* Field) OnChange(cb func(ev tcell.Event) bool) func() {
	return f.input.OnChange(cb)
}

func (f *Field) GetContent() string {
	return f.input.GetContent()
}

func (f *Field) SetInputType(t InputType) {
	f.input.SetInputType(t)
}

func (f *Field) GetInputType() InputType {
	return f.input.GetInputType()
}

func NewField(options *FieldOptions) *Field {
	// TODO: we can maybe add some padding by directly accessing the model and tampering wiht GetBounds
	field := &Field{}
	field.SetOrientation(views.Horizontal)

	opts := &InputOptions{InitialValue: options.InitialValue, Type: options.InputType, Disabled: options.Disabled}
	input := NewInput(opts)
	input.SetContent(options.InitialValue)
	input.SetInputType(options.InputType)

	label := views.NewSimpleStyledText()
	label.SetStyle(tcell.StyleDefault.Attributes(tcell.AttrBold))
	label.SetText(options.Label + ": ")

	field.AddWidget(label, 0)
	field.AddWidget(input, 1)

	field.input = input
	field.label = label

	return field
}
