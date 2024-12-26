package autocomplete

import (
	"sync"

	"github.com/gdamore/tcell/v2"
	"github.com/gdamore/tcell/v2/views"
	"github.com/mattn/go-runewidth"
	"github.com/shikaan/keydex/tui/components"
)

// Option

type option struct {
	model *optionModel
	once  sync.Once

	components.Focusable
	views.CellView
}

type optionModel struct {
	// Content as it appears outside
	content string
	// 0-spaced rune sequence, to facilitate printing of Unicode chars
	runes []rune

	width    int
	x        int
	style    tcell.Style
	hasFocus bool

	selectHandler func() bool
	focusHandler  func() bool
}

func (m *optionModel) GetCell(x, y int) (rune, tcell.Style, []rune, int) {
	if x >= len(m.runes) {
		return 0, m.style, nil, 1
	}

	char := m.runes[x]

	if char == 0 {
		return 0, m.style, nil, 1
	}

	return char, m.style, nil, runewidth.RuneWidth(char)
}

func (m *optionModel) GetBounds() (int, int) {
	return m.width, 1
}

func (m *optionModel) SetCursor(x, y int) {
	m.x = 0
}

func (m *optionModel) MoveCursor(x, y int) {
	m.x = 0
}

func (m *optionModel) GetCursor() (int, int, bool, bool) {
	return m.x, 0, true, false
}

func (i *option) Size() (int, int) {
	// Forces height 1, to fix problems on some terminals
	x, _ := i.CellView.Size()
	return x, 1
}

func (i *option) HasFocus() bool {
	return i.model.hasFocus
}

func (i *option) SetFocus(on bool) {
	i.Init()
	if on {
		i.SetContent("> " + i.model.content)
	} else {
		i.SetContent(i.model.content[2:])
	}

	i.model.hasFocus = on
	i.CellView.SetModel(i.model)

	if i.model.focusHandler != nil {
		i.model.focusHandler()
	}
}

func (i *option) SetContent(text string) {
	i.Init()
	m := i.model
	m.width = runewidth.StringWidth(text)
	m.content = text
	m.runes = []rune{}

	for _, c := range text {
		l := runewidth.RuneWidth(c)

		m.runes = append(m.runes, c)
		for i := 1; i < l; i++ {
			m.runes = append(m.runes, 0)
		}
	}

	i.CellView.SetModel(m)
}

func (i *option) GetContent() string {
	return i.model.content
}

func (i *option) HandleEvent(ev tcell.Event) bool {
	if !i.HasFocus() {
		return false
	}

	switch ev := ev.(type) {
	case *tcell.EventKey:
		if ev.Key() == tcell.KeyEnter {
			return i.model.selectHandler != nil && i.model.selectHandler()
		}
	}

	return false
}

func (i *option) OnSelect(cb func() bool) func() {
	i.model.selectHandler = cb
	return func() {
		i.model.selectHandler = nil
	}
}

func (i *option) OnFocus(cb func() bool) func() {
	i.model.focusHandler = cb
	return func() {
		i.model.focusHandler = nil
	}
}

func (i *option) Init() {
	i.once.Do(func() {
		m := &optionModel{}
		i.model = m
		i.CellView.Init()
		i.CellView.SetModel(m)
	})
}

func newOption() *option {
	i := &option{}
	i.Init()
	return i
}
