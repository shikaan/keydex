package components

import (
	"sync"

	"github.com/gdamore/tcell/v2"
	"github.com/gdamore/tcell/v2/views"
)

const CONTENT_WIDTH = 76

type Container struct {
	view views.View

	wview  views.ViewPort
	widget views.Widget

	once sync.Once

	views.WidgetWatchers
}

func (c *Container) Resize() {
	c.initialize(c.widget)
	c.layout()
	c.widget.Resize()
	c.PostEventWidgetResize(c)
}

func (c *Container) SetView(view views.View) {
	c.initialize(c.widget)
	c.view = view
	c.wview.SetView(view)
}

func (c *Container) Draw() {
	c.initialize(c.widget)
	c.layout()

	c.view.Clear()
	c.widget.Draw()
}

func (c *Container) Size() (int, int) {
	return c.widget.Size()
}

func (c *Container) SetContent(w views.Widget) {
	c.initialize(w)
	c.PostEventWidgetContent(c)
}

func (c *Container) HandleEvent(ev tcell.Event) bool {
	switch ev.(type) {
	case *views.EventWidgetContent:
		c.PostEventWidgetContent(c)
		return true
	}

	if c.widget.HandleEvent(ev) {
		return true
	}

	return false
}

func (c *Container) initialize(w views.Widget) {
	c.widget = w
	c.once.Do(func() {
		c.widget.SetView(&c.wview)
		c.widget.Watch(c)
	})
}

func (c *Container) layout() {
	w, _ := c.view.Size()

	_, wh := c.widget.Size()
	c.wview.Resize((w-CONTENT_WIDTH)/2, 0, CONTENT_WIDTH, wh)
	c.widget.Resize()
}
