package components

import (
	"sync"

	"github.com/gdamore/tcell/v2"
	"github.com/gdamore/tcell/v2/views"
)

type Option struct {
	model *optionModel
	once  sync.Once

	Focusable
	views.CellView
}

type optionModel struct {
	content   string
	width     int
	x         int
	style     tcell.Style
	hasFocus  bool
	inputType InputType
  
  selectHandler func(ev *tcell.EventKey) bool
}

func (m *optionModel) GetCell(x, y int) (rune, tcell.Style, []rune, int) {
  char := ' '
  if len(m.content) > x {
    char = rune(m.content[x])
  }

  return char, m.style, nil, 1
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

func (i *Option) HasFocus() bool {
	return i.model.hasFocus
}

func (i *Option) SetFocus(on bool) {
	i.Init()
  if on {
    i.SetContent("> " + i.model.content)
  } else {
    i.SetContent(i.model.content[2:])
  }

  i.model.hasFocus = on

	i.CellView.SetModel(i.model)
}

func (i *Option) SetContent(text string) {
	i.Init()
	m := i.model
	m.width = len(text)
	m.content = text

	i.CellView.SetModel(m)
}

func (i *Option) GetContent() string {
	return string(i.model.content)
}

func (i *Option) SetInputType(t InputType) {
  i.model.inputType = t
  i.Init()
}

func (i *Option) GetInputType() InputType {
  return i.model.inputType
}

func (i *Option) HandleEvent(ev tcell.Event) bool {
	if !i.HasFocus() {
		return false
	}

	switch ev := ev.(type) {
	case *tcell.EventKey:  
    if ev.Key() == tcell.KeyEnter {
      if i.model.selectHandler != nil {
        i.model.selectHandler(ev)
      }
    }	
  }
	
  return i.CellView.HandleEvent(ev)
}

func (i *Option) OnSelect (cb func (ev *tcell.EventKey) bool) func () {
  i.model.selectHandler = cb
  return func () {
    i.model.selectHandler = nil
  }
}

func (i *Option) Init() {
	i.once.Do(func() {
	  m := &optionModel{}
		i.model = m
		i.CellView.Init()
		i.CellView.SetModel(m)
	})
}

func NewOption() *Option {
	i := &Option{}
	i.Init()

	return i
}
