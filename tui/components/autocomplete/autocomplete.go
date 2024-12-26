package autocomplete

import (
	"fmt"
	"math"

	"github.com/gdamore/tcell/v2"
	"github.com/gdamore/tcell/v2/views"
	"github.com/lithammer/fuzzysearch/fuzzy"
	"github.com/mattn/go-runewidth"
	"github.com/shikaan/keydex/tui/components"
)

type AutoComplete struct {
	options AutoCompleteOptions

	list    *components.WithFocusables
	counter *views.SimpleStyledText
	search  *Search

	components.Focusable
	views.BoxLayout
}

type AutoCompleteOptions struct {
	Screen     tcell.Screen
	Entries    []string
	TotalCount int
	MaxX       int
	MaxY       int

	OnSelect func(entry string) bool
	OnFocus  func() bool
}

func NewAutoComplete(options AutoCompleteOptions) *AutoComplete {
	autoComplete := &AutoComplete{}
	autoComplete.SetOrientation(views.Vertical)
	autoComplete.options = options

	search := NewSearch()
	autoComplete.search = search
	autoComplete.AddWidget(search, 0)

	search.OnChange(func(ev tcell.Event) bool {
		entries := options.Entries
		content := search.GetContent()

		if len(content) > 0 {
			entries = fuzzy.FindFold(content, options.Entries)
		}

		autoComplete.drawList(entries)
		autoComplete.drawCounter(entries)

		return false
	})

	counter := views.NewSimpleStyledText()
	counter.SetAlignment(views.HAlignRight)
	autoComplete.counter = counter
	autoComplete.AddWidget(counter, 0)

	autoComplete.drawList(options.Entries)
	autoComplete.drawCounter(options.Entries)

	return autoComplete
}

func (ac *AutoComplete) OnFocus(cb func() bool) func() {
	return ac.search.OnFocus(cb)
}

func (ac *AutoComplete) SetFocus(on bool) {
	ac.search.SetFocus(on)
}

func (ac *AutoComplete) HasFocus() bool {
	return ac.search.HasFocus()
}

func (ac *AutoComplete) drawList(entries []string) {
	container := &components.WithFocusables{}
	container.SetOrientation(views.Vertical)

	// The 2 characters are used for when the entry is selected
	// and you have "> " prepended
	maxLineLength := ac.options.MaxX - 2

	if len(entries) == 0 {
		line := views.NewSimpleStyledText()
		line.SetText(runewidth.FillRight("--- No Results ---", ac.options.MaxX))
		container.AddWidget(line, 0)
	}

	for i, entry := range entries {
		if i >= ac.options.MaxY {
			break
		}

		line := newOption()
		line.SetContent((runewidth.FillRight(runewidth.Truncate(entry, maxLineLength, ""), maxLineLength)))

		// For memoization
		e := entry
		line.OnSelect(func() bool {
			return ac.options.OnSelect(e)
		})

		container.AddWidget(line, 0)
		if i == 0 {
			line.SetFocus(true)
		}
	}

	ac.RemoveWidget(ac.list)
	ac.list = container
	ac.AddWidget(ac.list, 0)
}

func (ac *AutoComplete) drawCounter(entries []string) {
	matched := int(math.Min(float64(len(entries)), float64(ac.options.MaxY)))
	counter := fmt.Sprintf("%d/%d", matched, len(ac.options.Entries))

	ac.counter.SetStyle(tcell.StyleDefault.Bold(matched == 0))

	ac.counter.SetText(fmt.Sprintf("% s", counter))
}
