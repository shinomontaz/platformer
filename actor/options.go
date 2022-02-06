package actor

import (
	"platformer/actor/statemachine"
	"platformer/common"
	"platformer/config"

	"github.com/faiface/pixel"
)

type Option func(*Actor)

func WithRun(speed float64) Option {
	return func(a *Actor) {
		a.runspeed = speed
	}
}

func WithWalk(speed float64) Option {
	return func(a *Actor) {
		a.walkspeed = speed
	}
}

func WithJump(force float64) Option {
	return func(a *Actor) {
		a.jumpforce = force
	}
}

func WithStatemachine(sm *statemachine.Machine) Option {
	return func(a *Actor) {
		a.sm = sm
	}
}

func WithAnimDir(animdir float64) Option {
	return func(a *Actor) {
		a.animdir = animdir
	}
}
func WithHP(hp int) Option {
	return func(a *Actor) {
		a.hp = hp
	}
}

func WithStrength(strength int) Option {
	return func(a *Actor) {
		a.strength = strength
	}
}

func WithPortrait(path string) Option {
	return func(a *Actor) {
		if path == "" {
			return
		}
		prt, err := common.LoadPicture(path)
		if err != nil {
			panic(err)
		}

		a.portrait = pixel.NewSprite(prt, pixel.R(0, 0, prt.Bounds().W(), prt.Bounds().H()))
	}
}

func WithSound(seffects []config.Soundeffect) Option {
	return func(a *Actor) {
		for _, se := range seffects {
			a.sounds[se.Type] = soundeffect{
				List: se.List,
			}
		}
	}
}
