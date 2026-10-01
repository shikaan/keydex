package field

import (
	"reflect"
	"testing"

	"github.com/gdamore/tcell/v2"
	"github.com/shikaan/keydex/tui/components/line"
)

func Test_inputModel_GetCell(t *testing.T) {
	tests := []struct {
		name          string
		fields        inputModel
		x             int
		y             int
		wantRune      rune
		wantRuneWidth int
	}{
		{"EMPTY_CELL when out of bounds", inputModel{cells: [][]rune{{'L', 'O', 'L'}}}, 2, 2, line.EMPTY_CELL, 1},
		{"EMPTY_CELL when cursor out of password bounds (password)",
			inputModel{cells: [][]rune{{'a', 'a', 'a', 'a', 'a', 'a', 'a', 'a', 'a', 'a'}}, hidden: true},
			PASSWORD_FIELD_LENGTH + 1, 0, // This is longer than a password, but shorter than field itself
			line.EMPTY_CELL, 1},
		{"EMPTY_CELL when cursor out of field bounds (password)",
			inputModel{cells: [][]rune{{'a', 'a', 'a', 'a', 'a', 'a', 'a', 'a', 'a', 'a'}}, hidden: true},
			14, 0, // This is longer than the field
			line.EMPTY_CELL, 1},
		{"* when password", inputModel{cells: [][]rune{{'L', 'O', 'L'}}, hidden: true}, 1, 0, '*', 1},
		{"* when non-byte password", inputModel{cells: [][]rune{{'🤖', '✅', '😂'}}, hidden: true}, 1, 0, '*', 1},
		{"returns byte with size 1 (only byte)", inputModel{cells: [][]rune{{'L', 'O', 'L'}}}, 1, 0, 'O', 1},
		{"returns emoji with size 2 (only emoji)", inputModel{cells: [][]rune{{'🤖', '✅', '😂'}}}, 1, 0, '✅', 2},
		{"returns byte with size 1 (mixed)", inputModel{cells: [][]rune{{'I', '🤖'}}}, 0, 0, 'I', 1},
		{"returns emoji with size 2 (mixed)", inputModel{cells: [][]rune{{'I', '🤖'}}}, 1, 0, '🤖', 2},
		{"returns emoji with size 2 (mixed)", inputModel{cells: [][]rune{{'I', '🤖'}}}, 1, 0, '🤖', 2},
		{"EMPTY_CELL with non-print characters", inputModel{cells: [][]rune{{'I', '\a'}}}, 1, 0, line.EMPTY_CELL, 1},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			m := &inputModel{
				content:         tt.fields.content,
				cells:           tt.fields.cells,
				width:           tt.fields.width,
				height:          tt.fields.height,
				x:               tt.fields.x,
				y:               tt.fields.y,
				style:           tt.fields.style,
				hasFocus:        tt.fields.hasFocus,
				hidden:          tt.fields.hidden,
				keyPressHandler: tt.fields.keyPressHandler,
				changeHandler:   tt.fields.changeHandler,
				focusHandler:    tt.fields.focusHandler,
			}
			gotRune, _, _, gotRuneWidth := m.GetCell(tt.x, tt.y)
			if gotRune != tt.wantRune {
				t.Errorf("inputModel.GetCell() got rune = %v, want %v", gotRune, tt.wantRune)
			}
			if gotRuneWidth != tt.wantRuneWidth {
				t.Errorf("inputModel.GetCell() got rune width = %v, want %v", gotRuneWidth, tt.wantRuneWidth)
			}
		})
	}
}

