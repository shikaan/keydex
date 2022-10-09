package style

import (
	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
)

const (
	OpenTreeIcon   = "▼"
	ClosedTreeIcon = "▶"
)

var dark = tview.Theme{
	PrimitiveBackgroundColor:    tcell.ColorBlack,
	ContrastBackgroundColor:     tcell.ColorBlue,
	MoreContrastBackgroundColor: tcell.ColorGreen,
	BorderColor:                 tcell.ColorWhite,
	TitleColor:                  tcell.ColorWhite,
	GraphicsColor:               tcell.ColorWhite,
	PrimaryTextColor:            tcell.ColorWhite,
	SecondaryTextColor:          tcell.ColorYellow,
	TertiaryTextColor:           tcell.ColorGreen,
	InverseTextColor:            tcell.ColorBlue,
	ContrastSecondaryTextColor:  tcell.ColorDarkBlue,
}

var light = tview.Theme{
	PrimitiveBackgroundColor:    tcell.ColorWhite,
	ContrastBackgroundColor:     tcell.ColorSilver,
	MoreContrastBackgroundColor: tcell.ColorGreen,
	BorderColor:                 tcell.ColorBlack,
	TitleColor:                  tcell.ColorBlack,
	GraphicsColor:               tcell.ColorBlack,
	PrimaryTextColor:            tcell.ColorBlack,
	SecondaryTextColor:          tcell.ColorYellow,
	TertiaryTextColor:           tcell.ColorGreen,
	InverseTextColor:            tcell.ColorBlue,
	ContrastSecondaryTextColor:  tcell.ColorDarkBlue,
}

type Theme = string

const (
	Light Theme = "light"
	Dark  Theme = "dark"
)

func SetTheme(theme Theme) {
	if theme == Dark {
		tview.Styles = dark
		return
	}

	tview.Styles = light
}
