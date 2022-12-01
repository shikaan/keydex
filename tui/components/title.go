package components

import (
	"fmt"

	"github.com/gdamore/tcell/v2"
	"github.com/gdamore/tcell/v2/views"
	"github.com/shikaan/kpcli/pkg/info"
)

type Title struct {
  views.TextBar
}

func  (t *Title) SetTitle(s string) {
  t.SetCenter(s, tcell.StyleDefault.Reverse(true))
}

func NewTitle(database string) *Title {
	tb := &Title{}
	tb.TextBar.SetStyle(tcell.StyleDefault.Reverse(true))
  left := fmt.Sprintf("  %s %s (%s)", info.NAME, info.VERSION, info.REVISION)
  tb.TextBar.SetLeft(left, tcell.StyleDefault.Reverse(true))
	tb.TextBar.SetRight(database, tcell.StyleDefault.Reverse(true))

	return tb
}
