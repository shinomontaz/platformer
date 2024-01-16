package state

import (
	"platformer/actor/statemachine"

	"github.com/shinomontaz/pixel"
)

type Common struct {
	id               int
	a                Actor
	anims            Animater
	trs              statemachine.Transition
	busy             bool
	iswater          bool
	iswaterresistant bool
}

func (s *Common) GetId() int {
	return s.id
}

func (s *Common) checkTransitions(e int, v *pixel.Vec) {
	if st, ok := s.trs.List[e]; ok {
		s.a.SetState(st)
	}
}

func (s *Common) Busy() bool {
	return s.busy
}

func (s *Common) SetWater(b bool) {
	s.iswater = b
}

func (s *Common) SetWaterResistant(b bool) {
	s.iswaterresistant = b
}
