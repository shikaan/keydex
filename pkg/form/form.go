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

type form struct {
  views.BoxLayout
}

func (r *root) HandleEvent(ev tcell.Event) bool {
	switch ev := ev.(type) {
	case *tcell.EventKey:
		if ev.Key() == tcell.KeyEscape || ev.Name() == "Ctrl+C"{
			app.Quit()
			return true
		}
    if ev.Key() == tcell.KeyTab {
      // TODO: This is ugly
      r.Widgets()[1].(*form).focusNext()
      return true
    }
  }
	return r.Panel.HandleEvent(ev)
}

func (m *form) focusNext() {
  shouldFocus := false
  firstIndex := -1

  for i, widget := range m.Widgets() {
    switch widget := widget.(type) {
    case *Field:
      if firstIndex < 0 {
        firstIndex = i
      }
      if shouldFocus {
        widget.SetFocus(true)
        return
      } else if widget.HasFocus() {
        widget.SetFocus(false)
        shouldFocus = true
      } 
    }
  }

  if shouldFocus {
    m.Widgets()[firstIndex].(*Field).SetFocus(true)
  }
}

func Run() {
  r := &root{}
  r.SetStyle(tcell.StyleDefault)

  title := NewTitle("This is the Title")
  r.SetTitle(title)

  main := &form{}
  main.SetOrientation(views.Vertical)

  field := NewField("label", "initial")
  main.AddWidget(field, 0)
  
  field2 := NewField("label2", "initial2")
  main.AddWidget(field2, 0)

  field.SetFocus(true)

  r.SetContent(main)

  status := NewStatus()
  r.SetStatus(status)

  app.SetRootWidget(r)
  if e := app.Run(); e != nil {
		fmt.Fprintln(os.Stderr, e.Error())
		os.Exit(1)
	}
}
