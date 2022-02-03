package ui

import (
	"github.com/faiface/pixel"
)

type Ui struct {
	ch       Characterer
	viewport pixel.Rect
	heart    *pixel.Sprite
}

func New(ch Characterer, viewport pixel.Rect, path string) *Ui {
	icon, err := loadPicture(path)
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

func (ui *Ui) Draw(t pixel.Target, pos pixel.Vec, cam pixel.Vec) {
	for i := 0; i < ui.ch.GetHp(); i++ {
		vec := pixel.V(ui.viewport.Min.X+float64(i+1)*32.0, ui.viewport.Max.Y-32)
		// draw heart
		ui.heart.Draw(t, pixel.IM.Moved(vec.Sub(cam)))
	}
}
