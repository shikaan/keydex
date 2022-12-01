package tui

import (
	"github.com/gdamore/tcell/v2"
	"github.com/gdamore/tcell/v2/views"
	"github.com/shikaan/kpcli/tui/components"
)

type Layout struct {
	Status *components.Status
	Title  *components.Title

	Screen tcell.Screen

	views.Panel
}

func (v *Layout) HandleEvent(ev tcell.Event) bool {
	switch ev := ev.(type) {
	case *tcell.EventKey:
		if ev.Name() == "Ctrl+X" {
			App.Quit()
			return true
		}
		if ev.Name() == "Ctrl+P" {
			App.NavigateTo(NewListView)
			return true
		}	
		if ev.Name() == "Ctrl+G" {
			App.NavigateTo(NewHelpView)
			return true
		}
		if ev.Name() == "Ctrl+C" {
      handled := v.Panel.HandleEvent(ev)

      if !handled {
        App.Notify("No field selected for copy. Use ^X to close.")
      }
      
      return true
		}
    if ev.Key() == tcell.KeyEsc {
			if App.State.Entry == nil {
				App.Notify("No entry selected yet.")
				return true
			}

			App.NavigateTo(NewHomeView)
			return true
		}
	}
	return v.Panel.HandleEvent(ev)
}

// Returns a Layout component responsible for the shell of the application
// and of the routing in between pages
func NewLayout(screen tcell.Screen) *Layout {
	l := &Layout{}
	title := components.NewTitle("")
	status := components.NewStatus()

	t := views.NewText()
	t.SetText(" ")
	l.SetMenu(t)

	l.SetStatus(status)
	l.Status = status

	l.SetTitle(title)
	l.Title = title

	l.Screen = screen

	return l
}
