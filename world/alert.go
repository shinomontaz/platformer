package world

import (
	"fmt"
	"math/rand"
	"platformer/common"

	"github.com/shinomontaz/pixel"
	"github.com/shinomontaz/pixel/text"
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

func InitAlerts() {
	alerts = make([]*Alert, 0)
	fnt := common.GetFont("regular")
	atlas = text.NewAtlas(fnt, text.ASCII)
}

func addAlert(pos pixel.Vec, force float64) *Alert {
	rect := pixel.R(pos.X-force, pos.Y-force, pos.X+force, pos.Y+force)
	al := &Alert{
		rect: rect,
		ttl:  1.0,
		txt:  randSeq([]rune("#$%&@*?arlTVXx"), 2+rand.Intn(3)) + "!",
	}
	alerts = append(alerts, al)
	return al
}

func updateAlerts(dt float64) {
	i := 0
	for _, al := range alerts {
		al.timer += dt
		al.rect = al.rect.Moved(pixel.Vec{-10 * dt, 10 * dt})
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
