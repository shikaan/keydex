package status

import (
	"github.com/gdamore/tcell/v2"
	"github.com/gdamore/tcell/v2/views"
	"github.com/shikaan/keydex/tui/components"
)

type Prompt struct {
	text  *views.Text
	model promptModel

	components.Focusable
	views.BoxLayout
}

type promptModel struct {
	hasFocus bool

	// Handle keypress events: triggered every time a key is pressed
	// Returns true if handled, false if needs cascading
	keyPressHandler func(ev *tcell.EventKey) bool
}

func (p *Prompt) SetText(s string) {
	p.text.SetText(s)
}

func (p *Prompt) SetFocus(on bool) {
	p.model.hasFocus = on
}

func (p *Prompt) OnKeyPress(cb func(ev *tcell.EventKey) bool) func() {
	p.model.keyPressHandler = cb
	return func() {
		p.model.keyPressHandler = nil
	}
}

func (p *Prompt) HandleEvent(ev tcell.Event) bool {
	if !p.model.hasFocus {
		return false
	}

	switch ev := ev.(type) {
	case *tcell.EventKey:
		return p.model.keyPressHandler != nil && p.model.keyPressHandler(ev)
	}

	return false
}

func newPrompt() *Prompt {
	p := &Prompt{}
	p.SetOrientation(views.Horizontal)
	p.SetStyle(tcell.StyleDefault.Reverse(true))

	p.text = views.NewText()
	p.text.SetStyle(tcell.StyleDefault.Reverse(true))

	p.InsertWidget(0, p.text, 0)

	return p
}