func Test_inputModel_SetCursor(t *testing.T) {
	tests := []struct {
		name   string
		fields inputModel
		x      int
		y      int
		wantX  int
		wantY  int
	}{
		{"x before beginning", inputModel{cells: [][]rune{{'T', 'e', 's'}, {'T', 'e', 's'}}}, -1, 0, 0, 0},
		{"x after ending", inputModel{cells: [][]rune{{'T', 'e', 's'}, {'T', 'e', 's'}}}, 10, 0, 3, 0},
		{"accommodates for extra x char", inputModel{cells: [][]rune{{'T', 'e', 's'}}, x: 2, y: 0}, 3, 0, 3, 0},
		{"y before beginning", inputModel{cells: [][]rune{{'T', 'e', 's'}, {'T', 'e', 's'}}}, 0, -1, 0, 0},
		{"y after ending", inputModel{cells: [][]rune{{'T', 'e', 's'}, {'T', 'e', 's'}}}, 0, 2, 0, 1},
		{"legit horizontal", inputModel{cells: [][]rune{{'T', 'e', 's'}, {'T', 'e', 's'}}}, 1, 0, 1, 0},
		{"legit vertical", inputModel{cells: [][]rune{{'T', 'e', 's'}, {'T', 'e', 's'}}}, 0, 1, 0, 1},
		{"legit both", inputModel{cells: [][]rune{{'T', 'e', 's'}, {'T', 'e', 's'}}}, 1, 1, 1, 1},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			m := &inputModel{
				content:         tt.fields.content,
				cells:           tt.fields.cells,
				width:           tt.fields.width,
				height:          tt.fields.height,
				x:               tt.fields.x,
				y:               tt.fields.y,
				style:           tt.fields.style,
				hasFocus:        tt.fields.hasFocus,
				hidden:          tt.fields.hidden,
				keyPressHandler: tt.fields.keyPressHandler,
				changeHandler:   tt.fields.changeHandler,
				focusHandler:    tt.fields.focusHandler,
			}
			m.SetCursor(tt.x, tt.y)
			if m.x != tt.wantX {
				t.Errorf("inputModel.SetCursor() got = %v, want %v", m.x, tt.wantX)
			}
			if m.y != tt.wantY {
				t.Errorf("inputModel.SetCursor() got1 = %v, want %v", m.y, tt.wantY)
			}
		})
	}
}

func Test_inputModel_MoveCursor(t *testing.T) {
	tests := []struct {
		name   string
		fields inputModel
		x      int
		y      int
		wantX  int
		wantY  int
	}{
		{"x before beginning", inputModel{cells: [][]rune{{'T', 'e', 's'}, {'T', 'e', 's'}}, x: 1, y: 1}, -2, 0, 0, 1},
		{"x after ending", inputModel{cells: [][]rune{{'T', 'e', 's'}, {'T', 'e', 's'}}, x: 1, y: 1}, 10, 0, 3, 1},
		{"accommodates for extra x char", inputModel{cells: [][]rune{{'T', 'e', 's'}}, x: 2, y: 0}, 1, 0, 3, 0},
		{"y before beginning", inputModel{cells: [][]rune{{'T', 'e', 's'}, {'T', 'e', 's'}}, x: 1, y: 1}, 0, -2, 1, 0},
		{"y after ending", inputModel{cells: [][]rune{{'T', 'e', 's'}, {'T', 'e', 's'}}, x: 1, y: 1}, 0, 1, 1, 1},
		{"legit horizontal", inputModel{cells: [][]rune{{'T', 'e', 's'}, {'T', 'e', 's'}}, x: 1, y: 1}, 1, 0, 2, 1},
		{"legit vertical", inputModel{cells: [][]rune{{'T', 'e', 's'}, {'T', 'e', 's'}}, x: 1, y: 0}, 0, 1, 1, 1},
		{"legit both", inputModel{cells: [][]rune{{'T', 'e', 's'}, {'T', 'e', 's'}}, x: 1, y: 0}, 1, 1, 2, 1},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			m := &inputModel{
				content:         tt.fields.content,
				cells:           tt.fields.cells,
				width:           tt.fields.width,
				height:          tt.fields.height,
				x:               tt.fields.x,
				y:               tt.fields.y,
				style:           tt.fields.style,
				hasFocus:        tt.fields.hasFocus,
				hidden:          tt.fields.hidden,
				keyPressHandler: tt.fields.keyPressHandler,
				changeHandler:   tt.fields.changeHandler,
				focusHandler:    tt.fields.focusHandler,
			}
			m.MoveCursor(tt.x, tt.y)
			if m.x != tt.wantX {
				t.Errorf("inputModel.MoveCursor() x = %v, want %v", m.x, tt.wantX)
			}
			if m.y != tt.wantY {
				t.Errorf("inputModel.MoveCursor() y = %v, want %v", m.y, tt.wantY)
			}
		})
	}
}

