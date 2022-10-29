package form

import (
	"github.com/gdamore/tcell/v2"
	"github.com/gdamore/tcell/v2/views"
)

type Input struct {
  label string
  
  views.TextArea
}

func (i *Input) HandleEvent(ev tcell.Event) bool {
  switch ev := ev.(type) {	
  case *tcell.EventKey:
		if ev.Key() == tcell.KeyCtrlE {
			println("asdfasdfasdfasdfasdfad")
			return true
		}
	}
	return i.TextArea.HandleEvent(ev)
}


func NewInput() *Input {
  text := &Input{ label: "test" }
  text.Init()
  text.SetContent("lolololo")
  text.EnableCursor(true)

  return text
}
