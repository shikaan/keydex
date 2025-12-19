package components

import (
	"strings"
	"sync"

	"github.com/gdamore/tcell/v2"
	"github.com/gdamore/tcell/v2/views"
	"github.com/mattn/go-runewidth"
)

type Text struct {
	model  *textModel
	once   sync.Once
	screen tcell.Screen

	views.CellView
}

type textModel struct {
	runes   [][]rune
	content string
	width   int
	pad     int
	height  int
	style   tcell.Style
}

func (m *textModel) GetBounds() (int, int) {
	return m.width, m.height
}

func (m *textModel) GetCell(x, y int) (rune, tcell.Style, []rune, int) {
	isOutOfBound := y < 0 || x < 0 || y >= len(m.runes) || x >= len(m.runes[y])

	if isOutOfBound {
		return 0, m.style, nil, 1
	}

	return m.runes[y][x], tcell.StyleDefault, nil, 1
}

func (m *textModel) GetCursor() (x, y int, hidden, active bool) {
	return 0, 0, false, false
}

func (m *textModel) MoveCursor(x, y int) {}
func (m *textModel) SetCursor(x, y int)  {}

func (t *Text) SetContent(content string) {
	m := t.model
	m.content = content
	t.Init(m.pad)

	result := [][]rune{}
	for line := range strings.SplitSeq(content, "\n") {
		result = append(result, chunk(line, m.width, m.pad)...)
	}

	m.runes = result
	m.height = len(m.runes)

	t.CellView.SetModel(m)
}

func chunk(s string, length, pad int) [][]rune {
	spacer := strings.Repeat(" ", pad)
	// This accounts for padding left and right
	lengthWithSpacer := (length - pad - pad)

	// This happens on the first render, when we still have no screen
	if lengthWithSpacer < 0 {
		return [][]rune{}
	}

	aux := s
	result := [][]rune{}

	for runewidth.StringWidth(aux) > lengthWithSpacer {
		line := aux[:lengthWithSpacer]
		lastSpace := strings.LastIndex(line, " ") + 1

		if lastSpace == 0 {
			lastSpace = lengthWithSpacer
		}

		result = append(result, []rune(spacer+aux[:lastSpace]))
		aux = aux[lastSpace:]
	}

	result = append(result, []rune(spacer+aux))

	return result
}

func (t *Text) Init(pad int) {
	t.once.Do(func() {
		m := &textModel{pad: pad}
		t.model = m
		t.CellView.Init()
		t.CellView.SetModel(m)
	})
}

func (t *Text) HandleEvent(ev tcell.Event) bool {
	switch ev.(type) {
	case *views.EventWidgetResize:
		m := t.model
		x, _ := t.screen.Size()
		m.width = x
		t.SetContent(m.content)
		return true
	}

	return t.CellView.HandleEvent(ev)
}

func NewText(screen tcell.Screen, pad int) *Text {
	t := &Text{screen: screen}
	t.Init(pad)

	return t
}
