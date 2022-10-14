package ui

import (
	"platformer/common"

	"github.com/shinomontaz/pixel"
)

type Ui struct {
	ch       Characterer
	viewport pixel.Rect
	heart    *pixel.Sprite
}

var loader *common.Loader

func Init(l *common.Loader) {
	loader = l
}

func New(ch Characterer, viewport pixel.Rect) *Ui {
	icon, err := loader.LoadPicture("icons/37.png")
	if err != nil {
		panic(err)
	}

	heart := pixel.NewSprite(icon, pixel.R(0, 0, icon.Bounds().W(), icon.Bounds().H()))
	ui := Ui{
		viewport: viewport,
		ch:       ch,
		heart:    heart,
	}

	return &ui
}

func (ui *Ui) Draw(t pixel.Target) {
	marginy := 16.0
	marginx := 20.0
	vec := pixel.V(ui.viewport.Min.X+marginy, ui.viewport.Max.Y-marginy)
	ui.ch.GetPortrait().Draw(t, pixel.IM.Moved(vec))

	for i := 0; i < ui.ch.GetHp(); i++ {
		vec = vec.Add(pixel.Vec{marginx, 0})
		ui.heart.Draw(t, pixel.IM.Moved(vec))
	}
}
