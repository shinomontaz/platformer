package background

import (
	"math"
	"platformer/common"

	"github.com/shinomontaz/pixel"
)

type Option func(*Back)

func WithSpeed(speed float64) Option {
	return func(b *Back) {
		b.speed = speed
	}
}

func WithOffset(vec pixel.Vec) Option {
	return func(b *Back) {
		b.offset = vec
	}
}

type Back struct {
	width   float64
	height  float64
	p       pixel.Picture
	vector1 pixel.Vec
	vector2 pixel.Vec

	part1 *pixel.Sprite
	part2 *pixel.Sprite

	pos        pixel.Vec
	start      pixel.Vec
	currCenter pixel.Vec
	offset     pixel.Vec
	steps      int
	speed      float64
}

func New(start pixel.Vec, viewport pixel.Rect, loader *common.Loader, path string, opts ...Option) *Back {
	width := viewport.W()
	height := viewport.H()

	b := Back{
		offset: pixel.ZV,
		width:  width,
		height: height,

		vector1: pixel.V(width/2+1, height/2),
		vector2: pixel.V(3*width/2-1, height/2),

		speed: 1,

		pos:   start,
		start: start,
		steps: 0,
	}

	for _, o := range opts {
		o(&b)
	}

	bg, err := loader.LoadPicture(path)

	if err != nil {
		panic(err)
	}

	b.p = bg

	b.part1 = pixel.NewSprite(b.p, pixel.R(0, 0, bg.Bounds().W()/2, bg.Bounds().H()))
	b.part2 = pixel.NewSprite(b.p, pixel.R(bg.Bounds().W()/2, 0, bg.Bounds().W(), bg.Bounds().H()))

	return &b
}

func (b *Back) Update(dt float64, pos pixel.Vec) {
	x := (pos.X - b.start.X) * b.speed
	steps := int(math.Abs(x / b.width))

	if steps != b.steps {
		b.vector1, b.vector2 = b.vector2, b.vector1
		b.steps = steps
	}

	y := (pos.Y - b.start.Y)
	b.pos = pixel.V(x-float64(steps)*b.width, y)

	b.currCenter = pos.Sub(pixel.Vec{b.width / 2, 150})
}

func (b *Back) Draw(t pixel.Target) {
	mtx1 := pixel.IM.ScaledXY(pixel.ZV, pixel.V(
		b.width/b.part1.Frame().W(),
		b.height/b.part1.Frame().H(),
	)).Moved(b.vector1.Sub(b.pos)).Moved(b.currCenter).Moved(b.offset)
	mtx2 := pixel.IM.ScaledXY(pixel.ZV, pixel.V(
		b.width/b.part2.Frame().W(),
		b.height/b.part2.Frame().H(),
	)).Moved(b.vector2.Sub(b.pos)).Moved(b.currCenter).Moved(b.offset)
	b.part1.Draw(t, mtx1)
	b.part2.Draw(t, mtx2)
}
