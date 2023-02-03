package talks

import (
	"fmt"
	"image/color"
	"unicode"

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
	txt := text.New(pos, atlas)
	txt.LineHeight = atlas.LineHeight() * 1.3
	txt.Color = a.col
	strChunks := splitToChunks(a.txt, 10) // use not a long strings and make it wordwrapped
	for _, st := range strChunks {
		fmt.Fprintln(txt, st)
	}
	txt.Draw(win, pixel.IM.Moved(center.Sub(camPos)))
}

func splitToChunks(s string, l int) []string {
	if len(s) <= l {
		return []string{s}
	}

	currStart := 0
	currLen := 0
	res := make([]string, 0)
	for i := range s {
		if currLen >= l && unicode.IsSpace(rune(s[i])) {
			res = append(res, s[currStart:i+1])
			currLen = 0
			currStart = i + 1
		}
		currLen++
	}

	res = append(res, s[currStart:])
	return res
}
