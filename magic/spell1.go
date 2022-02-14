package magic

import (
	"github.com/faiface/pixel"
)

type Deathstrike struct {
	anim         Animater
	sprite       *pixel.Sprite
	vel          pixel.Vec
	w            Worlder
	power        int
	source       pixel.Vec
	target       pixel.Vec
	ttl          float64
	timer        float64
	prepare_time float64
}

func (d *Deathstrike) Draw(t pixel.Target) {

}

func (d *Deathstrike) Update(dt float64) {
	d.timer += dt
}

// func (d *Deathstrike) Apply(a Actor) {

// }

func (d *Deathstrike) IsFinished() bool {
	if d.timer > d.ttl {
		return true
	}
	return false
}
