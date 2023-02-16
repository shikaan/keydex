package components

import (
	"github.com/gdamore/tcell/v2"
	"github.com/gdamore/tcell/v2/views"
)

type Focusable interface {
	SetFocus(on bool)
	HasFocus() bool
}

type WithFocusables struct {
	views.BoxLayout
}

func (wf *WithFocusables) HandleEvent(ev tcell.Event) bool {
	switch ev := ev.(type) {
	case *tcell.EventKey:
		if ev.Key() == tcell.KeyUp {
			wf.MoveFocus(-1)
			return true
		}
		if ev.Key() == tcell.KeyDown {
			wf.MoveFocus(1)
			return true
		}
	}

	return wf.BoxLayout.HandleEvent(ev)
}

// Moves focus by `offset` focusables
func (wf *WithFocusables) MoveFocus(offset int) {
	fs := wf.Focusables()
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
func (wf *WithFocusables) Focusables() []Focusable {
	ws := wf.Widgets()
	result := []Focusable{}

	for _, w := range ws {
		switch w := w.(type) {
		case Focusable:
			result = append(result, w)
		}
	}
	return result
}
