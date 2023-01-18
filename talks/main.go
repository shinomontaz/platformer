package talks

import (
	"math/rand"
	"platformer/activities"
	"platformer/common"

	"github.com/shinomontaz/pixel"
	"github.com/shinomontaz/pixel/pixelgl"
	"github.com/shinomontaz/pixel/text"
	"golang.org/x/image/colornames"
)

var alerts []*Alert
var atlas *text.Atlas

func Init() {
	alerts = make([]*Alert, 0)
	fnt := common.GetFont("regular")
	atlas = text.NewAtlas(fnt, text.ASCII)
}

func AddAlert(pos pixel.Vec, force float64) {
	rect := pixel.R(pos.X-force, pos.Y-force, pos.X+force, pos.Y+force)
	al := &Alert{
		rect: rect,
		ttl:  1.0,
		col:  colornames.Red,
		txt:  randSeq([]rune("#$%&@*?arlTVXx"), 2+rand.Intn(3)) + "!",
	}
	alerts = append(alerts, al)

	activities.Alert(al.GetRect())
}

func Update(dt float64) {
	i := 0
	for _, al := range alerts {
		if al.Update(dt) {
			alerts[i] = al
			i++
		}
	}
	alerts = alerts[:i]
}

func Draw(win *pixelgl.Window, camPos, center pixel.Vec) {
	for _, al := range alerts {
		al.Draw(win, camPos, center)
	}
}

//var letters = []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ")
func randSeq(letters []rune, n int) string {
	b := make([]rune, n)
	for i := range b {
		b[i] = letters[rand.Intn(len(letters))]
	}
	return string(b)
}
