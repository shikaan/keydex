package tui

import (
	"github.com/gdamore/tcell/v2"
	"github.com/gdamore/tcell/v2/views"
)

type View struct {
	views.Panel
}

type ViewCreator[K interface{}] func(screen tcell.Screen, options K) views.Widget

func (v *View) HandleEvent(ev tcell.Event) bool {
	switch ev := ev.(type) {
	case *tcell.EventKey:
		if ev.Name() == "Ctrl+X" {
			App.Quit()
			return true
		}
	}
	return v.Panel.HandleEvent(ev)
}

func NewView() *View {
	return &View{}
}
