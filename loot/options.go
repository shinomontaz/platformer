package loot

import (
	"platformer/config"

	"github.com/shinomontaz/pixel"
)

type Option func(*Loot)

func WithAnimDir(animdir float64) Option {
	return func(a *Loot) {
		a.animdir = animdir
	}
}

func WithGravity(grav float64) Option {
	return func(a *Loot) {
		a.grav = grav
	}
}

func WithSound(seffects []config.Soundeffect) Option {
	return func(a *Loot) {
		for _, se := range seffects {
			a.sounds[se.Type] = soundeffect{
				List: se.List,
			}
		}
	}
}

func WithVelocity(v pixel.Vec) Option {
	return func(a *Loot) {
		a.vel = v
	}
}

func WithPortrait(path string) Option {
	return func(a *Loot) {
		a.portrait = path
	}
}
