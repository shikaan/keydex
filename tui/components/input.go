package components

import (
	"sync"

	"github.com/gdamore/tcell/v2"
	"github.com/gdamore/tcell/v2/views"
	"github.com/shikaan/keydex/pkg/utils"
)

type Input struct {
	model *inputModel
	once  sync.Once

	Focusable
	views.CellView
}

type InputOptions struct {
	InitialValue string
	Type         InputType
}

type InputType int

const (
	InputTypeText InputType = iota
	InputTypePassword
)

// Model - Used internally by tcell/views to handle drawing

type inputModel struct {
	content   string
	width     int
	x         int
	style     tcell.Style
	hasFocus  bool
	inputType InputType

	keyPressHandler func(ev *tcell.EventKey) bool
	changeHandler   func(ev tcell.Event) bool
	focusHandler    func() bool
}

func (m *inputModel) GetCell(x, y int) (rune, tcell.Style, []rune, int) {
	isPassword := m.inputType == InputTypePassword
	isOutOfBound := y != 0 || x < 0 || x >= len(m.content)

	if isOutOfBound {
		return 0, m.style, nil, 1
	}

	char := rune(m.content[x])
	if isPassword {
		char = '*'
	}

	return char, m.style, nil, 1
}

func (m *inputModel) GetBounds() (int, int) {
	return m.width, 1
}

func (m *inputModel) limitCursor() {
	x, _ := m.GetBounds()

	if m.x > x {
		m.x = x
	}
	if m.x < 0 {
		m.x = 0
	}
}

func (m *inputModel) SetCursor(x, y int) {
	m.x = x
	m.limitCursor()
}

func (m *inputModel) MoveCursor(x, y int) {
	m.x += x
	m.limitCursor()
}

func (m *inputModel) GetCursor() (int, int, bool, bool) {
	return m.x, 0, true, m.hasFocus
}

// Input - Models an input (similar to HTML inputs)

func (i *Input) HasFocus() bool {
	return i.model.hasFocus
}

func (i *Input) SetFocus(on bool) {
	i.Init()
	i.model.hasFocus = on
	i.CellView.SetModel(i.model)
	if i.model.focusHandler != nil {
		i.model.focusHandler()
	}
}

func (i *Input) SetContent(text string) {
	i.Init()
	m := i.model
	m.width = len(text)
	m.content = text

	i.CellView.SetModel(m)
}

func (i *Input) GetContent() string {
	return string(i.model.content)
}

func (i *Input) SetInputType(t InputType) {
	i.model.inputType = t
	i.Init()
}

func (i *Input) GetInputType() InputType {
	return i.model.inputType
}

func (i *Input) HandleEvent(ev tcell.Event) bool {
	if !i.HasFocus() {
		return false
	}

	switch ev := ev.(type) {
	case *tcell.EventKey:
		handled := false

		if i.model.keyPressHandler != nil {
			handled = i.model.keyPressHandler(ev)
		}

		if handled {
			return handled
		}

		if ev.Key() == tcell.KeyRune {
			return i.handleContentUpdate(
				ev,
				func(c string, x int) (string, int) {
					return c[:x] + string(ev.Rune()) + c[x:], 1
				},
			)
		}

		if ev.Key() == tcell.KeyBackspace2 || ev.Key() == tcell.KeyBackspace {
			return i.handleContentUpdate(
				ev,
				func(c string, x int) (string, int) {
					safeX := utils.Max(0, x-1)

					return c[:safeX] + c[x:], -1
				},
			)
		}

		if ev.Key() == tcell.KeyDelete {
			return i.handleContentUpdate(
				ev,
				func(c string, x int) (string, int) {
					safeX := utils.Min(len(c), x+1)

					return c[:x] + c[safeX:], 0
				},
			)
		}

		// CellView (few lines below) would handle these events, preventing other
		// components (e.g., autocomplete) to handle them.
		// Not really nice, but not worth the complication of doing it nicer either
		if ev.Key() == tcell.KeyUp || ev.Key() == tcell.KeyDown {
			return false
		}
	}
	return i.CellView.HandleEvent(ev)
}

func (i *Input) handleContentUpdate(ev tcell.Event, cb func(content string, cursorPosition int) (string, int)) bool {
	m := i.GetModel()
	x, _, _, _ := m.GetCursor()

	content := i.GetContent()
	newContent, cursorOffset := cb(content, x)
	i.SetContent(newContent)
	m.MoveCursor(cursorOffset, 0)

	if i.model.changeHandler != nil {
		i.model.changeHandler(ev)
	}

	return true
}

func (i *Input) OnKeyPress(cb func(ev *tcell.EventKey) bool) func() {
	i.model.keyPressHandler = cb
	return func() {
		i.model.keyPressHandler = nil
	}
}

func (i *Input) OnChange(cb func(ev tcell.Event) bool) func() {
	i.model.changeHandler = cb
	return func() {
		i.model.changeHandler = nil
	}
}

func (i *Input) OnFocus(cb func() bool) func() {
	i.model.focusHandler = cb
	return func() {
		i.model.focusHandler = nil
	}
}

func (i *Input) Init() {
	i.once.Do(func() {
		m := &inputModel{}
		i.model = m
		i.CellView.Init()
		i.CellView.SetModel(m)
	})
}

func NewInput(options *InputOptions) *Input {
	i := &Input{}
	i.Init()
	i.model.inputType = options.Type

	return i
}
