package tui

import (
	"strings"
	"testing"

	"github.com/gdamore/tcell/v2"
	"github.com/gdamore/tcell/v2/views"
	"github.com/shikaan/keydex/tui/components"
)

func firstLine(s tcell.Screen) string {
	w, _ := s.Size()
	var b strings.Builder
	for x := range w {
		c, _, _, _ := s.GetContent(x, 0)
		if c == 0 {
			b.WriteRune(' ')
		} else {
			b.WriteRune(c)
		}
	}
	return b.String()
}

func TestSetDirtyAddsModifiedBadge(t *testing.T) {
	screen := tcell.NewSimulationScreen("UTF-8")
	if err := screen.Init(); err != nil {
		t.Fatalf("failed to init simulation screen: %v", err)
	}
	defer screen.Fini()

	title := components.NewTitle("TestDB")
	title.SetTitle("MyTitle")

	layout := &Layout{
		Title:  title,
		Screen: screen,
	}

	App.layout = layout
	App.screen = screen

	vp := views.NewViewPort(screen, 0, 0, 80, 1)
	title.SetView(vp)
	title.Draw()
	initial := firstLine(screen)
	if strings.Contains(initial, "[MODIFIED]") {
		t.Fatalf("expected no [MODIFIED] badge initially, got: %q", initial)
	}

	App.SetDirty(true)

	screen.Clear()
	title.Draw()
	after := firstLine(screen)

	if !strings.Contains(after, "[MODIFIED]") {
		t.Fatalf("expected [MODIFIED] badge in title after SetDirty(true), got: %q", after)
	}

	App.SetDirty(false)
	screen.Clear()
	title.Draw()
	final := firstLine(screen)

	if strings.Contains(final, "[MODIFIED]") {
		t.Fatalf("expected [MODIFIED] badge to be removed after SetDirty(false), got: %q", final)
	}
}
