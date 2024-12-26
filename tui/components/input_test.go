package components

import (
	"reflect"
	"testing"

	"github.com/gdamore/tcell/v2"
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
		{"EMPTY_CELL when out of bounds", inputModel{cells: [][]rune{{'L', 'O', 'L'}}}, 2, 2, EMPTY_CELL, 1},
		{"EMPTY_CELL when cursor out of password bounds (password)",
			inputModel{cells: [][]rune{{'a', 'a', 'a', 'a', 'a', 'a', 'a', 'a', 'a', 'a'}}, inputType: InputTypePassword},
			PASSWORD_FIELD_LENGTH + 1, 0, // This is longer than a password, but shorter than field itself
			EMPTY_CELL, 1},
		{"EMPTY_CELL when cursor out of field bounds (password)",
			inputModel{cells: [][]rune{{'a', 'a', 'a', 'a', 'a', 'a', 'a', 'a', 'a', 'a'}}, inputType: InputTypePassword},
			14, 0, // This is longer than the field
			EMPTY_CELL, 1},
		{"* when password", inputModel{cells: [][]rune{{'L', 'O', 'L'}}, inputType: InputTypePassword}, 1, 0, '*', 1},
		{"* when non-byte password", inputModel{cells: [][]rune{{'🤖', '✅', '😂'}}, inputType: InputTypePassword}, 1, 0, '*', 1},
		{"returns byte with size 1 (only byte)", inputModel{cells: [][]rune{{'L', 'O', 'L'}}}, 1, 0, 'O', 1},
		{"returns emoji with size 2 (only emoji)", inputModel{cells: [][]rune{{'🤖', '✅', '😂'}}}, 1, 0, '✅', 2},
		{"returns byte with size 1 (mixed)", inputModel{cells: [][]rune{{'I', '🤖'}}}, 0, 0, 'I', 1},
		{"returns emoji with size 2 (mixed)", inputModel{cells: [][]rune{{'I', '🤖'}}}, 1, 0, '🤖', 2},
		{"returns emoji with size 2 (mixed)", inputModel{cells: [][]rune{{'I', '🤖'}}}, 1, 0, '🤖', 2},
		{"EMPTY_CELL with non-print characters", inputModel{cells: [][]rune{{'I', '\a'}}}, 1, 0, EMPTY_CELL, 1},
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
				inputType:       tt.fields.inputType,
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
				inputType:       tt.fields.inputType,
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
				inputType:       tt.fields.inputType,
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
		{"get * with password", inputModel{cells: [][]rune{{'T', 'e'}}, inputType: InputTypePassword}, 1, 0, '*', 1},
		{"get byte char", inputModel{cells: [][]rune{{'T', 'e', 's'}}}, 1, 0, 'e', 1},
		{"get unicode char", inputModel{cells: [][]rune{{'T', '✅', PAD_BYTE, 's'}}}, 1, 0, '✅', 1},
		{"get unicode char on PAD_BYTE", inputModel{cells: [][]rune{{'T', '✅', PAD_BYTE, 's'}}}, 2, 0, '✅', 1},
		{"EMPTY_CELL when out of bounds", inputModel{cells: [][]rune{{'T', 's'}}}, 2, 0, EMPTY_CELL, -1},
		{"EMPTY_CELL when out of bounds (password)", inputModel{inputType: InputTypePassword, cells: [][]rune{{'T', 's'}}}, PASSWORD_FIELD_LENGTH + 1, 0, EMPTY_CELL, -1},
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
				inputType:       tt.fields.inputType,
				keyPressHandler: tt.fields.keyPressHandler,
				changeHandler:   tt.fields.changeHandler,
				focusHandler:    tt.fields.focusHandler,
			}
			gotRune, gotHOffset := m.GetRuneAtPosition(tt.x, tt.y)
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
		{"non-ASCII only", "🤖ß", [][]rune{{'🤖', PAD_BYTE, 'ß'}}, "🤖ß"},
		{"mixed input", "I💙ßs", [][]rune{{'I', '💙', PAD_BYTE, 'ß', 's'}}, "I💙ßs"},
		{"line breaks UNIX", "I💙ßs\n&🥔", [][]rune{{'I', '💙', PAD_BYTE, 'ß', 's'}, {'&', '🥔', PAD_BYTE}}, "I💙ßs\n&🥔"},
		{"line breaks Windows", "I💙ßs\r\n&🥔", [][]rune{{'I', '💙', PAD_BYTE, 'ß', 's'}, {'&', '🥔', PAD_BYTE}}, "I💙ßs\r\n&🥔"},
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
		name   string
		fields Input
		event  tcell.Event
		want   bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			i := &Input{
				model:     tt.fields.model,
				Focusable: tt.fields.Focusable,
			}
			if got := i.HandleEvent(tt.event); got != tt.want {
				t.Errorf("Input.HandleEvent() = %v, want %v", got, tt.want)
			}
		})
	}
}
