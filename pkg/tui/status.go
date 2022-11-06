package tui

import (
	"github.com/gdamore/tcell/v2"
	"github.com/gdamore/tcell/v2/views"
)

func NewStatus() *views.TextBar {
  tb := views.NewTextBar()
  tb.SetStyle(tcell.StyleDefault.Reverse(true))
  
  tb.SetLeft("[^C] Quit", tcell.StyleDefault.Reverse(true))
  tb.SetCenter("[Up] Move Up", tcell.StyleDefault.Reverse(true))
  tb.SetRight("[Down] Move Down", tcell.StyleDefault.Reverse(true))

  return tb
}
