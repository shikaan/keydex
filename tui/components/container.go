package components

import (
	"github.com/gdamore/tcell/v2"
	"github.com/gdamore/tcell/v2/views"
)

const MIN_WIDTH = 50

var spacer = &views.Spacer{}

type Container struct {
	Orientation views.Orientation

	screen tcell.Screen
	views.BoxLayout
}

func (c *Container) Draw() {
	w, _ := c.screen.Size()

	c.BoxLayout.RemoveWidget(spacer)
	c.BoxLayout.RemoveWidget(spacer)

	if w > MIN_WIDTH {
		c.BoxLayout.InsertWidget(0, spacer, 1)
		c.BoxLayout.InsertWidget(2, spacer, 1)
	}

	c.Resize()
	c.BoxLayout.Draw()
}

func (c *Container) SetContent(w views.Widget) {
	c.BoxLayout.InsertWidget(1, w, 1)
}

func (c *Container) InsertWidget() {}
func (c *Container) AddWidget()    {}
func (c *Container) RemoveWidget() {}

func NewContainer(screen tcell.Screen) *Container {
	container := &Container{}

	container.SetOrientation(views.Horizontal)
	container.screen = screen

	return container
}
