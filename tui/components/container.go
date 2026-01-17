package components

import (
	"sync"

	"github.com/gdamore/tcell/v2"
	"github.com/gdamore/tcell/v2/views"
	"github.com/shikaan/keydex/pkg/log"
)

const MIN_WIDTH = 80

type Container struct {
	view views.View

	wview  views.ViewPort
	widget views.Widget

	once sync.Once

	views.WidgetWatchers
}

func (c *Container) Resize() {
	log.Debug("Resize")
	c.initialize(c.widget)
	c.layout()
	c.widget.Resize()
	c.PostEventWidgetResize(c)
}

func (c *Container) SetView(view views.View) {
	log.Debug("SetView")
	c.initialize(c.widget)
	c.view = view
	c.wview.SetView(view)
}

func (c *Container) Draw() {
	c.initialize(c.widget)
	c.layout()

	w, h := c.view.Size()
	for y := range h {
		for x := range w {
			c.view.SetContent(x, y, ' ', nil, tcell.StyleDefault)
		}
	}

	c.widget.Draw()
}

func (c *Container) Size() (int, int) {
	return c.widget.Size()
}

func (c *Container) SetContent(w views.Widget) {
	c.initialize(w)
	c.PostEventWidgetContent(c)

	log.Debug("SetContent")
}

func (c *Container) HandleEvent(ev tcell.Event) bool {
	switch ev.(type) {
	case *views.EventWidgetContent:
		log.Debug("EventWidgetContent")
		c.PostEventWidgetContent(c)
		return true
	}

	if c.widget.HandleEvent(ev) {
		log.Debug("widget.HandleEvent")
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

	// TODO: make autocoplete to be of fixed size 80
	// TODO: make the form to be of fixed size smaller
	_, wh := c.widget.Size()
	c.wview.Resize((w-MIN_WIDTH)/2, 0, MIN_WIDTH, wh)
	c.widget.Resize()
}
