package tui

import (
	"github.com/gdamore/tcell/v2"
	"github.com/gdamore/tcell/v2/views"
	"github.com/lithammer/fuzzysearch/fuzzy"
	"github.com/shikaan/kpcli/pkg/tui/components"
)

type ListView struct {
  views.Panel
}

func NewListView(screen *tcell.Screen, props *State) views.Widget {
  c := components.NewContainer(*screen)
  v := &ListView{}
  l := newList(props.Database.GetEntryPaths())
  o := &components.InputOptions{}
  i := components.NewInput(o)

  i.OnChange(func(ev tcell.Event) bool {
    var entries []string

    if len(i.GetContent()) < 3 {
      entries = props.Database.GetEntryPaths()
    } else {
      entries = fuzzy.FindFold(i.GetContent(), props.Database.GetEntryPaths())
    }

    v.RemoveWidget(l)
    l = newList(entries)
    v.AddWidget(l,0)

    return true
  })

  i.SetFocus(true)
  v.AddWidget(i, 0)
  v.AddWidget(l, 0)

  c.SetContent(v)
  return c
}

func newList(entries []string) views.Widget {
  box := views.NewBoxLayout(views.Vertical)
  
  for i, e := range entries {
    line := views.NewTextArea()
    line.SetContent(e)
    box.AddWidget(line, 0)
  
    if i == 0 {
      line.EnableCursor(true)
      line.HideCursor(false)
    }
  }

  return box
}
