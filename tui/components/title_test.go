package components

import (
	"strings"
	"testing"

	"github.com/gdamore/tcell/v2"
	"github.com/gdamore/tcell/v2/views"
)

// TODO: worth extracting in a test package?
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

func TestSetTitleDoesNotOverrideSetDirty(t *testing.T) {
	screen := tcell.NewSimulationScreen("UTF-8")
	if err := screen.Init(); err != nil {
		t.Fatalf("failed to init simulation screen: %v", err)
	}
	defer screen.Fini()

	title := NewTitle("TestDB")
	vp := views.NewViewPort(screen, 0, 0, 80, 1)
	title.SetView(vp)

	title.SetTitle("MyTitle")
	title.Draw()
	initial := firstLine(screen)
	if strings.Contains(initial, "[MODIFIED]") {
		t.Fatalf("expected no [MODIFIED] badge initially, got: %q", initial)
	}

	title.SetDirty(true)
	screen.Clear()
	title.Draw()
	afterDirty := firstLine(screen)
	if !strings.Contains(afterDirty, "[MODIFIED]") {
		t.Fatalf("expected [MODIFIED] badge after SetDirty(true), got: %q", afterDirty)
	}

	title.SetTitle("UpdatedTitle")
	screen.Clear()
	title.Draw()
	afterTitleChange := firstLine(screen)
	if !strings.Contains(afterTitleChange, "[MODIFIED]") {
		t.Fatalf("expected [MODIFIED] badge to persist after SetTitle, got: %q", afterTitleChange)
	}
	if !strings.Contains(afterTitleChange, "UpdatedTitle") {
		t.Fatalf("expected 'UpdatedTitle' to appear after SetTitle, got: %q", afterTitleChange)
	}

	title.SetDirty(false)
	screen.Clear()
	title.Draw()
	afterClean := firstLine(screen)
	if strings.Contains(afterClean, "[MODIFIED]") {
		t.Fatalf("expected no [MODIFIED] badge after SetDirty(false), got: %q", afterClean)
	}
	if !strings.Contains(afterClean, "UpdatedTitle") {
		t.Fatalf("expected 'UpdatedTitle' to still appear, got: %q", afterClean)
	}
}
