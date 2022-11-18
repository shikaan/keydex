package tui

import (
	"fmt"

	"github.com/gdamore/tcell/v2"
	"github.com/gdamore/tcell/v2/views"
	"github.com/lithammer/fuzzysearch/fuzzy"

	"github.com/shikaan/kpcli/pkg/tui/components"
	"github.com/shikaan/kpcli/pkg/utils"
)

type ListView struct {
	views.Panel
}

func NewListView(screen tcell.Screen, props State) views.Widget {
	container := components.NewContainer(screen)
	view := &ListView{}
	paths := props.Database.GetEntryPaths()
	count := len(paths)
	maxX, maxY := getBoundaries(screen)
	list := newList(paths, count, maxX, maxY, screen, props)

	options := &components.InputOptions{}
	input := components.NewInput(options)
	input.SetFocus(true)

	view.AddWidget(input, 0)
	view.AddWidget(list, 0)

	input.OnChange(func(ev tcell.Event) bool {
		entries := paths
		content := input.GetContent()

		if content != "" {
			entries = fuzzy.FindFold(content, paths)
			view.RemoveWidget(list)
			list = newList(entries, count, maxX, maxY, screen, props)
			view.AddWidget(list, 0)
			view.Resize()
		}

		return true
	})

	container.SetContent(view)
	return container
}

type List struct {
  views.BoxLayout
}


func (f *List) HandleEvent(ev tcell.Event) bool {
	switch ev := ev.(type) {
	case *tcell.EventKey:
		if ev.Key() == tcell.KeyUp {
      f.MoveFocus(-1)
      return true
    }
    if ev.Key() == tcell.KeyDown {
      f.MoveFocus(1)
      return true
    }
  }
	
  return f.BoxLayout.HandleEvent(ev)
}

// Moves focus by `offset` fields
func (f *List) MoveFocus (offset int) {
  fs := f.Focusables()
  count := len(fs)
  current := -1

  for i, f := range fs {
    if f.HasFocus() {
      current = i
      f.SetFocus(false)
      break
    }
  }

  notFound := current == -1
  if notFound {
    return 
  }

  newIndex := (count + current + offset) % count

  fs[newIndex].SetFocus(true)
}

// Returns the subset of Widgets that can have focus
func (f *List) Focusables() []components.Focusable {
  ws := f.Widgets()
  result := []components.Focusable{}

  for _, w := range ws {
    switch w := w.(type) {
    case components.Focusable:
      result = append(result, w)
    }
  }
  return result
}

func newList(matchedEntries []string, total, maxX int, maxY int, screen tcell.Screen, state State) views.Widget {
	box := &List{}
  box.SetOrientation(views.Vertical)

	line := views.NewSimpleStyledText()
	// line.SetStyle(line.Style().Reverse(true))
	line.SetAlignment(views.HAlignRight)
  counter := fmt.Sprintf("%d/%d", len(matchedEntries), total)
	line.SetText(fmt.Sprintf("%*s", maxX, counter))

	box.AddWidget(line, 0)

	for i, e := range matchedEntries {
		if i >= maxY {
			break
		}

		line := components.NewOption()
		if len(e) > maxX {
			line.SetContent(e[:maxX])
		} else {
			line.SetContent(e)
		}

    entry := e
    line.OnSelect(func(ev *tcell.EventKey) bool {
      App.State.Entry = App.State.Database.GetEntry(entry)
      App.NavigateTo(NewEditView)

      return true
    })

		box.AddWidget(line, 0)
    if i == 0 {
      line.SetFocus(true)
    }
	}

	return box
}

func getBoundaries(screen tcell.Screen) (int, int) {
	x, y := (screen).Size()

	// one third of the screen width
	// all the height - title, status, search, counter, and notification
	return utils.Max(x / 3, components.MIN_WIDTH), y - 6
}
