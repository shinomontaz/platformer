package talks

import (
	"fmt"
	"image/color"

	"github.com/shinomontaz/pixel"
	"github.com/shinomontaz/pixel/pixelgl"
	"github.com/shinomontaz/pixel/text"
)

type Alert struct {
	rect  pixel.Rect
	timer float64
	ttl   float64
	txt   string
	col   color.Color
}

func (a *Alert) GetRect() pixel.Rect {
	return a.rect
}

func (a *Alert) Update(dt float64) bool {
	a.timer += dt
	a.rect = a.rect.Moved(pixel.Vec{-10 * dt, 10 * dt})
	return a.timer < a.ttl
}

func (a *Alert) Draw(win *pixelgl.Window, camPos, center pixel.Vec) {
	pos := a.rect.Center().Add(pixel.Vec{0.0, 40.0})
	// draw exclamation sign
	txt := text.New(pos, atlas)
	txt.Color = a.col
	fmt.Fprintln(txt, a.txt)
	txt.Draw(win, pixel.IM.Moved(center.Sub(camPos)))
}
