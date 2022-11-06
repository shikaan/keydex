package tui

import "github.com/gdamore/tcell/v2/views"

const MIN_HEIGHT = 50
const MIN_WIDTH = 50

var spacer = &views.Spacer{}

type Responsive struct {
  Orientation views.Orientation

  views.BoxLayout
}

func (r *Responsive) Draw() {
  w, h := Screen.Size()

  r.BoxLayout.RemoveWidget(spacer)
  r.BoxLayout.RemoveWidget(spacer)
  
  isHorizontallyBigEnough := w > MIN_WIDTH && r.Orientation == views.Horizontal
  isVerticallyBigEnough := h > MIN_HEIGHT && r.Orientation == views.Vertical

  if isHorizontallyBigEnough || isVerticallyBigEnough {
    r.BoxLayout.InsertWidget(0, spacer, 1)
    r.BoxLayout.InsertWidget(2, spacer, 1)
  } 

  r.Resize()
  r.BoxLayout.Draw()
}

func (r *Responsive) SetContent (w views.Widget) {
  r.BoxLayout.InsertWidget(1, w, 1)
}

func (r *Responsive) InsertWidget () {}
func (r *Responsive) AddWidget () {}
func (r *Responsive) RemoveWidget () {}

func (r *Responsive) SetOrientation(o views.Orientation) {
  r.Orientation = o
  r.BoxLayout.SetOrientation(o)
}

func NewResponsive(o views.Orientation) *Responsive {
  f := &Responsive{}
  f.SetOrientation(o)

  return f
}
