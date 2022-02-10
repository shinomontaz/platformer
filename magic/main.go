package magic

import (
	"github.com/faiface/pixel"
)

type Magic struct {
	//	anim   Animater
	source Actor
	target Actor
}

func (m *Magic) Cast() {

}

func (m *Magic) SetSpell(s string) {

}

func (m *Magic) SetTarget(t Actor) {
	m.target = t
}

func (m *Magic) Draw(t pixel.Target) {

}
