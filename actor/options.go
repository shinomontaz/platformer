package actor

import (
	"platformer/actor/statemachine"
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

func WithStatemachine(sm statemachine.Machine) Option {
	return func(a *Actor) {
		a.sm = sm
	}
}

func WithAnimDir(animdir float64) Option {
	return func(a *Actor) {
		a.animdir = animdir
	}
}
