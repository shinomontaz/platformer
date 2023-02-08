package loot

import (
	"platformer/config"
)

type Option func(*Loot)

func WithAnimDir(animdir float64) Option {
	return func(a *Loot) {
		a.animdir = animdir
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
