package components

import (
	"fmt"
	"time"

	"github.com/gdamore/tcell/v2"
	"github.com/gdamore/tcell/v2/views"
)

// Notifications will be cleared upon the first interaction
// after NOTIFICATION_MIN_DURATION_IN_SECONDS seconds
const NOTIFICATION_MIN_DURATION_IN_SECONDS = 5

type Status struct {
	notification *views.SimpleStyledTextBar
	helpLine     *views.TextBar

	views.BoxLayout
}

func (s *Status) Notify(st string) {
	s.notification.SetCenter(fmt.Sprintf("[ %s ]", st))
	go func() {
		time.Sleep(NOTIFICATION_MIN_DURATION_IN_SECONDS * time.Second)
		s.notification.SetCenter("")
	}()
}

func NewStatus() *Status {
	status := &Status{}
	status.SetOrientation(views.Vertical)

	notification := views.NewSimpleStyledTextBar()
	helpLine1 := makeLine("[^X] Quit", "[▴▾] Navigate", "[^O] Save")
	helpLine2 := makeLine("[^C] Copy field to clipboard", "[^P] Browse entries", "[^H] Help")

	status.AddWidget(notification, 1)
	status.AddWidget(helpLine1, 1)
	status.AddWidget(helpLine2, 1)

	status.notification = notification
	status.helpLine = helpLine1

	return status
}

func makeLine(left, center, right string) *views.TextBar {
	line := views.NewTextBar()
	line.SetStyle(tcell.StyleDefault.Reverse(true))

	line.SetLeft(left, tcell.StyleDefault.Reverse(true))
	line.SetCenter(center, tcell.StyleDefault.Reverse(true))
	line.SetRight(right, tcell.StyleDefault.Reverse(true))

	return line
}