func Test_inputModel_GetRuneAtPosition(t *testing.T) {
	tests := []struct {
		name        string
		fields      inputModel
		x           int
		y           int
		wantRune    rune
		wantHOffset int
	}{
		{"get * with password", inputModel{cells: [][]rune{{'T', 'e'}}, hidden: true}, 1, 0, '*', 1},
		{"get byte char", inputModel{cells: [][]rune{{'T', 'e', 's'}}}, 1, 0, 'e', 1},
		{"get unicode char", inputModel{cells: [][]rune{{'T', '✅', line.PAD_BYTE, 's'}}}, 1, 0, '✅', 1},
		{"get unicode char on PAD_BYTE", inputModel{cells: [][]rune{{'T', '✅', line.PAD_BYTE, 's'}}}, 2, 0, '✅', 1},
		{"EMPTY_CELL when out of bounds", inputModel{cells: [][]rune{{'T', 's'}}}, 2, 0, line.EMPTY_CELL, -1},
		{"EMPTY_CELL when out of bounds (password)", inputModel{hidden: true, cells: [][]rune{{'T', 's'}}}, PASSWORD_FIELD_LENGTH + 1, 0, line.EMPTY_CELL, -1},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			m := &inputModel{
				content:         tt.fields.content,
				cells:           tt.fields.cells,
				width:           tt.fields.width,
				height:          tt.fields.height,
				x:               tt.fields.x,
				y:               tt.fields.y,
				style:           tt.fields.style,
				hasFocus:        tt.fields.hasFocus,
				hidden:          tt.fields.hidden,
				keyPressHandler: tt.fields.keyPressHandler,
				changeHandler:   tt.fields.changeHandler,
				focusHandler:    tt.fields.focusHandler,
			}
			gotRune, gotHOffset := m.getRuneAtPosition(tt.x, tt.y)
			if gotRune != tt.wantRune {
				t.Errorf("inputModel.GetRuneAtPosition() got rune = %v, want %v", gotRune, tt.wantRune)
			}
			if gotHOffset != tt.wantHOffset {
				t.Errorf("inputModel.GetRuneAtPosition() got horizontal offset = %v, want %v", gotHOffset, tt.wantHOffset)
			}
		})
	}
}

func TestInput_SetContent(t *testing.T) {
	tests := []struct {
		name        string
		text        string
		wantCells   [][]rune
		wantContent string
	}{
		{"empty string", "", [][]rune{{}}, ""},
		{"bytes only", "test", [][]rune{{'t', 'e', 's', 't'}}, "test"},
		{"non-ASCII only", "🤖ß", [][]rune{{'🤖', line.PAD_BYTE, 'ß'}}, "🤖ß"},
		{"mixed input", "I💙ßs", [][]rune{{'I', '💙', line.PAD_BYTE, 'ß', 's'}}, "I💙ßs"},
		{"line breaks UNIX", "I💙ßs\n&🥔", [][]rune{{'I', '💙', line.PAD_BYTE, 'ß', 's'}, {'&', '🥔', line.PAD_BYTE}}, "I💙ßs\n&🥔"},
		{"line breaks Windows", "I💙ßs\r\n&🥔", [][]rune{{'I', '💙', line.PAD_BYTE, 'ß', 's'}, {'&', '🥔', line.PAD_BYTE}}, "I💙ßs\r\n&🥔"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			i := &Input{}
			i.SetContent(tt.text)
			if !reflect.DeepEqual(tt.wantCells, i.model.cells) {
				t.Errorf("Input.SetContent() got cells = %v, want %v", i.model.cells, tt.wantCells)
			}
			if tt.wantContent != i.model.content {
				t.Errorf("Input.SetContent() got content = %v, want %v", i.model.content, tt.wantContent)
			}
		})
	}
}

