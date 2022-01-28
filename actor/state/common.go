package state

import "platformer/actor/statemachine"

type Common struct {
	id    int
	a     Actor
	anims Animater
	trs   statemachine.Transition
}

func (s *Common) GetId() int {
	return s.id
}

func (s *Common) checkTransitions(e int) {
	if st, ok := s.trs.List[e]; ok {
		s.a.SetState(st)
	}
}
