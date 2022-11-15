package components

import (
	"github.com/gdamore/tcell/v2"
	"github.com/gdamore/tcell/v2/views"
)

type Title = views.TextBar

func NewTitle(title string) *views.TextBar {
	tb := views.NewTextBar()
	tb.SetStyle(tcell.StyleDefault.Reverse(true))
	tb.SetCenter(title, tcell.StyleDefault.Reverse(true))

	return tb
}
