package components

import (
	"strings"
	"sync"

	"github.com/gdamore/tcell/v2"
	"github.com/gdamore/tcell/v2/views"
	"github.com/mattn/go-runewidth"
	"golang.org/x/exp/slices"
)

type Input struct {
	model *inputModel
	once  sync.Once

	Focusable
	views.CellView
}

type InputOptions struct {
	InitialValue string
	Type         InputType
}

type InputType int

const (
	InputTypeText InputType = iota
	InputTypePassword
)

const EMPTY_CELL = 0
const PASSWORD_FIELD_LENGTH = 8

type inputModel struct {
	// This value is used only for caching purposes. It's the content as exposed outside,
	// but all the actual operations on the values need to be done on cells
	content string
	// Unicode chars can take more than one cell.
	// If a char takes two cells, its representation will be [char, 0].
	// For example: "😀" (len 2) is represented as []rune{😀, 0}
	cells []rune

	width     int
	x         int
	style     tcell.Style
	hasFocus  bool
	inputType InputType

	keyPressHandler func(ev *tcell.EventKey) bool
	changeHandler   func(ev tcell.Event) bool
	focusHandler    func() bool
}

func (m *inputModel) GetCell(x, y int) (rune, tcell.Style, []rune, int) {
	isPassword := m.inputType == InputTypePassword
	isOutOfBound := y != 0 || x < 0 || (!isPassword && x >= len(m.cells)) || (isPassword && x >= PASSWORD_FIELD_LENGTH)

	if isOutOfBound {
		return EMPTY_CELL, m.style, nil, 1
	}

	if isPassword {
		return '*', m.style, nil, 1
	}

	char := m.cells[x]
	if char == EMPTY_CELL {
		return EMPTY_CELL, m.style, nil, 1
	}

	return char, m.style, nil, runewidth.RuneWidth(char)
}

func (m *inputModel) GetBounds() (int, int) {
	return m.width, 1
}

// Prevents the cursor from going out of bounds
func (m *inputModel) limitCursor() {
	x, _ := m.GetBounds()

	if m.x > x {
		m.x = x
	}

	if m.x < 0 {
		m.x = 0
	}
}

func (m *inputModel) SetCursor(x, y int) {
	m.x = x
	m.limitCursor()
}

func (m *inputModel) MoveCursor(x, y int) {
	m.x += x
	m.limitCursor()
}

func (m *inputModel) GetCursor() (int, int, bool, bool) {
	return m.x, 0, true, m.hasFocus
}

// m.cells contains both runes and placeholder chars (0) to accommodate rendering.
// This method stably returns the rune at cursor, regardless of the 0s.
// It will however return 0 when cursor is out of bounds
func (m *inputModel) FindRuneAtPosition(cursorPosition int) (rune, int) {
	if cursorPosition < 0 || cursorPosition >= len(m.cells) {
		return EMPTY_CELL, -1
	}

	for j := cursorPosition; j >= 0; j-- {
		if m.cells[j] != 0 {
			return m.cells[j], j
		}
	}

	return EMPTY_CELL, -1
}

func (i *Input) HasFocus() bool {
	return i.model.hasFocus
}

func (i *Input) SetFocus(on bool) {
	i.Init()
	i.model.hasFocus = on
	i.CellView.SetModel(i.model)
	if i.model.focusHandler != nil {
		i.model.focusHandler()
	}
}

func (i *Input) SetContent(text string) {
	i.Init()
	m := i.model
	m.width = runewidth.StringWidth(text)
	m.content = text

	// Pad rune with 0 cells, in case the rune is longer than one cell
	m.cells = []rune{}
	for _, char := range text {
		cells := runewidth.RuneWidth(char)

		m.cells = append(m.cells, char)
		for i := 1; i < cells; i++ {
			m.cells = append(m.cells, 0)
		}
	}

	i.CellView.SetModel(m)
}

func (i *Input) GetContent() string {
	return i.model.content
}

func (i *Input) SetInputType(t InputType) {
	i.model.inputType = t
	i.model.x = 0
	i.Init()
}

func (i *Input) GetInputType() InputType {
	return i.model.inputType
}

func (i *Input) HandleEvent(ev tcell.Event) bool {
	if !i.HasFocus() {
		return false
	}

	switch ev := ev.(type) {
	case *tcell.EventKey:
		handled := false

		if i.model.keyPressHandler != nil {
			handled = i.model.keyPressHandler(ev)
		}

		if handled {
			return handled
		}

		// Don't allow interactions with password fields when hidden
		if i.model.inputType == InputTypePassword {
			return true
		}

		switch ev.Key() {
		case tcell.KeyLeft:
			_, p := i.model.FindRuneAtPosition(i.model.x - 1)
			i.model.SetCursor(p, 0)
			return true
		case tcell.KeyRight:
			char, _ := i.model.FindRuneAtPosition(i.model.x)
			i.model.MoveCursor(runewidth.RuneWidth(char), 0)
			return true
		case tcell.KeyRune:
			return i.handleContentUpdate(
				ev,
				func(c []rune, x int) ([]rune, int) {
					char := ev.Rune()
					return slices.Insert(c, x, char), runewidth.RuneWidth(char)
				},
			)
		case tcell.KeyBackspace2:
			fallthrough
		case tcell.KeyBackspace:
			return i.handleContentUpdate(
				ev,
				func(c []rune, x int) ([]rune, int) {
					char, _ := i.model.FindRuneAtPosition(x - 1)
					offset := runewidth.RuneWidth(char)
					return slices.Delete(c, x-offset, x), -offset
				},
			)
		case tcell.KeyDelete:
			return i.handleContentUpdate(
				ev,
				func(c []rune, x int) ([]rune, int) {
					char, _ := i.model.FindRuneAtPosition(x)
					offset := runewidth.RuneWidth(char)
					return slices.Delete(c, x, x+offset), 0
				},
			)
		}
	}

	return false
}

func (i *Input) handleContentUpdate(ev tcell.Event, cb func([]rune, int) ([]rune, int)) bool {
	cells, offset := cb(i.model.cells, i.model.x)

	i.SetContent(toString(cells))
	i.model.MoveCursor(offset, 0)

	if i.model.changeHandler != nil {
		i.model.changeHandler(ev)
	}

	return true
}

func (i *Input) OnKeyPress(cb func(ev *tcell.EventKey) bool) func() {
	i.model.keyPressHandler = cb
	return func() {
		i.model.keyPressHandler = nil
	}
}

func (i *Input) OnChange(cb func(ev tcell.Event) bool) func() {
	i.model.changeHandler = cb
	return func() {
		i.model.changeHandler = nil
	}
}

func (i *Input) OnFocus(cb func() bool) func() {
	i.model.focusHandler = cb
	return func() {
		i.model.focusHandler = nil
	}
}

func (i *Input) Init() {
	i.once.Do(func() {
		m := &inputModel{}
		i.model = m
		i.CellView.Init()
		i.CellView.SetModel(m)
	})
}

func NewInput(options *InputOptions) *Input {
	i := &Input{}
	i.Init()
	i.model.inputType = options.Type

	return i
}

// Takes a list of cells and returns a string, filtering out empty cells
func toString(cells []rune) string {
	b := &strings.Builder{}
	for _, cell := range cells {
		if cell != EMPTY_CELL {
			b.WriteRune(cell)
		}
	}
	return b.String()
}
