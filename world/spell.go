package world

import (
	"fmt"
	"platformer/actor"
	"platformer/ai"
	"platformer/magic"

	"github.com/faiface/pixel"
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
	fmt.Println("add spell!")
	// add spell here
	//	s string, source, target pixel.Vec
	sp := magic.Create(spell, owner.GetPos(), t)

	if ai.GetByObj(owner) != nil {
		enspells = append(enspells, sp)
	} else {
		plspells = append(plspells, sp)
	}
}

/*
func (m *Magic) Update(dt float64) {
	i := 0
	for _, s := range m.spells {
		s.Update(dt)
		if !s.IsFinished() {
			m.spells[i] = s
			i++
		}
	}

	m.spells = m.spells[:i]
}
*/

func updateSpells(dt float64, enemies []*actor.Actor, player *actor.Actor) {
	updatePlSpells(dt, enemies)
	updateEnSpells(dt, []*actor.Actor{player})
}

func updatePlSpells(dt float64, hittable []*actor.Actor) {
}

func updateEnSpells(dt float64, hittable []*actor.Actor) {
}

func drawSpells(t pixel.Target) {
	for _, s := range enspells {
		s.Draw(t)
	}
	for _, s := range plspells {
		s.Draw(t)
	}
}
