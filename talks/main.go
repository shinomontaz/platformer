package talks

import (
	"image/color"
	"math"
	"platformer/activities"
	"platformer/common"

	"github.com/shinomontaz/pixel"
	"github.com/shinomontaz/pixel/pixelgl"
	"github.com/shinomontaz/pixel/text"
	"golang.org/x/image/colornames"
)

var alerts []*Alert
var atlas *text.Atlas

func Init(loader *common.Loader) {
	alerts = make([]*Alert, 0)
	fnt := common.GetFont("regular16")
	atlas = text.NewAtlas(fnt, text.ASCII)
	initPhrases(loader)
}

func AddAlert(pos pixel.Vec, force float64) {
	al := addAlert(pos, colornames.Red, randSeq([]rune("#$%&@*?arlTVXx"), 2+int(common.GetRandFloat()*3))+"!", 1, force)
	alerts = append(alerts, al)
	activities.Alert(al.GetRect())
}

func addAlert(pos pixel.Vec, col color.Color, txt string, ttl, force float64) *Alert {
	return &Alert{
		rect: pixel.R(pos.X-force, pos.Y-force, pos.X+force, pos.Y+force),
		ttl:  ttl,
		col:  col,
		txt:  txt,
	}
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
		b[i] = letters[int(math.Round(common.GetRandFloat()*float64(len(letters)-1)))]
	}
	return string(b)
}
