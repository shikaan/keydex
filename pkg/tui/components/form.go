package components

import (
	"github.com/gdamore/tcell/v2"
	"github.com/gdamore/tcell/v2/views"
)

type Form struct {
  views.BoxLayout
}

func (f *Form) HandleEvent(ev tcell.Event) bool {
	switch ev := ev.(type) {
	case *tcell.EventKey:
		if ev.Key() == tcell.KeyUp {
      f.MoveFocus(-1)
      return true
    }
    if ev.Key() == tcell.KeyDown {
      f.MoveFocus(1)
      return true
    }
  }
	
  return f.BoxLayout.HandleEvent(ev)
}

// Moves focus by `offset` fields
func (f *Form) MoveFocus (offset int) {
  fs := f.Focusables()
  count := len(fs)
  current := -1

  for i, f := range fs {
    if f.HasFocus() {
      current = i
      f.SetFocus(false)
      break
    }
  }

  notFound := current == -1
  if notFound {
    return 
  }

  newIndex := (count + current + offset) % count

  fs[newIndex].SetFocus(true)
}

// Returns the subset of Widgets that can have focus
func (f *Form) Focusables() []Focusable {
  ws := f.Widgets()
  result := []Focusable{}

  for _, w := range ws {
    switch w := w.(type) {
    case Focusable:
      result = append(result, w)
    }
  }
  return result
}

func NewForm() *Form {
  f := &Form{}
  f.SetOrientation(views.Vertical)

  return f
}
