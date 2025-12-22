package autocomplete

import (
	"reflect"
	"testing"

	"github.com/gdamore/tcell/v2"
	"github.com/shikaan/keydex/tui/components/line"
)

type HandlerSpy struct {
	handler func() bool
	calls   int
}

func NewHandlerSpy(value bool) *HandlerSpy {
	spy := &HandlerSpy{}
	spy.handler = func() bool {
		spy.calls++
		return value
	}
	return spy
}

func TestOption_SetFocus(t *testing.T) {
	tests := []struct {
		name         string
		text         string
		on           bool
		want         string
		focusHandler *HandlerSpy
		wantCalls    int
	}{
		{"prefixes with a >", "test", true, "> test", nil, 0},
		{"removes >", "test", false, "test", nil, 0},
		{"triggers handler", "test", true, "> test", NewHandlerSpy(true), 1},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			i := newOption()
			i.SetContent(tt.text)
			if tt.focusHandler != nil {
				i.OnFocus(tt.focusHandler.handler)
			}
			i.SetFocus(tt.on)
			if i.model.content != tt.want {
				t.Errorf("Option.SetFocus() = %v, want %v", i.model.content, tt.want)
			}
			if i.model.hasFocus != tt.on {
				t.Errorf("Option.SetFocus() got hasFocus = %v, want hasFocus %v", i.model.hasFocus, tt.on)
			}
			if i.model.focusHandler != nil && tt.focusHandler.calls != tt.wantCalls {
				t.Errorf("Option.SetFocus() got calls %v to focusHandler, wanted %v", tt.focusHandler.calls, tt.wantCalls)
			}
		})
	}
}

func TestOption_GetContent(t *testing.T) {
	t.Run("retrieve content without caret", func(t *testing.T) {
		o := newOption()
		o.SetContent("hello world")
		o.SetFocus(true)

		got := o.GetContent()
		want := "hello world"
		if got != want {
			t.Errorf("Option.GetContent() = %q, want %q", got, want)
		}
	})

	t.Run("retrieve content", func(t *testing.T) {
		o := newOption()
		o.SetContent("hello world")
		o.SetFocus(false)

		got := o.GetContent()
		want := "hello world"
		if got != want {
			t.Errorf("Option.GetContent() = %q, want %q", got, want)
		}
	})
}

func TestOption_HandleEvent(t *testing.T) {
	tests := []struct {
		name          string
		event         tcell.Event
		focus         bool
		selectHandler *HandlerSpy
		wantCalls     int
		wantHandled   bool
	}{
		{"skips events without focus", tcell.NewEventKey(tcell.KeyEnter, 0, 0), false, NewHandlerSpy(true), 0, false},
		{"handles select event (with handler)", tcell.NewEventKey(tcell.KeyEnter, 0, 0), true, NewHandlerSpy(true), 1, true},
		{"doesn't handle select event (with handler)", tcell.NewEventKey(tcell.KeyEnter, 0, 0), true, NewHandlerSpy(false), 1, false},
		{"doesn't handle select event (without handler)", tcell.NewEventKey(tcell.KeyEnter, 0, 0), true, nil, 0, false},
		{"trigger select handler with other events", tcell.NewEventKey(tcell.KeyUp, 0, 0), true, NewHandlerSpy(true), 0, false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			o := newOption()
			if tt.selectHandler != nil {
				o.OnSelect(tt.selectHandler.handler)
			}
			o.SetFocus(tt.focus)
			got := o.HandleEvent(tt.event)
			if got != tt.wantHandled {
				t.Errorf("Option.HandleEvent() = %v, want %v", got, tt.wantHandled)
			}
			if o.model.selectHandler != nil && tt.selectHandler.calls != tt.wantCalls {
				t.Errorf("Option.HandleEvent() got calls %v to selectHandler, wanted %v", tt.selectHandler.calls, tt.wantCalls)
			}
		})
	}
}
func TestOptionModel_GetCell(t *testing.T) {
	tests := []struct {
		name      string
		content   string
		x, y      int
		wantRune  rune
		wantStyle tcell.Style
		wantComb  []rune
		wantWidth int
	}{
		{"out of bounds y", "test", 0, 1, line.EMPTY_CELL, tcell.StyleDefault, nil, 1},
		{"negative x", "test", -1, 0, line.EMPTY_CELL, tcell.StyleDefault, nil, 1},
		{"x out of bounds", "test", 5, 0, line.EMPTY_CELL, tcell.StyleDefault, nil, 1},
		{"valid cell", "test", 1, 0, 'e', tcell.StyleDefault, nil, 1},
		{"empty cell", "\x00", 0, 0, 0, tcell.StyleDefault, nil, 1},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			m := &optionModel{
				content: tt.content,
				runes:   line.NewPaddedLine(tt.content),
				style:   tcell.StyleDefault,
			}
			gotRune, gotStyle, gotComb, gotWidth := m.GetCell(tt.x, tt.y)
			if gotRune != tt.wantRune {
				t.Errorf("optionModel.GetCell() gotRune = %v, want %v", gotRune, tt.wantRune)
			}
			if gotStyle != tt.wantStyle {
				t.Errorf("optionModel.GetCell() gotStyle = %v, want %v", gotStyle, tt.wantStyle)
			}
			if !reflect.DeepEqual(gotComb, tt.wantComb) {
				t.Errorf("optionModel.GetCell() gotComb = %v, want %v", gotComb, tt.wantComb)
			}
			if gotWidth != tt.wantWidth {
				t.Errorf("optionModel.GetCell() gotWidth = %v, want %v", gotWidth, tt.wantWidth)
			}
		})
	}
}
