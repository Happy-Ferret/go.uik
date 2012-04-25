package widgets

import (
	"code.google.com/p/draw2d/draw2d"
	"github.com/skelterjohn/geom"
	"github.com/skelterjohn/go.uik"
	"image/color"
)

type Checker chan bool

type Checkbox struct {
	uik.Block

	state, pressed, pressHover bool
}

func NewCheckbox(size geom.Coord) (c *Checkbox) {
	c = new(Checkbox)
	c.Initialize()
	c.Size = size

	go c.handleEvents()

	c.Paint = func(gc draw2d.GraphicContext) {
		c.draw(gc)
	}

	c.SetSizeHint(uik.SizeHint{
		MinSize:       size,
		PreferredSize: size,
		MaxSize:       size,
	})

	return
}

func (c *Checkbox) draw(gc draw2d.GraphicContext) {
	gc.Clear()
	gc.SetStrokeColor(color.Black)
	if c.pressed {
		if c.pressHover {
			gc.SetFillColor(color.RGBA{200, 0, 0, 255})
		} else {
			gc.SetFillColor(color.RGBA{155, 0, 0, 255})
		}
	} else {
		gc.SetFillColor(color.RGBA{255, 0, 0, 255})
	}

	// Draw background rect
	x, y := gc.LastPoint()
	gc.MoveTo(0, 0)
	gc.LineTo(c.Size.X, 0)
	gc.LineTo(c.Size.X, c.Size.Y)
	gc.LineTo(0, c.Size.Y)
	gc.Close()
	gc.FillStroke()

	// Draw inner rect
	if c.state {
		gc.SetFillColor(color.Black)
		gc.MoveTo(5, 5)
		gc.LineTo(c.Size.X-5, 5)
		gc.LineTo(c.Size.X-5, c.Size.Y-5)
		gc.LineTo(5, c.Size.Y-5)
		gc.Close()
		gc.FillStroke()
	}

	gc.MoveTo(x, y)
}

func (c *Checkbox) handleEvents() {
	for {
		select {
		case e := <-c.UserEvents:
			switch e := e.(type) {
			case uik.MouseDownEvent:
				c.pressed = true
				c.pressHover = true
				c.Invalidate()
			case uik.MouseUpEvent:
				if c.pressHover {
					c.state = !c.state
					c.Invalidate()
				}
				c.pressHover = false
				c.pressed = false
			case uik.MouseDraggedEvent:
				if !c.pressed {
					break
				}
				hover := c.Bounds().ContainsCoord(e.Where())
				// uik.Report(c.ID, "mde")
				if hover != c.pressHover {
					c.pressHover = hover
					c.Invalidate()
					// uik.Report("invalidate pressHover")
				}
			default:
				c.Block.HandleEvent(e)
			}
		}
	}
}
