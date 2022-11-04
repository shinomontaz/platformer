package magic

import (
	"platformer/actor"
	"platformer/common"

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

func Create(s string, owner *actor.Actor, target pixel.Vec, objs []common.Objecter) Speller {
	switch s {
	case "deathstrike":
		return NewDeathstrike(s, owner, target, objs)
	default:
		return nil
	}
}
