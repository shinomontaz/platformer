package world

import (
	"fmt"
	"math/rand"
	"platformer/common"

	"github.com/faiface/pixel"
	"github.com/faiface/pixel/text"
	"golang.org/x/image/colornames"
)

type Alert struct {
	rect  pixel.Rect
	timer float64
	ttl   float64
	txt   string
}

var alerts []*Alert
var atlas *text.Atlas

func init() {
	alerts = make([]*Alert, 0)
	fnt, err := common.LoadFont("assets/fonts/GravityBold8.ttf", 8)
	if err != nil {
		panic(err)
	}
	atlas = text.NewAtlas(fnt, text.ASCII)
	//	atlas = text.NewAtlas(basicfont.Face7x13, text.ASCII)
}

func addAlert(pos pixel.Vec, force float64) *Alert {
	rect := pixel.R(pos.X-force, pos.Y-force, pos.X+force, pos.Y+force)
	al := &Alert{
		rect: rect,
		ttl:  2.0,
		txt:  randSeq([]rune("#$%&@*?arlTVXx"), 4) + "!",
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
	for _, al := range alerts {
		pos := al.rect.Center().Add(pixel.Vec{0.0, 40.0})
		// draw exclamation sign
		txt := text.New(pos, atlas)
		txt.Color = colornames.Red
		fmt.Fprintln(txt, al.txt)
		txt.Draw(t, pixel.IM)
	}
}

func (a *Alert) GetRect() pixel.Rect {
	return a.rect
}

//var letters = []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ")

func randSeq(letters []rune, n int) string {
	b := make([]rune, n)
	for i := range b {
		b[i] = letters[rand.Intn(len(letters))]
	}
	return string(b)
}
