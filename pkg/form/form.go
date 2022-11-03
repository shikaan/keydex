package form

import (
	"fmt"
	"os"

	"github.com/gdamore/tcell/v2"
	"github.com/gdamore/tcell/v2/views"
)

var app = &views.Application{}

type root struct {
  views.Panel
}

func (r *root) HandleEvent(ev tcell.Event) bool {
	switch ev := ev.(type) {
	case *tcell.EventKey:
		if ev.Key() == tcell.KeyEscape || ev.Name() == "Ctrl+C"{
			app.Quit()
			return true
		}
	}
	return r.Panel.HandleEvent(ev)
}

func Run() {
  r := &root{}
  r.SetStyle(tcell.StyleDefault)

  title := NewTitle("This is the Title")
  r.SetTitle(title)

  main := views.NewBoxLayout(views.Vertical)

  field := NewField("label", "initial")
  main.AddWidget(field, 0)
  
  r.SetContent(main)

  status := NewStatus()
  r.SetStatus(status)

  app.SetRootWidget(r)
  if e := app.Run(); e != nil {
		fmt.Fprintln(os.Stderr, e.Error())
		os.Exit(1)
	}
}
