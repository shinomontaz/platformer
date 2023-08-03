package background

import (
	"image/color"
	"platformer/common"

	"github.com/shinomontaz/pixel"
	"github.com/shinomontaz/pixel/imdraw"
)

type Pback struct {
	width  float64
	height float64
	layers []*Back
	imd    *imdraw.IMDraw
}

func NewParallax(start pixel.Vec, viewport pixel.Rect, loader *common.Loader) *Pback {
	width := viewport.W()
	height := viewport.H()

	p := Pback{
		width:  width,
		height: height,
		layers: make([]*Back, 0),
	}
	p.layers = append(p.layers, New(start, viewport, loader, "back/rocs1/1.png", WithSpeed(0)))
	p.layers = append(p.layers, New(start, viewport, loader, "back/rocs1/2.png", WithSpeed(0.1)))
	p.layers = append(p.layers, New(start, viewport, loader, "back/rocs1/3.png", WithSpeed(0.2), WithOffset(pixel.Vec{0, 150})))
	p.layers = append(p.layers, New(start, viewport, loader, "back/rocs1/4.png", WithSpeed(0.7), WithOffset(pixel.Vec{0, 100})))
	p.layers = append(p.layers, New(start, viewport, loader, "back/rocs1/5.png", WithOffset(pixel.Vec{0, 50})))

	// p.layers = append(p.layers, New(start, viewport, loader, "back/nature/2.png", WithSpeed(0.1)))
	// p.layers = append(p.layers, New(start, viewport, loader, "back/nature/3.png", WithSpeed(0.7)))
	// p.layers = append(p.layers, New(start, viewport, loader, "back/nature/4.png"))

	p.imd = imdraw.New(nil)
	//	p.imd.Color = color.RGBA{101, 186, 227, 1} // nature

	return &p
}

func (p *Pback) Update(dt float64, pos pixel.Vec) {
	p.imd.Clear()
	p.imd.Color = color.RGBA{207, 199, 223, 1} // rocs1
	p.imd.Push(pixel.Vec{pos.X - p.width/2, pos.Y - p.height/2 + 150}, pixel.Vec{pos.X + p.width, pos.Y + p.height/2 + 150})
	p.imd.Rectangle(0)

	for _, l := range p.layers {
		l.Update(dt, pos)
	}
}

func (p *Pback) Draw(t pixel.Target) {
	p.imd.Draw(t)
	for _, l := range p.layers {
		l.Draw(t)
	}
}
