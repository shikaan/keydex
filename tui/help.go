package tui

import (
	"strings"

	"github.com/gdamore/tcell/v2"
	"github.com/gdamore/tcell/v2/views"
	"github.com/shikaan/kpcli/pkg/info"
	"github.com/shikaan/kpcli/tui/components"
)

// Line breaking the border is to accommodate for the space introduced
// with the useage of the constant. Do not align it!
const welcomeBanner = `
        ┌─────────────────────────────────────────────────────────────┐
        │ Welcome to ` + info.NAME +`!                                           │
        │                                                             │  
        │   Press Ctrl+P to start browsing the database or keep on    │
        │   reading these instructions for more information.          │
        └─────────────────────────────────────────────────────────────┘
`

func NewHelpView (screen tcell.Screen) views.Widget {
  c := components.NewContainer(screen)
  t := views.NewSimpleStyledText()

  // This line is showed only when reference is missing
  // as that signifies this is the first time in this 
  // session this page has been opened
  var firstAccessLine string

  if App.State.Reference == "" {
    firstAccessLine = welcomeBanner
  }

  t.SetText(`

` + firstAccessLine + `

` + strings.Title(info.NAME) + ` Help Text

` + strings.Title(info.NAME) + ` is designed to be a simple, easy-to-use, terminal-based password manager 
for the KeePass (https://keepass.info/) database format. The user interface is
highly inspired to GNU Nano (https://www.nano-editor.org/).

The top line displays the current version and contextual information about the
current view. At the bottom there are three lines: the two at the bottom are a
list of available commands, the third is a notification line - used to report
informational messages.

All the commands - except for navigation - are issued by pressing a combination
of Ctrl and another letter. We use the caret (^) symbol to indicate Ctrl. For 
example, ^C means Ctrl+C.

The following functions are available in ` + info.NAME +`:

^X    Closes the application. Asks for confirmation with unsaved changes.

^P    Opens a fuzzy finder to open other entries in the database.  

^O    Saves current state to the opened file.

^C    If a field is focused, copies the content to the clipboard.

^R    Reveals hidden fields such as password.

^G    Opens this help.

End of help.
`)

  c.SetContent(t)

  return c
}
