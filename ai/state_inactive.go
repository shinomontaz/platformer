package ai

import "github.com/shinomontaz/pixel"

type StateInactive struct {
	id     int
	w      Worlder
	ai     *Ai
	isbusy bool
}

func NewInactive(ai *Ai, w Worlder) *StateInactive {
	return &StateInactive{
		id: INACTIVE,
		ai: ai,
		w:  w,
	}
}

func (s *StateInactive) Update(dt float64) {
	// do nothing
}

func (s *StateInactive) Start(poi pixel.Vec) {

}

func (s *StateInactive) Listen(e int, v pixel.Vec) {

}

func (s *StateInactive) GetVec() pixel.Vec {
	return pixel.ZV
}

func (s *StateInactive) IsAlerted() bool {
	return false
}
