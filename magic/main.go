package magic

import (
	"platformer/actor"

	"github.com/faiface/pixel"
)

var w Worlder

func SetWorld(ww Worlder) {
	w = ww
}

func Create(s string, owner *actor.Actor, target pixel.Vec) Speller {
	switch s {
	case "deathstrike":
		return NewDeathstrike(s, owner, target)
	default:
		return nil
	}
}
