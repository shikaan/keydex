package components

import (
	"fmt"
	"sync"

	"github.com/gdamore/tcell/v2"
	"github.com/gdamore/tcell/v2/views"
	"github.com/lithammer/fuzzysearch/fuzzy"
)

type AutoComplete struct {
  options AutoCompleteOptions
  matchedEntries []string
  
  list *WithFocusables
  counter *views.SimpleStyledText
  input *Input

  views.BoxLayout
}

type AutoCompleteOptions struct {
  Screen tcell.Screen
  Entries []string
  TotalCount int
  MaxX int
  MaxY int

  OnSelect func (entry string) bool
}

func NewAutoComplete(options AutoCompleteOptions) views.Widget {
	autoComplete := &AutoComplete{}
  autoComplete.SetOrientation(views.Vertical)
  autoComplete.options = options

	input := NewInput(&InputOptions{})
	autoComplete.input = input
  autoComplete.AddWidget(input, 0)

  input.SetFocus(true)
  input.OnChange(func(ev tcell.Event) bool {
    entries := options.Entries
		content := input.GetContent()

		if content != "" {
			entries = fuzzy.FindFold(content, options.Entries)
		}

    autoComplete.drawList(entries)
    autoComplete.drawCounter(entries)

		return false
	})
	
  counter := views.NewSimpleStyledText()
	counter.SetAlignment(views.HAlignRight) 
  counter.SetStyle(tcell.StyleDefault.Reverse(true))
  autoComplete.counter = counter
	autoComplete.AddWidget(counter, 0)

  autoComplete.drawList(options.Entries)
  autoComplete.drawCounter(options.Entries)

	return autoComplete
}

func (ac *AutoComplete) drawList(entries []string) {
  d := &WithFocusables{}
  d.SetOrientation(views.Vertical)

  for i, e := range entries {
		if i >= ac.options.MaxY {
			break
		}

		line := newOption()
		if len(e) > ac.options.MaxX {
			line.SetContent(e[:ac.options.MaxX])
		} else {
			line.SetContent(e)
		}

    entry := e
    line.OnSelect(func() bool {
      return ac.options.OnSelect(entry) 
    })

		d.AddWidget(line, 0)
    if i == 0 {
      line.SetFocus(true)
    }
	}

  ac.RemoveWidget(ac.list)
  ac.list = d
  ac.AddWidget(ac.list, 0)
}

func (ac *AutoComplete) drawCounter(entries []string) {
  matched := len(entries)
  counter := fmt.Sprintf("%d/%d", matched, len(ac.options.Entries))
  
  ac.counter.SetStyle(tcell.StyleDefault.Bold(matched == 0))
	
  ac.counter.SetText(fmt.Sprintf("% s", counter))
}

// Option

type option struct {
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
  
  selectHandler func() bool
  focusHandler func() bool
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
	m.width = len(text)
	m.content = text

	i.CellView.SetModel(m)
}

func (i *option) GetContent() string {
	return string(i.model.content)
}

func (i *option) SetInputType(t InputType) {
  i.model.inputType = t
  i.Init()
}

func (i *option) GetInputType() InputType {
  return i.model.inputType
}

func (i *option) HandleEvent(ev tcell.Event) bool {
	if !i.HasFocus() {
		return false
	}

	switch ev := ev.(type) {
	case *tcell.EventKey:  
    if ev.Key() == tcell.KeyEnter {
      if i.model.selectHandler != nil {
        i.model.selectHandler()
      }
    }	
  }
	
  return i.CellView.HandleEvent(ev)
}

func (i *option) OnSelect (cb func () bool) func () {
  i.model.selectHandler = cb
  return func () {
    i.model.selectHandler = nil
  }
}

func (i *option) OnFocus (cb func () bool) func () {
  i.model.focusHandler = cb
  return func () {
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
