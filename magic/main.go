package magic

import (
	"github.com/faiface/pixel"
)

var w Worlder

func SetWorld(ww Worlder) {
	w = ww
}

func Create(s string, source, target pixel.Vec) Speller {
	switch s {
	case "deathstrike":
		return NewDeathstrike(s, source, target)
	default:
		return nil
	}
}
