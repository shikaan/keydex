package components

import (
	"fmt"
	"math"
	"sync"

	"github.com/gdamore/tcell/v2"
	"github.com/gdamore/tcell/v2/views"
	"github.com/lithammer/fuzzysearch/fuzzy"
	"github.com/mattn/go-runewidth"
)

type AutoComplete struct {
	options AutoCompleteOptions

	list    *WithFocusables
	counter *views.SimpleStyledText
	input   *Input

	Focusable
	views.BoxLayout
}

type AutoCompleteOptions struct {
	Screen     tcell.Screen
	Entries    []string
	TotalCount int
	MaxX       int
	MaxY       int

	OnSelect func(entry string) bool
	OnFocus  func() bool
}

func NewAutoComplete(options AutoCompleteOptions) *AutoComplete {
	autoComplete := &AutoComplete{}
	autoComplete.SetOrientation(views.Vertical)
	autoComplete.options = options

	input := NewInput(&InputOptions{})
	autoComplete.input = input
	autoComplete.AddWidget(input, 0)

	input.OnChange(func(ev tcell.Event) bool {
		entries := options.Entries
		content := input.GetContent()

		if len(content) > 0 {
			entries = fuzzy.FindFold(content, options.Entries)
		}

		autoComplete.drawList(entries)
		autoComplete.drawCounter(entries)

		return false
	})

	counter := views.NewSimpleStyledText()
	counter.SetAlignment(views.HAlignRight)
	autoComplete.counter = counter
	autoComplete.AddWidget(counter, 0)

	autoComplete.drawList(options.Entries)
	autoComplete.drawCounter(options.Entries)

	return autoComplete
}

func (ac *AutoComplete) OnFocus(cb func() bool) func() {
	return ac.input.OnFocus(cb)
}

func (ac *AutoComplete) SetFocus(on bool) {
	ac.input.SetFocus(on)
}

func (ac *AutoComplete) HasFocus() bool {
	return ac.input.HasFocus()
}

func (ac *AutoComplete) drawList(entries []string) {
	container := &WithFocusables{}
	container.SetOrientation(views.Vertical)

	// The 2 characters are used for when the entry is selected
	// and you have "> " prepended
	maxLineLength := ac.options.MaxX - 2

	if len(entries) == 0 {
		line := views.NewSimpleStyledText()
		line.SetText(runewidth.FillRight("--- No Results ---", ac.options.MaxX))
		container.AddWidget(line, 0)
	}

	for i, entry := range entries {
		if i >= ac.options.MaxY {
			break
		}

		line := newOption()
		line.SetContent((runewidth.FillRight(runewidth.Truncate(entry, maxLineLength, ""), maxLineLength)))

		// For memoization
		e := entry
		line.OnSelect(func() bool {
			return ac.options.OnSelect(e)
		})

		container.AddWidget(line, 0)
		if i == 0 {
			line.SetFocus(true)
		}
	}

	ac.RemoveWidget(ac.list)
	ac.list = container
	ac.AddWidget(ac.list, 0)
}

func (ac *AutoComplete) drawCounter(entries []string) {
	matched := int(math.Min(float64(len(entries)), float64(ac.options.MaxY)))
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
	// Content as it appears outside
	content string
	// 0-spaced rune sequence, to facilitate printing of Unicode chars
	runes []rune

	width     int
	x         int
	style     tcell.Style
	hasFocus  bool
	inputType InputType

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
