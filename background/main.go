package background

import (
	"math"
	"platformer/common"

	"github.com/shinomontaz/pixel"
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

func (b *Back) Draw(t pixel.Target, pos pixel.Vec, viewport pixel.Rect) {
	x, _ := pos.XY()
	steps := int(math.Abs(x / b.width))

	width := viewport.W()
	height := viewport.H()
	x, y := viewport.Min.X, viewport.Min.Y
	b.vector1 = pixel.V(x+width/2+1, y+height/2)
	b.vector2 = pixel.V(x+3*width/2, y+height/2)

	if steps != b.steps {
		b.vector1, b.vector2 = b.vector2, b.vector1

		b.steps = steps
	}

	x = math.Mod(x, b.width)

	b.pos = pixel.V(x, 0)

	mtx1 := pixel.IM.ScaledXY(pixel.ZV, pixel.V(
		b.width/b.part1.Frame().W(),
		b.height/b.part1.Frame().H(),
	)).Moved(b.vector1.Sub(b.pos))
	mtx2 := pixel.IM.ScaledXY(pixel.ZV, pixel.V(
		b.width/b.part2.Frame().W(),
		b.height/b.part2.Frame().H(),
	)).Moved(b.vector2.Sub(b.pos))
	b.part1.Draw(t, mtx1)
	b.part2.Draw(t, mtx2)
}
