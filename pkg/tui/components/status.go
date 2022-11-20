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

type Status struct {
	notification *views.SimpleStyledTextBar
	helpLine     views.Widget

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
	helpLine1 := makeLine("^X Quit", "▴▾ Navigate", "^P Browse", "^O Save")
	helpLine2 := makeLine("^C Copy", "^H Reveal", "ESC Close", "^G Help")

	status.AddWidget(notification, 1)
	status.AddWidget(helpLine1, 1)
	status.AddWidget(helpLine2, 1)

	status.notification = notification
	status.helpLine = helpLine1

	return status
}

func makeLine(blocks ...string) views.Widget {
	line := views.NewBoxLayout(views.Horizontal)

	for _, block := range blocks {
		blockElement := views.NewText()
    blockElement.SetText(block)

    spaceIndex := runewidth.StringWidth(strings.Split(block, " ")[0])

    for i := 0; i < spaceIndex; i++ {
      blockElement.SetStyleAt(i, tcell.StyleDefault.Reverse(true))
    }

		line.AddWidget(blockElement, 1.0/float64(len(blocks)))
	}

	return line
}
