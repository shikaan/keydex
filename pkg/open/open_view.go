package open

import (
	"fmt"
	"time"

	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
	"github.com/shikaan/kpcli/pkg/kdbx"
)

func header(database kdbx.Database) tview.Primitive {
	view := tview.NewTextView()

	view.SetBackgroundColor(tview.Styles.ContrastBackgroundColor)
	view.SetTextAlign(tview.AlignCenter)

	fmt.Fprintf(view, "kpcli: %s", database.Name)

	return view
}

func footer() tview.Primitive {
	view := tview.NewTextView()

	view.SetBackgroundColor(tview.Styles.MoreContrastBackgroundColor)
	fmt.Fprint(view, "kpcli: use arrows to navigate, Ctrl+ww to switch between panes")

	return view
}

func sidebar(groups []kdbx.Group, selectEntry func(kdbx.Entry)) *tview.TreeView {
	root := tview.NewTreeNode(groups[0].Name)
	tree := makeTree(root, groups[0].Groups, groups[0].Entries)

	tree.SetGraphics(false)

	tree.SetSelectedFunc(func(node *tview.TreeNode) {
		switch node.GetReference().(type) {

		case kdbx.Entry:
			selectEntry(node.GetReference().(kdbx.Entry))
			break
		case kdbx.Group:
			node.SetExpanded(!node.IsExpanded())

			newTitle := string([]rune(node.GetText())[1:])

			if node.IsExpanded() {
				node.SetText("▼" + newTitle)
			} else {
				node.SetText("▶" + newTitle)
			}
		}
	})

	//TODO: hoe to ighlight it's in use? idea: blinking cursor
	tree.SetFocusFunc(func() { tree.SetBackgroundColor(tcell.ColorGray) })
	tree.SetBlurFunc(func() { tree.SetBackgroundColor(tcell.ColorBlack) })

	return tree
}

func makeTree(root *tview.TreeNode, groups []kdbx.Group, entries []kdbx.Entry) *tview.TreeView {
	tree := tview.NewTreeView().SetRoot(root).SetCurrentNode(root)

	for _, g := range groups {
		node := tview.NewTreeNode("▶ " + g.Name)
		node.SetReference(g)
		node.SetExpanded(false)
		makeTree(node, g.Groups, g.Entries)
		root.AddChild(node)
	}

	for _, e := range entries {
		node := tview.NewTreeNode("  " + e.GetTitle())
		node.SetReference(e)
		node.SetSelectable(true)
		root.AddChild(node)
	}

	return tree
}

func main() tview.Primitive {
	pages := tview.NewPages()

	help := tview.NewTextView()

	fmt.Fprint(help, `
    Welcome to kpcli!

    Select an entry in the sidebar to see the details.
    `)

	pages.AddAndSwitchToPage("help", help, true)
	return pages
}

func entryPage(e kdbx.Entry, onSelect func(content string)) tview.Primitive {
	form := tview.NewForm()

	pwdIndex := e.GetPasswordIndex()

	for idx, v := range e.Values {
		content := e.GetContent(v.Key)

		// NOt painting empty field. Might become a problem when
		// adding a new field with the same label?
		if content != "" {
			field := tview.NewInputField()

			field.SetLabel(v.Key)
			field.SetText(content)
			field.SetFieldWidth(0)

			if idx == pwdIndex {
				field.SetMaskCharacter('*')
			}

			field.SetAcceptanceFunc(func(textToCheck string, lastChar rune) bool { return false })

			field.SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
				if event.Key() == tcell.KeyCtrlW {
					return nil
				}
				if event.Key() == tcell.KeyEnter {
					onSelect(content)
					return nil
				}
				if event.Key() == tcell.KeyUp {
					i, _ := form.GetFocusedItemIndex()
					form.SetFocus(i - 1)
				}

				return event
			})

			form.AddFormItem(field)
		}
	}

	return form
}

func Render(database kdbx.Database) error {
	app := tview.NewApplication()
	root := tview.NewPages()
	flex := tview.NewFlex()

	// TODO: make me a toaster
	modal := tview.NewModal()
	modal.SetText("Copied to clipboard")
	modal.AddButtons([]string{"ok"})

	root.
		AddPage("main", flex, true, true).
		AddPage("modal", modal, true, false)

	m := main().(*tview.Pages)

	modal.SetDoneFunc(func(buttonIndex int, buttonLabel string) {
		root.HidePage("modal")
		app.SetFocus(m)
	})

	onCopyFieldContent := func(fieldContent string) {
		root.ShowPage("modal")
	}

	onSelectedEntry := func(e kdbx.Entry) {
		id := e.GetTitle()

		if m.HasPage(id) {
			m.SwitchToPage(id)
		} else {
			m.AddAndSwitchToPage(id, entryPage(e, onCopyFieldContent), true)
		}

		app.SetFocus(m)
	}

	sb := sidebar(database.Groups, onSelectedEntry)

	flex.SetDirection(tview.FlexRow).
		AddItem(header(database), 1, 0, false).
		AddItem(tview.NewFlex().SetDirection(tview.FlexColumn).
			AddItem(sb, 0, 3, true).
			AddItem(m, 0, 9, false),
			0,
			11,
			false).
		AddItem(footer(), 1, 0, false)

	var isSelecting bool
	app.SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
		if event.Key() == tcell.KeyCtrlW {
			if isSelecting {
				if sb.HasFocus() {
					app.SetFocus(m)
				} else {
					app.SetFocus(sb)
				}
			} else {
				isSelecting = true
				go func() {
					time.Sleep(1 * time.Second)
					isSelecting = false
				}()
			}
			return nil
		}

		return event
	})

	return app.SetRoot(root, true).SetFocus(sb).Run()
}
