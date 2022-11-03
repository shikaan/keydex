package form

import (
	"github.com/gdamore/tcell/v2"
	"github.com/gdamore/tcell/v2/views"
)

type Field struct {
  input Input
  label views.SimpleStyledText

  views.BoxLayout
}

func NewField(label, initialValue string) *Field {
  field := &Field{ }
  field.SetOrientation(0)

  i := NewInput()
  i.SetContent(initialValue)

  l := views.NewSimpleStyledText()
  l.SetStyle(tcell.StyleDefault.Attributes(tcell.AttrBold))
  l.SetText(label + ": ")

  field.AddWidget(l, 0)
  field.AddWidget(i, 1)

  field.input = *i
  field.label = *l

  return field
}
