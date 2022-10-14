package menu

import (
	"image/color"
	"math"

	"github.com/shinomontaz/pixel"
	"github.com/shinomontaz/pixel/pixelgl"
	"golang.org/x/image/colornames"
)

var (
	uTime float32
)

type Back struct {
	rgba          color.Color
	animSpriteNum int
	currtime      float64
	rect          pixel.Rect
	intercanvas   *pixelgl.Canvas
	canvas        *pixelgl.Canvas
	bg            *pixel.Sprite
}

func NewBack(rect pixel.Rect) *Back {
	b := &Back{
		intercanvas: pixelgl.NewCanvas(pixel.R(0, 0, rect.W(), rect.H())),
		canvas:      pixelgl.NewCanvas(pixel.R(0, 0, rect.W(), rect.H())),

		rect: rect,
		rgba: colornames.Black,
	}

	b.init()

	return b
}

func (b *Back) init() {
	bg, err := loader.LoadPicture("gamebackground.png")
	if err != nil {
		panic(err)
	}

	b.bg = pixel.NewSprite(bg, pixel.R(0, 0, bg.Bounds().W(), bg.Bounds().H()))
}

func (b *Back) Update(dt float64) {
	b.currtime += dt
	b.animSpriteNum = int(math.Floor(b.currtime / 0.2))

	uTime += float32(dt)
}

func (b *Back) Draw(t pixel.Target) {
	b.intercanvas.Clear(pixel.RGB(0, 0, 0))
	b.canvas.Clear(pixel.RGB(0, 0, 0))

	b.bg.Draw(b.intercanvas, pixel.IM.Moved(b.intercanvas.Bounds().Center()))

	b.intercanvas.Draw(b.canvas, pixel.IM.Moved(b.canvas.Bounds().Center()))
	b.canvas.Draw(t, pixel.IM.Moved(b.rect.Center()))
}
