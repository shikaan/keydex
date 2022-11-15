package tui

import (
	"github.com/gdamore/tcell/v2"
	"github.com/gdamore/tcell/v2/views"
	"github.com/shikaan/kpcli/pkg/tui/components"
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
		if ev.Key() == tcell.KeyEsc {
		  App.NavigateTo(NewEditView)	
			return true
		}
	}
	return v.Panel.HandleEvent(ev)
}

func NewLayout(screen tcell.Screen) *Layout {
	l := &Layout{}
	title := components.NewTitle("kpcli")
	status := components.NewStatus()

	l.SetStatus(status)
	l.Status = status

	l.SetTitle(title)
	l.Title = title

  l.Screen = screen

	return l
}
