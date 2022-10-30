package form

import (
	"strings"
	"sync"

	"github.com/gdamore/tcell/v2"
	"github.com/gdamore/tcell/v2/views"
)

// TextArea is a pannable 2 dimensional text widget. It wraps both
// the view and the model for text in a single, convenient widget.
// Text is provided as an array of strings, each of which represents
// a single line to display.  All text in the TextArea has the same
// style.  An optional soft cursor is available.
type TextArea struct {
	model *linesModel
	once  sync.Once
	views.CellView
}

type linesModel struct {
	runes  []rune
	width  int
	height int
	x      int
	y      int
	hide   bool
	cursor bool
	style  tcell.Style
}

func (m *linesModel) GetCell(x int) (rune, tcell.Style, []rune, int) {
	if x < 0 || x >= len(m.runes) {
		return 0, m.style, nil, 1
	}
	// XXX: extend this to support combining and full width chars
	return m.runes[x], m.style, nil, 1
}

func (m *linesModel) GetBounds() (int, int) {
	return m.width, m.height
}

func (m *linesModel) limitCursor() {
	if m.x > m.width-1 {
		m.x = m.width - 1
	}
	if m.y > m.height-1 {
		m.y = m.height - 1
	}
	if m.x < 0 {
		m.x = 0
	}
	if m.y < 0 {
		m.y = 0
	}
}

func (m *linesModel) SetCursor(x, y int) {
	m.x = x
	m.y = y
	m.limitCursor()
}

func (m *linesModel) MoveCursor(x, y int) {
	m.x += x
	m.y += y
	m.limitCursor()
}

func (m *linesModel) GetCursor() (int, int, bool, bool) {
	return m.x, m.y, m.cursor, !m.hide
}

func (ta *TextArea) SetStyle(style tcell.Style) {
	ta.model.style = style
	ta.CellView.SetStyle(style)
}

// EnableCursor enables a soft cursor in the TextArea.
func (ta *TextArea) EnableCursor(on bool) {
	ta.Init()
	ta.model.cursor = on
}

// HideCursor hides or shows the cursor in the TextArea.
// If on is true, the cursor is hidden.  Note that a cursor is only
// shown if it is enabled.
func (ta *TextArea) HideCursor(on bool) {
	ta.Init()
	ta.model.hide = on
}

// SetContent is used to set the textual content, passed as a
// single string.  Lines within the string are delimited by newlines.
func (ta *TextArea) SetContent(text string) {
	ta.Init()
	lines := strings.Split(strings.Trim(text, "\n"), "\n")
	ta.SetLines(lines)
}

// Init initializes the TextArea.
func (ta *TextArea) Init() {
	ta.once.Do(func() {
		lm := &linesModel{runes: [][]rune{}, width: 0}
		ta.model = lm
		ta.CellView.Init()
		ta.CellView.SetModel(lm)
	})
}

// NewTextArea creates a blank TextArea.
func NewTextArea() *TextArea {
	ta := &TextArea{}
	ta.Init()
	return ta
}

