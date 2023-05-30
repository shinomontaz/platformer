package talks

import (
	"fmt"
	"image/color"
	"unicode"

	"github.com/shinomontaz/pixel"
	"github.com/shinomontaz/pixel/pixelgl"
	"github.com/shinomontaz/pixel/text"
	"golang.org/x/image/colornames"
)

var borderColor = colornames.Whitesmoke

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
	pos := a.rect.Center().Add(pixel.Vec{-1.0, 40.0})
	txt := text.New(pos, atlas)
	txt.LineHeight = atlas.LineHeight() * 1.1
	txt.Color = borderColor
	strChunks := splitToChunks(a.txt, 10) // use not a long strings and make it wordwrapped
	for _, st := range strChunks {
		fmt.Fprintln(txt, st)
	}

	pos = a.rect.Center().Add(pixel.Vec{1.0, 40.0})
	txt2 := text.New(pos, atlas)
	txt2.LineHeight = atlas.LineHeight() * 1.1
	txt2.Color = borderColor
	for _, st := range strChunks {
		fmt.Fprintln(txt2, st)
	}

	pos = a.rect.Center().Add(pixel.Vec{0.0, 41.0})
	txt3 := text.New(pos, atlas)
	txt3.LineHeight = atlas.LineHeight() * 1.1
	txt3.Color = borderColor
	for _, st := range strChunks {
		fmt.Fprintln(txt3, st)
	}

	pos = a.rect.Center().Add(pixel.Vec{0.0, 39.0})
	txt4 := text.New(pos, atlas)
	txt4.LineHeight = atlas.LineHeight() * 1.1
	txt4.Color = borderColor
	for _, st := range strChunks {
		fmt.Fprintln(txt4, st)
	}

	pos = a.rect.Center().Add(pixel.Vec{0.0, 40.0})
	txt5 := text.New(pos, atlas)
	txt5.LineHeight = atlas.LineHeight() * 1.1
	txt5.Color = a.col
	for _, st := range strChunks {
		fmt.Fprintln(txt5, st)
	}

	txt.Draw(win, pixel.IM.Moved(center.Sub(camPos)))
	txt2.Draw(win, pixel.IM.Moved(center.Sub(camPos)))
	txt3.Draw(win, pixel.IM.Moved(center.Sub(camPos)))
	txt4.Draw(win, pixel.IM.Moved(center.Sub(camPos)))
	txt5.Draw(win, pixel.IM.Moved(center.Sub(camPos)))
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
