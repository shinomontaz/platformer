package magic

import (
	"platformer/actor"

	"github.com/shinomontaz/pixel"
)

var w Worlder

func SetWorld(ww Worlder) {
	w = ww
}

type SpellProf struct {
	r     pixel.Rect
	sound string
}

var profiles map[string]SpellProf = make(map[string]SpellProf)

func Load(name string, cfg Spellprofile) {
	profiles[name] = SpellProf{
		r:     cfg.GetHitbox(),
		sound: cfg.GetSound(),
	}
}

func Create(s string, owner *actor.Actor, target pixel.Vec) Speller {
	switch s {
	case "deathstrike":
		return NewDeathstrike(s, owner, target)
	default:
		return nil
	}
}
