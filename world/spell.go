package world

import (
	"platformer/actor"
	"platformer/ai"
	"platformer/magic"

	"github.com/shinomontaz/pixel"
)

type Speller interface {
	Update(dt float64)
	Draw(t pixel.Target)
	IsFinished() bool
}

var enspells []Speller
var plspells []Speller

func init() {
	enspells = make([]Speller, 0)
	plspells = make([]Speller, 0)
}

func AddSpell(owner *actor.Actor, t pixel.Vec, spell string) {
	// add spell here
	//	s string, source, target pixel.Vec
	sp := magic.Create(spell, owner, t)

	if ai.GetByObj(owner) != nil {
		enspells = append(enspells, sp)
	} else {
		plspells = append(plspells, sp)
	}
}

func updateSpells(dt float64, enemies []*actor.Actor, player *actor.Actor) {
	updatePlSpells(dt, enemies)
	updateEnSpells(dt, []*actor.Actor{player})
}

func updatePlSpells(dt float64, hittable []*actor.Actor) {
	i := 0
	for _, s := range plspells {
		// for _, hh := range hittable {
		// 	if hh == s.GetOwner() {
		// 		continue
		// 	}
		// if s.GetHitted(hh) {
		// 	continue
		// }

		//			r := hh.GetRect()
		// if s.GetRect().Intersects(r) {
		// 	vec := pixel.ZV // TODO: detect hit vector
		// 	vec.X = -1
		// 	if r.Center().X > b.dir.A.X {
		// 		vec.X = 1
		// 	}
		// 	hh.Hit(vec, s.GetPower())
		// 	b.hitted[hh] = struct{}{}
		// }
		//		}
		s.Update(dt)
		if !s.IsFinished() {
			plspells[i] = s
			i++
		}
	}

	plspells = plspells[:i]
}

func updateEnSpells(dt float64, hittable []*actor.Actor) {
	i := 0
	for _, s := range enspells {
		s.Update(dt)
		if !s.IsFinished() {
			enspells[i] = s
			i++
		}
	}

	enspells = enspells[:i]
}

func drawSpells(t pixel.Target) {
	for _, s := range enspells {
		s.Draw(t)
	}
	for _, s := range plspells {
		s.Draw(t)
	}
}