func Test_toString(t *testing.T) {
	tests := []struct {
		name  string
		cells [][]rune
		want  string
	}{
		{"bytes", [][]rune{{'l', 'o', 'l'}}, "lol"},
		{"non-ascii", [][]rune{{'✅', 'æ'}}, "✅æ"},
		{"multi-line", [][]rune{{'✅', 'æ'}, {'l'}}, "✅æ\nl"},
		{"non-print chars", [][]rune{{'✅', '\a', '\f'}, {'l'}}, "✅\nl"},
		{"empty lines", [][]rune{{'t', 'e'}, {}}, "te\n"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := toString(tt.cells); got != tt.want {
				t.Errorf("toString() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestInput_handleCellsUpdate(t *testing.T) {
	t.Run("updates content and moves cursor", func(t *testing.T) {
		i := &Input{model: &inputModel{}}
		handled := i.handleCellsUpdate(&tcell.EventKey{}, func() (int, int) {
			i.model.cells = [][]rune{{'l'}}
			return 1, 0
		})

		wantCells := [][]rune{{'l'}}
		wantX := 1
		wantY := 0
		if !reflect.DeepEqual(i.model.cells, wantCells) {
			t.Errorf("Input.handleCellsUpdate() got cells = %v, want %v", i.model.cells, wantCells)
		}
		if i.model.x != wantX {
			t.Errorf("Input.handleCellsUpdate() got x = %v, want %v", i.model.x, wantX)
		}
		if i.model.y != wantY {
			t.Errorf("Input.handleCellsUpdate() got y = %v, want %v", i.model.y, wantY)
		}
		if !handled {
			t.Errorf("Input.handleCellsUpdate() expected event to be handled")
		}
	})

	t.Run("uses the change handler (handled)", func(t *testing.T) {
		triggered := false
		i := NewInput(&InputOptions{})
		i.model.changeHandler = func(ev tcell.Event) bool {
			triggered = true
			return true
		}
		handled := i.handleCellsUpdate(&tcell.EventKey{}, func() (int, int) { return 0, 0 })

		if !triggered {
			t.Errorf("Input.handleCellsUpdate() expected change handler to be triggered")
		}

		if !handled {
			t.Errorf("Input.handleCellsUpdate() change to be handled")
		}
	})

	t.Run("uses the change handler (unhandled)", func(t *testing.T) {
		triggered := false
		i := NewInput(&InputOptions{})
		i.model.changeHandler = func(ev tcell.Event) bool {
			triggered = true
			return false
		}
		handled := i.handleCellsUpdate(&tcell.EventKey{}, func() (int, int) { return 0, 0 })

		if !triggered {
			t.Errorf("Input.handleCellsUpdate() expected change handler to be triggered")
		}

		if handled {
			t.Errorf("Input.handleCellsUpdate() change not to be handled")
		}
	})
}
func TestInput_HandleEvent(t *testing.T) {
	tests := []struct {
		name        string
		content     string
		hasFocus    bool
		disabled    bool
		x, y        int
		event       tcell.Event
		wantHandled bool
		wantX       int
		wantY       int
		wantCells   [][]rune
	}{
		{
			name:        "no focus",
			content:     "hello",
			hasFocus:    false,
			event:       tcell.NewEventKey(tcell.KeyRune, 'a', 0),
			wantHandled: false,
		},
		{
			name:        "key left",
			hasFocus:    true,
			content:     "abc",
			x:           1,
			y:           0,
			event:       tcell.NewEventKey(tcell.KeyLeft, 0, 0),
			wantHandled: true,
			wantX:       0,
			wantY:       0,
		},
		{
			name:        "key left (non-ascii)",
			content:     "a✅c",
			hasFocus:    true,
			x:           2,
			y:           0,
			event:       tcell.NewEventKey(tcell.KeyLeft, 0, 0),
			wantHandled: true,
			wantX:       1,
			wantY:       0,
		},
		{
			name:        "key right",
			content:     "abc",
			x:           1,
			y:           0,
			hasFocus:    true,
			event:       tcell.NewEventKey(tcell.KeyRight, 0, 0),
			wantHandled: true,
			wantX:       2,
			wantY:       0,
		},
		{
			name:        "key right (non-ascii)",
			content:     "a✅c",
			hasFocus:    true,
			x:           1,
			y:           0,
			event:       tcell.NewEventKey(tcell.KeyRight, 0, 0),
			wantHandled: true,
			wantX:       3,
			wantY:       0,
		},
		{
			name:        "key down (within bounds)",
			content:     "abc\ndef",
			hasFocus:    true,
			x:           1,
			y:           0,
			event:       tcell.NewEventKey(tcell.KeyDown, 0, 0),
			wantHandled: true,
			wantX:       1,
			wantY:       1,
		},
		{
			name:        "key down (out of bounds)",
			content:     "abc\ndef",
			hasFocus:    true,
			x:           1,
			y:           1,
			event:       tcell.NewEventKey(tcell.KeyDown, 0, 0),
			wantHandled: false,
			wantX:       1,
			wantY:       1,
		},
		{
			name:        "key up",
			content:     "abc\ndef",
			hasFocus:    true,
			x:           1,
			y:           1,
			event:       tcell.NewEventKey(tcell.KeyUp, 0, 0),
			wantHandled: true,
			wantX:       1,
			wantY:       0,
		},
		{
			name:        "key up (out of bounds)",
			content:     "abc\ndef",
			hasFocus:    true,
			x:           1,
			y:           0,
			event:       tcell.NewEventKey(tcell.KeyUp, 0, 0),
			wantHandled: false,
			wantX:       1,
			wantY:       0,
		},
		{
			name:        "key enter (end of line)",
			content:     "abc",
			hasFocus:    true,
			x:           3,
			y:           0,
			event:       tcell.NewEventKey(tcell.KeyEnter, 0, 0),
			wantHandled: true,
			wantX:       0,
			wantY:       1,
			wantCells:   [][]rune{{'a', 'b', 'c'}, {}},
		},
		{
			name:        "key enter (mid-line)",
			content:     "abc",
			hasFocus:    true,
			x:           1,
			y:           0,
			event:       tcell.NewEventKey(tcell.KeyEnter, 0, 0),
			wantHandled: true,
			wantX:       0,
			wantY:       1,
			wantCells:   [][]rune{{'a'}, {'b', 'c'}},
		},
		{
			name:        "key rune (byte)",
			content:     "abc",
			hasFocus:    true,
			x:           1,
			y:           0,
			event:       tcell.NewEventKey(tcell.KeyRune, 'd', 0),
			wantHandled: true,
			wantX:       2,
			wantY:       0,
			wantCells:   [][]rune{{'a', 'd', 'b', 'c'}},
		},
		{
			name:        "key rune (non-ascii)",
			content:     "abc",
			hasFocus:    true,
			x:           1,
			y:           0,
			event:       tcell.NewEventKey(tcell.KeyRune, '✅', 0),
			wantHandled: true,
			wantX:       3,
			wantY:       0,
			wantCells:   [][]rune{{'a', '✅', 0, 'b', 'c'}},
		},
		{
			name:        "key backspace",
			content:     "abc",
			hasFocus:    true,
			x:           2,
			y:           0,
			event:       tcell.NewEventKey(tcell.KeyBackspace, 0, 0),
			wantHandled: true,
			wantX:       1,
			wantY:       0,
			wantCells:   [][]rune{{'a', 'c'}},
		},
		{
			name:        "key backspace (start of line)",
			content:     "abc\nd",
			hasFocus:    true,
			x:           0,
			y:           1,
			event:       tcell.NewEventKey(tcell.KeyBackspace, 0, 0),
			wantHandled: true,
			wantX:       3,
			wantY:       0,
			wantCells:   [][]rune{{'a', 'b', 'c', 'd'}},
		},
		{
			name:        "key backspace (start of block)",
			content:     "abc\nd",
			hasFocus:    true,
			x:           0,
			y:           0,
			event:       tcell.NewEventKey(tcell.KeyBackspace, 0, 0),
			wantHandled: true,
			wantX:       0,
			wantY:       0,
			wantCells:   [][]rune{{'a', 'b', 'c'}, {'d'}},
		},
		{
			name:        "key delete",
			content:     "abc",
			hasFocus:    true,
			x:           1,
			y:           0,
			event:       tcell.NewEventKey(tcell.KeyDelete, 0, 0),
			wantHandled: true,
			wantX:       1,
			wantY:       0,
			wantCells:   [][]rune{{'a', 'c'}},
		},
		{
			name:        "key delete (end of line)",
			content:     "abc\nd",
			hasFocus:    true,
			x:           3,
			y:           0,
			event:       tcell.NewEventKey(tcell.KeyDelete, 0, 0),
			wantHandled: true,
			wantX:       3,
			wantY:       0,
			wantCells:   [][]rune{{'a', 'b', 'c', 'd'}},
		},
		{
			name:        "key delete (end of block)",
			content:     "abc\nd",
			hasFocus:    true,
			x:           1,
			y:           1,
			event:       tcell.NewEventKey(tcell.KeyDelete, 0, 0),
			wantHandled: true,
			wantX:       1,
			wantY:       1,
			wantCells:   [][]rune{{'a', 'b', 'c'}, {'d'}},
		},
		{
			name:        "disabled input (rune)",
			content:     "abc",
			hasFocus:    true,
			disabled:    true,
			x:           1,
			y:           0,
			event:       tcell.NewEventKey(tcell.KeyRune, 'd', 0),
			wantHandled: false,
			wantX:       1,
			wantY:       0,
			wantCells:   [][]rune{{'a', 'b', 'c'}},
		},
		{
			name:        "disabled input (delete)",
			content:     "abc",
			hasFocus:    true,
			disabled:    true,
			x:           1,
			y:           0,
			event:       tcell.NewEventKey(tcell.KeyDelete, 0, 0),
			wantHandled: false,
			wantX:       1,
			wantY:       0,
			wantCells:   [][]rune{{'a', 'b', 'c'}},
		},
		{
			name:        "disabled input (enter)",
			content:     "abc",
			hasFocus:    true,
			disabled:    true,
			x:           3,
			y:           0,
			event:       tcell.NewEventKey(tcell.KeyEnter, 0, 0),
			wantHandled: false,
			wantX:       3,
			wantY:       0,
			wantCells:   [][]rune{{'a', 'b', 'c'}},
		},
		{
			name:        "disabled input (backspace)",
			content:     "abc",
			hasFocus:    true,
			disabled:    true,
			x:           1,
			y:           0,
			event:       tcell.NewEventKey(tcell.KeyBackspace, 0, 0),
			wantHandled: false,
			wantX:       1,
			wantY:       0,
			wantCells:   [][]rune{{'a', 'b', 'c'}},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			i := NewInput(&InputOptions{
				Disabled: tt.disabled,
				InitialValue: tt.content,
			})
			i.SetFocus(tt.hasFocus)
			i.SetCursor(tt.x, tt.y)
			handled := i.HandleEvent(tt.event)
			if handled != tt.wantHandled {
				t.Errorf("Input.HandleEvent() handled = %v, want %v", handled, tt.wantHandled)
			}
			if tt.wantHandled {
				if i.model.x != tt.wantX {
					t.Errorf("Input.HandleEvent() x = %v, want %v", i.model.x, tt.wantX)
				}
				if i.model.y != tt.wantY {
					t.Errorf("Input.HandleEvent() y = %v, want %v", i.model.y, tt.wantY)
				}
				if tt.wantCells != nil && !reflect.DeepEqual(i.model.cells, tt.wantCells) {
					t.Errorf("Input.HandleEvent() cells = %v, want %v", i.model.cells, tt.wantCells)
				}
			}
		})
	}
}

func TestInput_SetHidden(t *testing.T) {
	tests := []struct {
		name           string
		initHidden     bool
		newHidden      bool
		initialContent string
		wantWidth      int
		wantHeight     int
		wantX          int
		wantY          int
	}{
		{
			name:           "set to hidden",
			initHidden:     false,
			newHidden:      true,
			initialContent: "test",
			wantWidth:      PASSWORD_FIELD_LENGTH,
			wantHeight:     1,
			wantX:          0,
			wantY:          0,
		},
		{
			name:           "set to revealed",
			initHidden:     true,
			newHidden:      false,
			initialContent: "test",
			wantWidth:      4,
			wantHeight:     1,
			wantX:          0,
			wantY:          0,
		},
		{
			name:           "set to hidden (multiline)",
			initHidden:     false,
			newHidden:      true,
			initialContent: "test\nline",
			wantWidth:      PASSWORD_FIELD_LENGTH,
			wantHeight:     1,
			wantX:          0,
			wantY:          0,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			i := NewInput(&InputOptions{
				Hidden: tt.initHidden, 
				InitialValue: tt.initialContent,
			})
			i.SetHidden(tt.newHidden)

			if i.model.hidden != tt.newHidden {
				t.Errorf("Input.SetInputType() inputType = %v, want %v", i.model.hidden, tt.newHidden)
			}
			if i.model.width != tt.wantWidth {
				t.Errorf("Input.SetInputType() width = %v, want %v", i.model.width, tt.wantWidth)
			}
			if i.model.height != tt.wantHeight {
				t.Errorf("Input.SetInputType() height = %v, want %v", i.model.height, tt.wantHeight)
			}
			if i.model.x != tt.wantX {
				t.Errorf("Input.SetInputType() x = %v, want %v", i.model.x, tt.wantX)
			}
			if i.model.y != tt.wantY {
				t.Errorf("Input.SetInputType() y = %v, want %v", i.model.y, tt.wantY)
			}
		})
	}
}
