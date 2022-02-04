package ui

import "github.com/faiface/pixel"

type Characterer interface {
	GetHp() int
	GetPortrait() *pixel.Sprite
}
