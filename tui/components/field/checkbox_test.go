package field

import (
	"testing"

	"github.com/gdamore/tcell/v2"
	"github.com/shikaan/keydex/tui/components/line"
)

func Test_checkboxModel_GetCell(t *testing.T) {
	checked := [3]rune{'[', 'X', ']'}
	tests := []struct {
		name          string
		fields        checkboxModel
		x             int
		y             int
		wantRune      rune
		wantRuneWidth int
	}{
		{"opening bracket", checkboxModel{cells: checked}, 0, 0, '[', 1},
		{"checked mark", checkboxModel{cells: checked}, 1, 0, 'X', 1},
		{"unchecked mark", checkboxModel{cells: [3]rune{'[', ' ', ']'}}, 1, 0, ' ', 1},
		{"closing bracket", checkboxModel{cells: checked}, 2, 0, ']', 1},
		{"EMPTY_CELL when x out of bounds", checkboxModel{cells: checked}, 3, 0, line.EMPTY_CELL, 1},
		{"EMPTY_CELL when y out of bounds", checkboxModel{cells: checked}, 1, 1, line.EMPTY_CELL, 1},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			m := &checkboxModel{
				content:         tt.fields.content,
				hasFocus:        tt.fields.hasFocus,
				disabled:        tt.fields.disabled,
				cells:           tt.fields.cells,
				keyPressHandler: tt.fields.keyPressHandler,
				changeHandler:   tt.fields.changeHandler,
				focusHandler:    tt.fields.focusHandler,
			}
			gotRune, _, _, gotRuneWidth := m.GetCell(tt.x, tt.y)
			if gotRune != tt.wantRune {
				t.Errorf("checkboxModel.GetCell() got rune = %v, want %v", gotRune, tt.wantRune)
			}
			if gotRuneWidth != tt.wantRuneWidth {
				t.Errorf("checkboxModel.GetCell() got rune width = %v, want %v", gotRuneWidth, tt.wantRuneWidth)
			}
		})
	}
}

func Test_checkboxModel_GetCursor(t *testing.T) {
	tests := []struct {
		name        string
		hasFocus    bool
		wantX       int
		wantY       int
		wantEnabled bool
		wantVisible bool
	}{
		{"cursor visible with focus", true, 1, 0, true, true},
		{"cursor hidden without focus", false, 1, 0, true, false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			m := &checkboxModel{hasFocus: tt.hasFocus}
			gotX, gotY, gotEnabled, gotVisible := m.GetCursor()
			if gotX != tt.wantX {
				t.Errorf("checkboxModel.GetCursor() x = %v, want %v", gotX, tt.wantX)
			}
			if gotY != tt.wantY {
				t.Errorf("checkboxModel.GetCursor() y = %v, want %v", gotY, tt.wantY)
			}
			if gotEnabled != tt.wantEnabled {
				t.Errorf("checkboxModel.GetCursor() enabled = %v, want %v", gotEnabled, tt.wantEnabled)
			}
			if gotVisible != tt.wantVisible {
				t.Errorf("checkboxModel.GetCursor() visible = %v, want %v", gotVisible, tt.wantVisible)
			}
		})
	}
}

func TestCheckbox_SetContent(t *testing.T) {
	tests := []struct {
		name      string
		value     bool
		wantCells [3]rune
	}{
		{"checked", true, [3]rune{'[', 'X', ']'}},
		{"unchecked", false, [3]rune{'[', ' ', ']'}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := &Checkbox{}
			c.SetContent(tt.value)
			if c.model.content != tt.value {
				t.Errorf("Checkbox.SetContent() got content = %v, want %v", c.model.content, tt.value)
			}
			if c.model.cells != tt.wantCells {
				t.Errorf("Checkbox.SetContent() got cells = %v, want %v", c.model.cells, tt.wantCells)
			}
			if c.GetContent() != tt.value {
				t.Errorf("Checkbox.GetContent() = %v, want %v", c.GetContent(), tt.value)
			}
		})
	}
}

