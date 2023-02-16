package components

import (
	"fmt"
	"strings"
	"time"

	"github.com/gdamore/tcell/v2"
	"github.com/gdamore/tcell/v2/views"
	"github.com/mattn/go-runewidth"
)

// Notifications will be cleared upon the first interaction
// after NOTIFICATION_MIN_DURATION_IN_SECONDS seconds
const NOTIFICATION_MIN_DURATION_IN_SECONDS = 5
const EMPTY_NOTIFICATION = " "

type Status struct {
	model statusModel

	notification *views.TextBar
	helpLines    [2]views.Widget

	prompt       *prompt
	confirmLines [2]views.Widget

	views.BoxLayout
}

type statusModel struct {
	isConfirming bool
	onAccept     func()
	onReject     func()
}

func (s *Status) Notify(st string) {
	s.notification.SetCenter(fmt.Sprintf("[ %s ]", st), tcell.StyleDefault)
	go func() {
		time.Sleep(NOTIFICATION_MIN_DURATION_IN_SECONDS * time.Second)
		s.notification.SetCenter(EMPTY_NOTIFICATION, tcell.StyleDefault)
	}()
}

func (s *Status) Confirm(message string, onAccept func(), onReject func()) {
	s.model.isConfirming = true

	s.prompt.SetText(message)
	s.prompt.SetFocus(true)
	s.model.onAccept = onAccept
	s.model.onReject = onReject

	s.RemoveWidget(s.notification)
	for _, l := range s.helpLines {
		s.RemoveWidget(l)
	}

	s.InsertWidget(0, s.prompt, 0)
	for i, l := range s.confirmLines {
		s.InsertWidget(i+1, l, 0)
	}

	s.Resize()
}

func (s *Status) reset() {
	s.model.isConfirming = false
	s.prompt.SetFocus(false)

	s.RemoveWidget(s.prompt)
	for _, l := range s.confirmLines {
		s.RemoveWidget(l)
	}

	s.InsertWidget(0, s.notification, 0)
	for i, l := range s.helpLines {
		s.InsertWidget(i+1, l, 0)
	}
	s.Resize()
}

func NewStatus() *Status {
	status := &Status{}
	status.SetOrientation(views.Vertical)

	status.notification = views.NewTextBar()
	// Prevents jumps on the first render
	status.notification.SetCenter(EMPTY_NOTIFICATION, tcell.StyleDefault)

	status.helpLines[0] = newLine("^X Exit", "▴▾ Navigate", "^P  Browse       ", "^O Save")
	status.helpLines[1] = newLine("^C Copy", "^R Reveal  ", "ESC To Last Entry", "^G Help")

	status.prompt = newPrompt()
	status.prompt.OnKeyPress(func(ev *tcell.EventKey) bool {
		if !status.model.isConfirming {
			return status.BoxLayout.HandleEvent(ev)
		}

		if ev.Rune() == 'y' || ev.Rune() == 'Y' {
			if status.model.onAccept != nil {
				status.model.onAccept()
			}
			status.reset()
			return true
		}

		if ev.Rune() == 'n' || ev.Rune() == 'N' || ev.Name() == "Ctrl+C" || ev.Key() == tcell.KeyESC {
			if status.model.onReject != nil {
				status.model.onReject()
			}
			status.reset()
			return true
		}

		return true
	})
	status.confirmLines[0] = newLine("Y Yes")
	status.confirmLines[1] = newLine("N No")

	status.reset()

	return status
}

func newLine(blocks ...string) views.Widget {
	l := views.NewBoxLayout(views.Horizontal)

	for _, block := range blocks {
		blockElement := views.NewText()
		blockElement.SetText(block)

		spaceIndex := runewidth.StringWidth(strings.Split(block, " ")[0])

		for i := 0; i < spaceIndex; i++ {
			blockElement.SetStyleAt(i, tcell.StyleDefault.Reverse(true))
		}

		l.AddWidget(blockElement, 1.0/float64(len(blocks)))
	}

	return l
}

type prompt struct {
	text  *views.Text
	input *Input

	views.BoxLayout
}

func (p *prompt) SetText(s string) {
	p.text.SetText(s)
}

func (p *prompt) SetFocus(on bool) {
	p.input.SetFocus(on)
}

func (p *prompt) OnKeyPress(cb func(ev *tcell.EventKey) bool) func() {
	return p.input.OnKeyPress(cb)
}

func newPrompt() *prompt {
	p := &prompt{}
	p.SetOrientation(views.Horizontal)
	p.SetStyle(tcell.StyleDefault.Reverse(true))

	p.text = views.NewText()
	p.text.SetStyle(tcell.StyleDefault.Reverse(true))

	p.input = NewInput(&InputOptions{})
	p.input.SetStyle(tcell.StyleDefault.Reverse(true))

	p.InsertWidget(0, p.text, 0)
	p.InsertWidget(1, p.input, 1)

	return p
}
