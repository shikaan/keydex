package components

import (
	"bufio"
	"fmt"
	"strings"
	"sync"
	"unicode"

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
const PAD_BYTE = 0
const PASSWORD_FIELD_LENGTH = 8
const LINE_BREAK = '\n'

type inputModel struct {
	// A read-only string representation of the content for outside to prevent
	// on-the-fly conversions on every frame. Its value is used by other
	// components (e.g., clipboard, kdbx).
	// Again, it's read-only: modifications will not be preserved.
	content string

	// Holds the runes that make up the content of the field. It's structured
	// as a grid to account for multi-line input fields. A rune is accessed by
	// `cell[lineIndex][columnIndex]` (aka `cell[y][x]`).
	//
	// Unicode chars can take more than one cell. To have cells representing
	// correctly the grid on which the cursor moves, we right-pad unicode chars
	// that are larger than a cell with PAD_BYTE.
	//
	// If a char takes two cells, its representation will be [char, PAD_BYTE].
	// For example: "😀" (len 2) is represented as []rune{😀, PAD_BYTE}
	cells [][]rune

	// Display width of the input being displayed
	width int
	// Display height of the input being displayed
	height int
	// X coordinate of the cursor position
	x int
	// Y coordnate of the cursor position
	y int
	// Style. See tcell.Style for details
	style tcell.Style
	// Returns true if the field is focused
	hasFocus bool
	// Whether the field is a password field or a regular one
	inputType InputType

	// Handle a keypress event. Returns true if handled, false if needs cascading
	keyPressHandler func(ev *tcell.EventKey) bool
	// Handle a change event. Returns true if handled, false if needs cascading
	changeHandler func(ev tcell.Event) bool
	// Handle a focus event. Returns true if handled, false if needs cascading
	focusHandler func() bool
}

func (m *inputModel) GetCell(x, y int) (rune, tcell.Style, []rune, int) {
	if m.isOutOfBounds(x, y) {
		return EMPTY_CELL, m.style, nil, 1
	}

	if m.inputType == InputTypePassword {
		return '*', m.style, nil, 1
	}

	char := m.cells[y][x]
	if unicode.IsPrint(char) {
		return char, m.style, nil, runewidth.RuneWidth(char)
	}

	return EMPTY_CELL, m.style, nil, 1
}

func (m *inputModel) GetBounds() (int, int) {
	return m.width, m.height
}

// Prevents the cursor from going out of bounds
func (m *inputModel) limitCursor() {
	m.y = clamp(m.y, 0, len(m.cells)-1)
	if len(m.cells) == 0 {
		m.x = 0
	} else {
		// Not len-1 to allow an extra spot to backspace last char of the line
		m.x = clamp(m.x, 0, len(m.cells[m.y]))
	}
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

// m.cells contains both runes and padding zeros to accommodate rendering.
// This method stably returns the rune at cursor, regardless of the 0s.
// It will however return 0 when cursor is out of bounds
func (m *inputModel) GetRuneAtPosition(x, y int) (rune, int) {
	if m.isOutOfBounds(x, y) {
		return EMPTY_CELL, -1
	}

	if m.inputType == InputTypePassword {
		return '*', x
	}

	for j := x; j >= 0; j-- {
		if m.cells[y][j] != PAD_BYTE {
			return m.cells[y][j], j
		}
	}

	return EMPTY_CELL, -1
}

func (m *inputModel) isOutOfBounds(x, y int) bool {
	if m.inputType == InputTypePassword {
		return x < 0 || x >= PASSWORD_FIELD_LENGTH || y != 0
	}

	return y < 0 || x < 0 || y >= len(m.cells) || x >= len(m.cells[y])
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
	lines := getLines(text)
	m.height = len(lines)
	m.width = 0
	m.cells = make([][]rune, m.height)

	for lineIndex, line := range lines {
		m.cells[lineIndex] = []rune{}
		m.width = max(m.width, runewidth.StringWidth(line))

		for _, char := range line {
			cells := runewidth.RuneWidth(char)

			m.cells[lineIndex] = append(m.cells[lineIndex], char)
			// Pad rune with PAD_BYTE, in case the rune is longer than one cell
			for i := 1; i < cells; i++ {
				m.cells[lineIndex] = append(m.cells[lineIndex], PAD_BYTE)
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
			_, p := i.model.GetRuneAtPosition(i.model.x-1, i.model.y)
			i.model.SetCursor(p, i.model.y)
			return true
		case tcell.KeyRight:
			char, _ := i.model.GetRuneAtPosition(i.model.x, i.model.y)
			i.model.MoveCursor(runewidth.RuneWidth(char), 0)
			return true
		case tcell.KeyDown:
			if i.model.y < i.model.height-1 {
				_, p := i.model.GetRuneAtPosition(i.model.x, i.model.y+1)
				i.model.SetCursor(p, i.model.y+1)
				return true
			}
			return false
		case tcell.KeyUp:
			if i.model.y > 0 {
				_, p := i.model.GetRuneAtPosition(i.model.x, i.model.y-1)
				i.model.SetCursor(p, i.model.y-1)
				return true
			}
			return false
		case tcell.KeyRune:
			return i.handleCellsUpdate(
				ev,
				func() int {
					c, x, y := i.model.cells, i.model.x, i.model.y
					char := ev.Rune()
					c[y] = slices.Insert(c[y], x, char)
					return runewidth.RuneWidth(char)
				},
			)
		case tcell.KeyBackspace2:
			fallthrough
		case tcell.KeyBackspace:
			return i.handleCellsUpdate(
				ev,
				func() int {
					c, x, y := i.model.cells, i.model.x, i.model.y
					char, _ := i.model.GetRuneAtPosition(x-1, y)
					offset := runewidth.RuneWidth(char)
					c[y] = slices.Delete(c[y], x-offset, x)
					return -offset
				},
			)
		case tcell.KeyDelete:
			return i.handleCellsUpdate(
				ev,
				func() int {
					c, x, y := i.model.cells, i.model.x, i.model.y
					char, _ := i.model.GetRuneAtPosition(x, y)
					offset := runewidth.RuneWidth(char)
					c[y] = slices.Delete(c[y], x+offset, x)
					return 0
				},
			)
		}
	}

	return false
}

func (i *Input) handleCellsUpdate(ev tcell.Event, updateCells func() int) bool {
	// Warning: this is order dependent!
	deltaX := updateCells()
	i.SetContent(toString(i.model.cells))
	i.model.MoveCursor(deltaX, 0)

	if i.model.changeHandler != nil {
		return i.model.changeHandler(ev)
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

// Takes a list of cells and returns a string, cleaning up pad bytes
// and using the correct line endings, depending on the platform
func toString(cells [][]rune) string {
	b := &strings.Builder{}
	for lineIndex, line := range cells {
		for _, cell := range line {
			if unicode.IsPrint(cell) {
				b.WriteRune(cell)
			}
		}
		if lineIndex != len(cells)-1 {
			fmt.Fprintln(b, "") // Cross-platform way of adding a line-ending
		}
	}
	return b.String()
}

// Clamps an integer between minVale and maxValue (both included)
func clamp(n, minValue, maxValue int) int {
	return max(min(n, maxValue), minValue)
}

// Breaks a text in lines in a platform independent way
func getLines(text string) []string {
	reader := strings.NewReader(text)
	scanner := bufio.NewScanner(reader)
	lines := []string{}
	for scanner.Scan() {
		lines = append(lines, scanner.Text())
	}
	return lines
}
