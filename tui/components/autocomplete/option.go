package autocomplete

import (
	"sync"
	"unicode"

	"github.com/gdamore/tcell/v2"
	"github.com/gdamore/tcell/v2/views"
	"github.com/mattn/go-runewidth"
	"github.com/shikaan/keydex/tui/components"
	"github.com/shikaan/keydex/tui/components/line"
)

type Option struct {
	model *optionModel
	once  sync.Once

	components.Focusable
	views.CellView
}

type optionModel struct {
	content       string
	runes         line.PaddedLine
	width         int
	style         tcell.Style
	hasFocus      bool
	selectHandler func() bool
	focusHandler  func() bool
}

func (m *optionModel) GetCell(x, y int) (rune, tcell.Style, []rune, int) {
	if y != 0 || x < 0 || x >= len(m.runes) {
		return line.EMPTY_CELL, m.style, nil, 1
	}

	if char := m.runes[x]; unicode.IsPrint(char) {
		return char, m.style, nil, runewidth.RuneWidth(char)
	}

	return line.EMPTY_CELL, m.style, nil, 1
}

func (m *optionModel) GetBounds() (int, int) {
	return m.width, 1
}

func (m *optionModel) SetCursor(x, y int)  {}
func (m *optionModel) MoveCursor(x, y int) {}

func (m *optionModel) GetCursor() (int, int, bool, bool) {
	return 0, 0, true, false
}

func (i *Option) Size() (int, int) {
	// Forces height 1, to fix problems on some terminals
	x, _ := i.CellView.Size()
	return x, 1
}

func (i *Option) HasFocus() bool {
	return i.model.hasFocus
}

func (i *Option) SetFocus(on bool) {
	if i.model.hasFocus == on {
		return
	}

	i.Init()
	if on {
		i.SetContent("> " + i.model.content)
	} else {
		if len(i.model.content) >= 2 {
			i.SetContent(i.model.content[2:])
		}
	}

	i.model.hasFocus = on
	i.CellView.SetModel(i.model)

	if i.model.focusHandler != nil {
		i.model.focusHandler()
	}
}

func (i *Option) SetContent(text string) {
	i.Init()
	i.model.width = runewidth.StringWidth(text)
	i.model.content = text
	i.model.runes = line.NewPaddedLine(text)
	i.CellView.SetModel(i.model)
}

func (i *Option) GetContent() string {
	if i.model.hasFocus {
		return i.model.content[2:] // skip caret
	}
	return i.model.content
}

func (i *Option) HandleEvent(ev tcell.Event) bool {
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

func (i *Option) OnSelect(cb func() bool) func() {
	i.model.selectHandler = cb
	return func() {
		i.model.selectHandler = nil
	}
}

func (i *Option) OnFocus(cb func() bool) func() {
	i.model.focusHandler = cb
	return func() {
		i.model.focusHandler = nil
	}
}

func (i *Option) Init() {
	i.once.Do(func() {
		i.model = &optionModel{}
		i.CellView.Init()
		i.CellView.SetModel(i.model)
	})
}

func newOption() *Option {
	i := &Option{}
	i.Init()
	return i
}
