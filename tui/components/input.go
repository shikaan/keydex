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
const LINE_BREAK = '\n'

type inputModel struct {
	// This value is used only for caching purposes. It's the content as exposed outside,
	// but all the actual operations on the values need to be done on cells
	content string
	// Unicode chars can take more than one cell.
	// If a char takes two cells, its representation will be [char, 0].
	// For example: "😀" (len 2) is represented as []rune{😀, 0}
	cells [][]rune

	width     int
	height    int
	x         int
	y         int
	style     tcell.Style
	hasFocus  bool
	inputType InputType

	keyPressHandler func(ev *tcell.EventKey) bool
	changeHandler   func(ev tcell.Event) bool
	focusHandler    func() bool
}

func (m *inputModel) GetCell(x, y int) (rune, tcell.Style, []rune, int) {
	isPassword := m.inputType == InputTypePassword
	isOutOfBoundY := y < 0 || (!isPassword && y >= len(m.cells)) || (isPassword && y >= 1)
	isOutOfBound := isOutOfBoundY || x < 0 || (!isPassword && x >= len(m.cells[y])) || (isPassword && x >= PASSWORD_FIELD_LENGTH)

	if isOutOfBound {
		return EMPTY_CELL, m.style, nil, 1
	}

	if isPassword {
		return '*', m.style, nil, 1
	}

	char := m.cells[y][x]
	if isEmptyCell(char) {
		return EMPTY_CELL, m.style, nil, 1
	}

	return char, m.style, nil, runewidth.RuneWidth(char)
}

func (m *inputModel) GetBounds() (int, int) {
	return m.width, m.height
}

// Prevents the cursor from going out of bounds
func (m *inputModel) limitCursor() {
	m.y = clamp(m.y, 0, len(m.cells))
	// Not len-1 to allow an extra spot to backspace last char of the line
	m.x = clamp(m.x, 0, len(m.cells[m.y]))
}

func (m *inputModel) SetCursor(x, y int) {
	m.x = x
	m.y = y
	m.limitCursor()
}

func (m *inputModel) MoveCursor(x, y int) {
	m.x += x
	m.y += y
	m.limitCursor()
}

func (m *inputModel) GetCursor() (int, int, bool, bool) {
	return m.x, m.y, true, m.hasFocus
}

// m.cells contains both runes and placeholder chars (0) to accommodate rendering.
// This method stably returns the rune at cursor, regardless of the 0s.
// It will however return 0 when cursor is out of bounds
func (m *inputModel) FindRuneAtPosition(x, y int) (rune, int) {
	lines := len(m.cells)

	if y < 0 || x < 0 || lines == 0 || y >= lines || x >= len(m.cells[0]) {
		return EMPTY_CELL, -1
	}

	for j := x; j >= 0; j-- {
		if !isEmptyCell(m.cells[y][j]) {
			return m.cells[y][j], j
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
	m.content = text
	lines := strings.Split(text, "\n") // TODO: is this OS independent?
	m.height = len(lines)
	m.width = 0
	m.cells = make([][]rune, m.height)

	for line, text := range lines {
		m.cells[line] = []rune{}
		m.width = max(m.width, runewidth.StringWidth(text))

		for _, char := range text {
			cells := runewidth.RuneWidth(char)

			m.cells[line] = append(m.cells[line], char)
			// Pad rune with 0 cells, in case the rune is longer than one cell
			for i := 1; i < cells; i++ {
				m.cells[line] = append(m.cells[line], 0)
			}
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
	i.model.y = 0
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

		// Don't allow interactions with password fields when hidden,
		if i.model.inputType == InputTypePassword {
			return false
		}

		switch ev.Key() {
		case tcell.KeyLeft:
			_, p := i.model.FindRuneAtPosition(i.model.x-1, i.model.y)
			i.model.SetCursor(p, i.model.y)
			return true
		case tcell.KeyRight:
			char, _ := i.model.FindRuneAtPosition(i.model.x, i.model.y)
			i.model.MoveCursor(runewidth.RuneWidth(char), 0)
			return true
		case tcell.KeyDown:
			if i.model.y < i.model.height-1 {
				i.model.MoveCursor(0, 1)
				return true
			}
			return false
		case tcell.KeyUp:
			if i.model.y > 0 {
				i.model.MoveCursor(0, -1)
				return true
			}
			return false
		case tcell.KeyRune:
			return i.handleContentUpdate(
				ev,
				func(c [][]rune, x int, y int) ([][]rune, int) {
					char := ev.Rune()
					c[y] = slices.Insert(c[y], x, char)
					return c, runewidth.RuneWidth(char)
				},
			)
		case tcell.KeyBackspace2:
			fallthrough
		case tcell.KeyBackspace:
			return i.handleContentUpdate(
				ev,
				func(c [][]rune, x, y int) ([][]rune, int) {
					char, _ := i.model.FindRuneAtPosition(x-1, y)
					offset := runewidth.RuneWidth(char)
					c[y] = slices.Delete(c[y], x-offset, x)
					return c, -offset
				},
			)
		case tcell.KeyDelete:
			return i.handleContentUpdate(
				ev,
				func(c [][]rune, x, y int) ([][]rune, int) {
					char, _ := i.model.FindRuneAtPosition(x, y)
					offset := runewidth.RuneWidth(char)
					c[y] = slices.Delete(c[y], x+offset, x)
					return c, 0
				},
			)
		}
	}

	return false
}

func (i *Input) handleContentUpdate(ev tcell.Event, cb func(initialCells [][]rune, x int, y int) (cells [][]rune, width int)) bool {
	cells, offsetX := cb(i.model.cells, i.model.x, i.model.y)

	i.SetContent(toString(cells))
	i.model.MoveCursor(offsetX, 0)

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
func toString(cells [][]rune) string {
	b := &strings.Builder{}
	for lineIndex, line := range cells {
		for _, cell := range line {
			if cell != EMPTY_CELL {
				b.WriteRune(cell)
			}
		}
		if lineIndex != len(cells)-1 {
			b.WriteRune(LINE_BREAK)
		}
	}
	return b.String()
}

func isEmptyCell(c rune) bool {
	return c == EMPTY_CELL
}

func clamp(n, minValue, maxValue int) int {
	return max(min(n, maxValue), minValue)
}
