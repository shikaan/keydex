package components

import (
	"fmt"

	"github.com/gdamore/tcell/v2"
	"github.com/gdamore/tcell/v2/views"
	"github.com/shikaan/keydex/pkg/info"
)

type Title struct {
	content string

	views.TextBar
}

func (t *Title) SetTitle(title string) {
	t.content = title
	t.SetCenter(title, tcell.StyleDefault.Reverse(true))
}

func (t *Title) SetDirty(dirty bool) {
	if dirty {
		t.SetCenter(t.content+" [MODIFIED]", tcell.StyleDefault.Reverse(true))
	} else {
		t.SetCenter(t.content, tcell.StyleDefault.Reverse(true))
	}
}

func NewTitle(database string) *Title {
	tb := &Title{}
	tb.TextBar.SetStyle(tcell.StyleDefault.Reverse(true))
	left := fmt.Sprintf("  %s %s (%s)", info.NAME, info.VERSION, info.REVISION)
	tb.TextBar.SetLeft(left, tcell.StyleDefault.Reverse(true))
	right := fmt.Sprintf("%s  ", database)
	tb.TextBar.SetRight(right, tcell.StyleDefault.Reverse(true))

	return tb
}
