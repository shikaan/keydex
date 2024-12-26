package autocomplete

import (
	"strings"
	"sync"
	"unicode"

	"github.com/gdamore/tcell/v2"
	"github.com/gdamore/tcell/v2/views"
	"github.com/mattn/go-runewidth"
	"github.com/shikaan/keydex/tui/components"
	"golang.org/x/exp/slices"
)

type Search struct {
	model *searchModel
	once  sync.Once

	components.Focusable
	views.CellView
}

// Very similar to the input model, but in the autocomplete context
// this behaves differently. It's been extracted because they are evolving
// in different ways.
type searchModel struct {
	content       string
	cells         components.PaddedLine
	width         int
	x             int
	style         tcell.Style
	hasFocus      bool
	changeHandler func(ev tcell.Event) bool
	focusHandler  func() bool
}

func (m *searchModel) GetCell(x, y int) (rune, tcell.Style, []rune, int) {
	if y != 0 || x < 0 || x >= len(m.cells) {
		return components.EMPTY_CELL, m.style, nil, 1
	}

	char := m.cells[x]
	if unicode.IsPrint(char) {
		return char, m.style, nil, runewidth.RuneWidth(char)
	}

	return components.EMPTY_CELL, m.style, nil, 1
}

func (m *searchModel) GetBounds() (int, int) {
	return m.width, 1
}

func (m *searchModel) SetCursor(x, y int) {
	m.x = max(min(x+m.x, len(m.cells)), 0)
}

func (m *searchModel) MoveCursor(x, y int) {
	m.x = max(min(x+m.x, len(m.cells)), 0)
}

func (m *searchModel) GetCursor() (int, int, bool, bool) {
	return m.x, 0, true, true
}

func (s *Search) OnChange(cb func(ev tcell.Event) bool) func() {
	s.model.changeHandler = cb
	return func() {
		s.model.changeHandler = nil
	}
}

func (s *Search) OnFocus(cb func() bool) func() {
	s.model.focusHandler = cb
	return func() {
		s.model.focusHandler = nil
	}
}

func (s *Search) HasFocus() bool {
	return s.model.hasFocus
}

func (s *Search) SetFocus(on bool) {
	s.Init()
	s.model.hasFocus = on
	s.CellView.SetModel(s.model)
	if s.model.focusHandler != nil {
		s.model.focusHandler()
	}
}

func (s *Search) SetContent(text string) {
	s.Init()
	m := s.model
	m.width = runewidth.StringWidth(text)
	m.content = text
	m.cells = components.NewPaddedLine(text)

	s.CellView.SetModel(m)
}

func (s *Search) GetContent() string {
	return s.model.content
}

func (s *Search) HandleEvent(ev tcell.Event) bool {
	if !s.HasFocus() {
		return false
	}

	switch ev := ev.(type) {
	case *tcell.EventKey:
		switch ev.Key() {
		case tcell.KeyRune:
			return s.handleContentUpdate(ev, func() int {
				char := ev.Rune()
				s.model.cells = slices.Insert(s.model.cells, s.model.x, char)
				return runewidth.RuneWidth(char)
			})
		case tcell.KeyBackspace2:
			fallthrough
		case tcell.KeyBackspace:
			return s.handleContentUpdate(ev, func() int {
				char, _ := components.GetRune(s.model.cells, s.model.x-1)
				offset := runewidth.RuneWidth(char)
				s.model.cells = slices.Delete(s.model.cells, s.model.x-offset, s.model.x)
				return -offset
			})
		}
	}

	return false
}

func (s *Search) handleContentUpdate(ev tcell.Event, cb func() int) bool {
	offset := cb()
	s.SetContent(toString(s.model.cells))
	s.model.MoveCursor(offset, 0)

	if s.model.changeHandler != nil {
		return s.model.changeHandler(ev)
	}

	return true
}

func (i *Search) Init() {
	i.once.Do(func() {
		i.model = &searchModel{}
		i.CellView.Init()
		i.CellView.SetModel(i.model)
	})
}

func NewSearch() *Search {
	s := &Search{}
	s.Init()
	return s
}

// Takes a list of cells and returns a string, filtering out empty cells
func toString(cells []rune) string {
	b := &strings.Builder{}
	for _, cell := range cells {
		if cell != components.PAD_BYTE {
			b.WriteRune(cell)
		}
	}
	return b.String()
}
