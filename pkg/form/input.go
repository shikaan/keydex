package form

import (
	"sync"

	"github.com/gdamore/tcell/v2"
	"github.com/gdamore/tcell/v2/views"
)

type Input struct {
	model *linesModel
	once  sync.Once
	views.CellView
}

type linesModel struct {
	runes  []rune
	width  int
	x      int
	hide   bool
	cursor bool
	style  tcell.Style
}

func (m *linesModel) GetCell(x, y int) (rune, tcell.Style, []rune, int) {
	if x < 0 || x >= len(m.runes) {
		return 0, m.style, nil, 1
	}

	return m.runes[x], m.style, nil, 1
}

func (m *linesModel) GetBounds() (int, int) {
	return m.width, 1 
}

func (m *linesModel) limitCursor() {
	if m.x > m.width-1 {
		m.x = m.width - 1
	}
	if m.x < 0 {
		m.x = 0
	}
}

func (m *linesModel) SetCursor(x, y int) {
	m.x = x
	m.limitCursor()
}

func (m *linesModel) MoveCursor(x, y int) {
	m.x += x
	m.limitCursor()
}

func (m *linesModel) GetCursor() (int, int, bool, bool) {
	return m.x, 0, true, !m.hide
}

func (i *Input) SetStyle(style tcell.Style) {
	i.model.style = style
	i.CellView.SetStyle(style)
}

func (i *Input) HideCursor(on bool) {
	i.Init()
	i.model.hide = on
}

func (i *Input) SetContent(text string) {	
	i.Init()
	m := i.model
	m.width = len(text)
	m.runes = []rune(text)

  i.CellView.SetModel(m)
}

func (i *Input) GetContent() []rune {
  return i.model.runes
}


func (i *Input) HandleEvent(ev tcell.Event) bool {
  switch ev := ev.(type) {	
  case *tcell.EventKey:
		if ev.Key() == tcell.KeyRune {
		  i.typeAtCursor(ev.Rune())	
      return true
		}
	}
	return i.CellView.HandleEvent(ev)
}

func (i *Input) typeAtCursor(char rune) {
  m := i.GetModel()
  x, _, on, visible := m.GetCursor()

  if on && visible {
    c := i.GetContent()
    previus := c[:x]
    next := c[x:]
    newContent := append(previus, char)
    newContent = append(newContent, next...)
    i.SetContent(string(newContent))
    m.MoveCursor(1, 0)
  }
}

func (i *Input) Init() {
	i.once.Do(func() {
    lm := &linesModel{runes: []rune{}, width: 0, cursor: true}
		i.model = lm
		i.CellView.Init()
		i.CellView.SetModel(lm)
	})
}

func NewInput() *Input {
	i := &Input{}
	i.Init()
	return i
}

