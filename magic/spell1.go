package magic

import (
	"math"
	"platformer/animation"

	"github.com/faiface/pixel"
)

type Deathstrike struct {
	anim          Animater
	sprite        *pixel.Sprite
	animSpriteNum int
	vel           pixel.Vec
	w             Worlder
	power         int
	source        pixel.Vec
	target        pixel.Vec
	ttl           float64
	timer         float64
	prepare_time  float64
	rect          pixel.Rect
}

func NewDeathstrike(s string, source, target pixel.Vec) *Deathstrike {
	return &Deathstrike{
		anim:         animation.Get(s),
		w:            w,
		power:        1,
		source:       source,
		target:       target,
		ttl:          3.0,
		prepare_time: 1.0,
		rect:         pixel.R(0, 0, 48, 64),
	}
}

func (d *Deathstrike) Draw(t pixel.Target) {
	if d.sprite == nil {
		d.sprite = pixel.NewSprite(nil, pixel.Rect{})
	}
	pic, rect := d.anim.GetSprite("spell", d.animSpriteNum)
	d.sprite.Set(pic, rect)

	drawrect := d.rect.ResizedMin(pixel.Vec{d.rect.W(), d.rect.H()})
	d.sprite.Draw(t, pixel.IM.
		ScaledXY(pixel.ZV, pixel.V(
			drawrect.W()/d.sprite.Frame().W(),
			drawrect.H()/d.sprite.Frame().H(),
		)).Moved(drawrect.Center()),
	)
}

func (d *Deathstrike) Update(dt float64) {
	d.timer += dt
	d.animSpriteNum = int(math.Floor(d.timer / 0.1))
}

// func (d *Deathstrike) Apply(a Actor) {

// }

func (d *Deathstrike) IsFinished() bool {
	return d.timer > d.ttl
}
