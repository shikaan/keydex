package form

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
  label string
  initialValue string
  mask string
  labelWidth int16
  fieldWidth int16
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
  field := &Field{ }
  field.SetOrientation(0)

  i := NewInput()
  i.SetContent(options.initialValue)

  l := views.NewSimpleStyledText()
  l.SetStyle(tcell.StyleDefault.Attributes(tcell.AttrBold))
  l.SetText(options.label + ": ")

  field.AddWidget(l, 0)
  field.AddWidget(i, 1)

  field.input = i
  field.label = l

  return field
}
