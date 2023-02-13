package components

import "github.com/gdamore/tcell/v2/views"

type Form struct {
	WithFocusables
}

func NewForm() *Form {
	f := &Form{}
	f.SetOrientation(views.Vertical)

	return f
}
