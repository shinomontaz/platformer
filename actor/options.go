package actor

import (
	"platformer/actor/statemachine"
	"platformer/common"

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
