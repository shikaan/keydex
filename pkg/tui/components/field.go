package components

import (
	"github.com/gdamore/tcell/v2"
	"github.com/gdamore/tcell/v2/views"
)

type Field struct {
  input *Input
  label *views.SimpleStyledText

  Focusable
  views.BoxLayout
}

type FieldOptions struct {
  Label string
  InitialValue string
  InputType InputType 
}

func (f *Field) HasFocus() bool {
  return f.input.HasFocus()
}

func (f *Field) SetFocus(on bool) {
  f.input.SetFocus(on)
}

func (f *Field) HandleEvent(ev tcell.Event) bool {
  if !f.HasFocus() {
    return false
  }

  return f.input.HandleEvent(ev)
}

func (f *Field) OnKeyPress (cb func (ev *tcell.EventKey) bool) func () {
  return f.input.OnKeyPress(cb)
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
  field := &Field{ }
  field.SetOrientation(views.Horizontal)

  opts := &InputOptions{ InitialValue: options.InitialValue, Type: options.InputType }
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
