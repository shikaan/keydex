package field

import (
	"github.com/gdamore/tcell/v2"
	"github.com/gdamore/tcell/v2/views"
	"github.com/shikaan/keydex/tui/components"
)

type FieldType int

const (
	FieldTypeText FieldType = iota
	FieldTypePassword
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
	FieldType    FieldType
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

func (f *Field) OnChange(cb func(ev tcell.Event) bool) func() {
	return f.input.OnChange(cb)
}

func (f *Field) GetContent() string {
	return f.input.GetContent()
}

func (f *Field) SetFieldType(t FieldType) {
	f.input.SetHidden(t == FieldTypePassword)
}

func (f *Field) GetFieldType() FieldType {
	if f.input.IsHidden() {
		return FieldTypePassword
	}
	return FieldTypeText
}

func NewField(options *FieldOptions) *Field {
	field := &Field{}
	field.SetOrientation(views.Horizontal)

	hidden := options.FieldType == FieldTypePassword

	opts := &InputOptions{
		InitialValue: options.InitialValue,
		Hidden:       hidden,
		Disabled:     options.Disabled,
	}
	input := NewInput(opts)
	input.SetContent(options.InitialValue)

	label := views.NewSimpleStyledText()
	label.SetStyle(tcell.StyleDefault.Attributes(tcell.AttrBold))
	label.SetText(options.Label + ": ")

	field.AddWidget(label, 0)
	field.AddWidget(input, 1)

	field.input = input
	field.label = label

	return field
}
