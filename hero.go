package main

import "github.com/faiface/pixel"

const (
	WALKING = iota
	JUMPING
	STANDING
	FIRING
	DYING
	DEAD
)

type Hero struct {
	phys  *phys
	sheet pixel.Picture
	anims map[string][]pixel.Rect

	state   string
	counter float64

	frame pixel.Rect

	sprite *pixel.Sprite
}

func (h *Hero) draw(t pixel.Target) {
	if h.sprite == nil {
		h.sprite = pixel.NewSprite(nil, pixel.Rect{})
	}
	// draw the correct frame with the correct position and direction
	ga.sprite.Set(ga.sheet, ga.frame)
	ga.sprite.Draw(t, pixel.IM.
		ScaledXY(pixel.ZV, pixel.V(
			phys.rect.W()/ga.sprite.Frame().W(),
			phys.rect.H()/ga.sprite.Frame().H(),
		)).
		ScaledXY(pixel.ZV, pixel.V(-ga.dir, 1)).
		Moved(phys.rect.Center()),
	)
}
