package world

import (
	"platformer/actor"

	"github.com/faiface/pixel"
)

type Speller interface {
	Update(dt float64)
	Draw(t pixel.Target)
}

var enspells []Speller
var plspells []Speller

func init() {
	enspells = make([]Speller, 0)
	plspells = make([]Speller, 0)
}

func AddSpell(owner *actor.Actor, t pixel.Vec, spell string) {

}
