package background

import (
	"github.com/shinomontaz/pixel"
)

type Pback struct {
	width  float64
	height float64
	p1     *Back
	p2     *Back
	p3     *Back
	p4     *Back
	p5     *Back
	p6     *Back
	p7     *Back
}

func NewParallax(start pixel.Vec, viewport pixel.Rect) *Pback {
	width := viewport.W()
	height := viewport.H()

	p := Pback{
		width:  width,
		height: height,
	}
	p.p1 = New(start, viewport, "assets/back/1.png")
	p.p2 = New(start, viewport, "assets/back/2.png")
	p.p3 = New(start, viewport, "assets/back/3.png")
	p.p4 = New(start, viewport, "assets/back/4.png")
	p.p5 = New(start, viewport, "assets/back/5.png")
	p.p6 = New(start, viewport, "assets/back/6.png")
	p.p7 = New(start, viewport, "assets/back/8.png")

	//background.New(lastPos, currBounds.Moved(pixel.Vec{0, 100}), "assets/gamebackground.png")
	return &p
}

func (p *Pback) Draw(t pixel.Target, pos pixel.Vec, viewport pixel.Rect) {
	// p.p1.Draw(t, pos, pixel.Vec{cam.X * 0.01, cam.Y})
	// p.p2.Draw(t, pos, pixel.Vec{cam.X * 0.05, cam.Y})
	// p.p3.Draw(t, pos, pixel.Vec{cam.X * 0.1, cam.Y})
	// p.p4.Draw(t, pos, pixel.Vec{cam.X * 0.2, cam.Y})
	// p.p5.Draw(t, pos, pixel.Vec{cam.X * 0.5, cam.Y})
	// p.p6.Draw(t, pos, pixel.Vec{cam.X * 0.8, cam.Y})
	// p.p7.Draw(t, pos, pixel.Vec{cam.X * 1, cam.Y})

	p.p1.Draw(t, pos, viewport)
	p.p2.Draw(t, pos, viewport)
	p.p3.Draw(t, pos, viewport)
	// p.p4.Draw(t, pos, cam)
	// p.p5.Draw(t, pos, cam)
	// p.p6.Draw(t, pos, cam)
	p.p7.Draw(t, pos, viewport)

}
