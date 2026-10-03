package field

import (
	"sync"

	"github.com/gdamore/tcell/v2"
	"github.com/gdamore/tcell/v2/views"
	"github.com/shikaan/keydex/tui/components"
	"github.com/shikaan/keydex/tui/components/line"
)

type Checkbox struct {
	model *checkboxModel
	once  sync.Once

	components.Focusable
	views.CellView
}

type CheckboxOptions struct {
	InitialValue bool
	Disabled     bool
}

type checkboxModel struct {
	content  bool
	hasFocus bool
	disabled bool
	cells    [3]rune

	keyPressHandler func(ev *tcell.EventKey) bool
	changeHandler   func(ev tcell.Event) bool
	focusHandler    func() bool
}

func (m *checkboxModel) GetCell(x, y int) (rune, tcell.Style, []rune, int) {
	if x > 2 || y > 0 {
		return line.EMPTY_CELL, tcell.StyleDefault, nil, 1
	}

	return m.cells[x], tcell.StyleDefault, nil, 1
}

func (m *checkboxModel) GetBounds() (int, int) {
	return 3, 1
}

func (m *checkboxModel) SetCursor(x, y int) {
}

func (m *checkboxModel) MoveCursor(x, y int) {
}

func (m *checkboxModel) GetCursor() (int, int, bool, bool) {
	return 1, 0, true, m.hasFocus
}

func (c *Checkbox) HasFocus() bool {
	return c.model.hasFocus
}

func (c *Checkbox) SetFocus(on bool) {
	c.Init()
	c.model.hasFocus = on
	c.CellView.SetModel(c.model)
	if c.model.focusHandler != nil {
		c.model.focusHandler()
	}
}

func (c *Checkbox) SetContent(value bool) {
	c.Init()
	c.model.content = value

	m := c.model
	m.cells[0] = '['
	if m.content {
		m.cells[1] = 'X'
	} else {
		m.cells[1] = ' '
	}
	m.cells[2] = ']'

	c.CellView.SetModel(m)
}

func (c *Checkbox) GetContent() bool {
	return c.model.content
}

func (c *Checkbox) HandleEvent(ev tcell.Event) bool {
	if !c.HasFocus() {
		return false
	}

	switch ev := ev.(type) {
	case *tcell.EventKey:
		handled := false

		if c.model.keyPressHandler != nil {
			handled = c.model.keyPressHandler(ev)
		}

		if handled {
			return handled
		}

		// Prevent all input-changing actions when the field is disabled
		if c.model.disabled {
			return false
		}

		switch ev.Key() {
		case tcell.KeyEnter:
			c.SetContent(!c.model.content)
			if c.model.changeHandler != nil {
				return c.model.changeHandler(ev)
			}

			return true
		}
	}

	return false
}

func (c *Checkbox) OnKeyPress(cb func(ev *tcell.EventKey) bool) func() {
	c.model.keyPressHandler = cb
	return func() {
		c.model.keyPressHandler = nil
	}
}

func (c *Checkbox) OnChange(cb func(ev tcell.Event) bool) func() {
	c.model.changeHandler = cb
	return func() {
		c.model.changeHandler = nil
	}
}

func (c *Checkbox) OnFocus(cb func() bool) func() {
	c.model.focusHandler = cb
	return func() {
		c.model.focusHandler = nil
	}
}

func (c *Checkbox) Init() {
	c.once.Do(func() {
		c.model = newCheckboxModel()
		c.CellView.Init()
		c.CellView.SetModel(c.model)
	})
}

func newCheckboxModel() *checkboxModel {
	m := &checkboxModel{}
	return m
}

func NewCheckbox(options *CheckboxOptions) *Checkbox {
	c := &Checkbox{}
	c.Init()
	c.model.disabled = options.Disabled
	c.SetContent(options.InitialValue)
	return c
}
