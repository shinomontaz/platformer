package background

import (
	"math"
	"platformer/common"

	"github.com/faiface/pixel"
)

type Back struct {
	width    float64
	height   float64
	p        pixel.Picture
	viewport pixel.Rect
	vector1  pixel.Vec
	vector2  pixel.Vec

	part1 *pixel.Sprite
	part2 *pixel.Sprite

	pos   pixel.Vec
	steps int
}

func New(start pixel.Vec, viewport pixel.Rect, path string) *Back {
	width := viewport.W()
	height := viewport.H()

	x, y := viewport.Min.X, viewport.Min.Y
	b := Back{
		width:    width,
		height:   height,
		viewport: viewport,

		vector1: pixel.V(x+width/2+1, y+height/2),
		vector2: pixel.V(x+3*width/2, y+height/2),

		pos:   start,
		steps: 0,
	}

	bg, err := common.LoadPicture(path)
	if err != nil {
		panic(err)
	}

	b.p = bg

	b.part1 = pixel.NewSprite(b.p, pixel.R(0, 0, bg.Bounds().W()/2, bg.Bounds().H()))
	b.part2 = pixel.NewSprite(b.p, pixel.R(bg.Bounds().W()/2, 0, bg.Bounds().W(), bg.Bounds().H()))

	return &b
}

func (b *Back) Draw(t pixel.Target, pos pixel.Vec, cam pixel.Vec) {
	x, _ := pos.XY()
	cam.Y = 0
	steps := int(math.Abs(x / b.width))

	if steps != b.steps {
		b.vector1, b.vector2 = b.vector2, b.vector1

		b.steps = steps
	}

	x = math.Mod(x, b.width)

	b.pos = pixel.V(x, 0)

	// b.part1.Draw(t, pixel.IM.Moved(b.vector1.Sub(b.pos).Sub(cam)).ScaledXY(b.vector1.Sub(b.pos).Sub(cam), pixel.Vec{b.width / b.p.Bounds().W() / 2, b.height / b.p.Bounds().H()}))
	// b.part2.Draw(t, pixel.IM.Moved(b.vector2.Sub(b.pos).Sub(cam)).ScaledXY(b.vector2.Sub(b.pos).Sub(cam), pixel.Vec{b.width / b.p.Bounds().W() / 2, b.height / b.p.Bounds().H()}))

	b.part1.Draw(t, pixel.IM.Moved(b.vector1.Sub(b.pos).Sub(cam)))
	b.part2.Draw(t, pixel.IM.Moved(b.vector2.Sub(b.pos).Sub(cam)))
}
