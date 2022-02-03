package world

import (
	"fmt"

	"github.com/faiface/pixel"
	"github.com/faiface/pixel/imdraw"
	"github.com/faiface/pixel/text"
	"golang.org/x/image/colornames"
	"golang.org/x/image/font/basicfont"
)

type Alert struct {
	rect  pixel.Rect
	timer float64
	ttl   float64
}

var alerts []*Alert
var atlas *text.Atlas

func init() {
	alerts = make([]*Alert, 0)
	atlas = text.NewAtlas(basicfont.Face7x13, text.ASCII)
}

func addAlert(pos pixel.Vec, force float64) *Alert {
	rect := pixel.R(pos.X-force, pos.Y-force, pos.X+force, pos.Y+force)
	al := &Alert{
		rect: rect,
		ttl:  2.0,
	}
	alerts = append(alerts, al)
	return al
}

func updateAlerts(dt float64) {
	i := 0
	for _, al := range alerts {
		al.timer += dt
		if al.timer < al.ttl {
			alerts[i] = al
			i++
		}
	}
	alerts = alerts[:i]
}

func drawAlerts(t pixel.Target) {
	imd := imdraw.New(nil)
	imd.Color = colornames.Whitesmoke
	for _, al := range alerts {
		pos := al.rect.Center().Add(pixel.Vec{0.0, 40.0})
		imd.Push(pos.Add(pixel.Vec{5, 8}))
		imd.Circle(10, 0)

		// draw exclamation sign
		txt := text.New(pos, atlas)
		txt.Color = colornames.Red
		fmt.Fprintln(txt, "!")
		imd.Draw(t)
		txt.Draw(t, pixel.IM.Scaled(txt.Orig, 1.5))
	}
}

func (a *Alert) GetRect() pixel.Rect {
	return a.rect
}
