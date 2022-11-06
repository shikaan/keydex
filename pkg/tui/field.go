package tui

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

func NewField(options *FieldOptions) *Field {
  // TODO: we can maybe add some padding by directly accessing the model and tampering wiht GetBounds
  field := &Field{ }
  field.SetOrientation(views.Horizontal)

  o := &InputOptions{ InitialValue: options.InitialValue, Type: options.InputType }
  i := NewInput(o)
  i.SetContent(options.InitialValue)

  l := views.NewSimpleStyledText()
  l.SetStyle(tcell.StyleDefault.Attributes(tcell.AttrBold))
  l.SetText(options.Label + ": ")

  field.AddWidget(l, 0)
  field.AddWidget(i, 1)

  field.input = i
  field.label = l

  return field
}
