package ui

import "github.com/shinomontaz/pixel"

type Characterer interface {
	GetHp() int
	GetPortrait() *pixel.Sprite
}
