package submenu

import (
	"github.com/faiface/pixel"
	"github.com/faiface/pixel/text"
)

const (
	MAIN = iota
	VIDEO
	SOUND
	INFO
)

type Item struct {
	title    string
	action   func()
	selected bool
	txt      *text.Text
}

type Menu interface {
	SetState(int)
	GetRect() pixel.Rect
	OnVideo(isfullscreen bool)
}
