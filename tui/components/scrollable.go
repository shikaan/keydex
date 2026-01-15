package components

import (
	"strings"
	"sync"

	"github.com/gdamore/tcell/v2"
	"github.com/gdamore/tcell/v2/views"
)

// Scrollable is a scrollable text widget with fixed width and height.
// It displays only the visible portion of content based on scroll offset.
type Scrollable struct {
	model *scrollableModel
	once  sync.Once
	views.CellView
}

type scrollableModel struct {
	runes   [][]rune
	width   int
	height  int
	offsetY int
}

func (m *scrollableModel) GetCell(x, y int) (rune, tcell.Style, []rune, int) {
	// y is in viewport coordinates (0 to fixedHeight-1)
	// Translate to content coordinates
	contentY := y + m.offsetY

	// Check if we're within content bounds
	if contentY < 0 || contentY >= len(m.runes) {
		return 0, tcell.StyleDefault, nil, 1
	}

	line := m.runes[contentY]
	if x < 0 || x >= len(line) {
		return 0, tcell.StyleDefault, nil, 1
	}

	return line[x], tcell.StyleDefault, nil, 1
}

func (m *scrollableModel) GetBounds() (int, int) {
	// Return viewport dimensions
	return m.width, m.height
}

func (m *scrollableModel) GetCursor() (int, int, bool, bool) {
	// No cursor
	return 0, 0, false, false
}

func (m *scrollableModel) SetCursor(x, y int) {
	// No-op: no cursor support
}

func (m *scrollableModel) MoveCursor(x, y int) {
	// No-op: no cursor support
}

// SetLines sets the content to display.
func (s *Scrollable) SetLines(lines []string) {
	s.Init(s.model.height, s.model.width)
	runes := make([][]rune, 0, len(lines))
	for _, l := range lines {
		runes = append(runes, []rune(l))
	}
	s.model.runes = runes
	s.model.offsetY = 0
	s.CellView.SetModel(s.model)
}

// SetContent sets the content from a single string with newlines.
func (s *Scrollable) SetContent(text string) {
	s.Init(s.model.height, s.model.width)
	lines := strings.Split(strings.Trim(text, "\n"), "\n")
	s.SetLines(lines)
}

// Init initializes the Scrollable.
func (s *Scrollable) Init(width, height int) {
	s.once.Do(func() {
		m := &scrollableModel{
			runes:   [][]rune{},
			width:   width,
			height:  height,
			offsetY: 0,
		}
		s.model = m
		s.CellView.Init()
		s.CellView.SetModel(m)
	})
}

// NewScrollable creates a new Scrollable with fixed width and height.
func NewScrollable(width, height int) *Scrollable {
	s := &Scrollable{}
	s.Init(width, height)
	return s
}

// HandleEvent handles key events for scrolling
func (s *Scrollable) HandleEvent(ev tcell.Event) bool {
	switch ev := ev.(type) {
	case *tcell.EventKey:
		if s.model.height >= len(s.model.runes) {
			return true
		}

		switch ev.Key() {
		case tcell.KeyUp:
			s.model.offsetY = max(s.model.offsetY-1, 0)
			return true
		case tcell.KeyDown:
			s.model.offsetY = min(s.model.offsetY+1, max(0, len(s.model.runes)-s.model.height))
			return true
		}
	}

	return s.CellView.HandleEvent(ev)
}

func (s *Scrollable) Resize() {
	s.CellView.Resize()
}

func (s *Scrollable) SetSize(width, height int) {
	s.model.height = height
	s.model.width = width
}