func TestCheckbox_HandleEvent(t *testing.T) {
	tests := []struct {
		name        string
		initial     bool
		hasFocus    bool
		disabled    bool
		event       tcell.Event
		wantHandled bool
		wantContent bool
	}{
		{"no focus", false, false, false, tcell.NewEventKey(tcell.KeyEnter, 0, 0), false, false},
		{"enter checks", false, true, false, tcell.NewEventKey(tcell.KeyEnter, 0, 0), true, true},
		{"enter unchecks", true, true, false, tcell.NewEventKey(tcell.KeyEnter, 0, 0), true, false},
		{"other key is ignored", false, true, false, tcell.NewEventKey(tcell.KeyRune, 'x', 0), false, false},
		{"non-key event is ignored", false, true, false, tcell.NewEventResize(10, 10), false, false},
		{"disabled (enter)", false, true, true, tcell.NewEventKey(tcell.KeyEnter, 0, 0), false, false},
		{"disabled and checked (enter)", true, true, true, tcell.NewEventKey(tcell.KeyEnter, 0, 0), false, true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := NewCheckbox(&CheckboxOptions{InitialValue: tt.initial, Disabled: tt.disabled})
			c.SetFocus(tt.hasFocus)
			handled := c.HandleEvent(tt.event)
			if handled != tt.wantHandled {
				t.Errorf("Checkbox.HandleEvent() handled = %v, want %v", handled, tt.wantHandled)
			}
			if c.GetContent() != tt.wantContent {
				t.Errorf("Checkbox.HandleEvent() content = %v, want %v", c.GetContent(), tt.wantContent)
			}
		})
	}
}

func TestCheckbox_HandleEvent_handlers(t *testing.T) {
	enter := tcell.NewEventKey(tcell.KeyEnter, 0, 0)

	t.Run("key press handler handles the event (skips toggle)", func(t *testing.T) {
		c := NewCheckbox(&CheckboxOptions{})
		c.SetFocus(true)
		c.OnKeyPress(func(ev *tcell.EventKey) bool { return true })

		handled := c.HandleEvent(enter)

		if !handled {
			t.Errorf("Checkbox.HandleEvent() expected event to be handled")
		}
		if c.GetContent() {
			t.Errorf("Checkbox.HandleEvent() expected content not to change")
		}
	})

	t.Run("key press handler does not handle the event (toggles)", func(t *testing.T) {
		c := NewCheckbox(&CheckboxOptions{})
		c.SetFocus(true)
		c.OnKeyPress(func(ev *tcell.EventKey) bool { return false })

		handled := c.HandleEvent(enter)

		if !handled {
			t.Errorf("Checkbox.HandleEvent() expected event to be handled")
		}
		if !c.GetContent() {
			t.Errorf("Checkbox.HandleEvent() expected content to change")
		}
	})

	t.Run("uses the change handler (handled)", func(t *testing.T) {
		triggered := false
		c := NewCheckbox(&CheckboxOptions{})
		c.SetFocus(true)
		c.OnChange(func(ev tcell.Event) bool {
			triggered = true
			return true
		})

		handled := c.HandleEvent(enter)

		if !triggered {
			t.Errorf("Checkbox.HandleEvent() expected change handler to be triggered")
		}
		if !handled {
			t.Errorf("Checkbox.HandleEvent() expected change to be handled")
		}
	})

	t.Run("uses the change handler (unhandled)", func(t *testing.T) {
		triggered := false
		c := NewCheckbox(&CheckboxOptions{})
		c.SetFocus(true)
		c.OnChange(func(ev tcell.Event) bool {
			triggered = true
			return false
		})

		handled := c.HandleEvent(enter)

		if !triggered {
			t.Errorf("Checkbox.HandleEvent() expected change handler to be triggered")
		}
		if handled {
			t.Errorf("Checkbox.HandleEvent() expected change not to be handled")
		}
		if !c.GetContent() {
			t.Errorf("Checkbox.HandleEvent() expected content to change regardless of the handler result")
		}
	})

	t.Run("does not use the change handler when disabled", func(t *testing.T) {
		triggered := false
		c := NewCheckbox(&CheckboxOptions{Disabled: true})
		c.SetFocus(true)
		c.OnChange(func(ev tcell.Event) bool {
			triggered = true
			return true
		})

		c.HandleEvent(enter)

		if triggered {
			t.Errorf("Checkbox.HandleEvent() expected change handler not to be triggered")
		}
	})
}
