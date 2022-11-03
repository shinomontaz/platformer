package ai

import (
	"platformer/events"

	"github.com/shinomontaz/pixel"
)

type StateIdle2 struct {
	id     int
	w      Worlder
	ai     *Ai
	isbusy bool
}

func NewIdle2(ai *Ai, w Worlder) *StateIdle2 {
	return &StateIdle2{
		id: IDLE,
		ai: ai,
		w:  w,
	}
}

func (s *StateIdle2) Update(dt float64) {
	if s.isbusy {
		return
	}
}

func (s *StateIdle2) Start(poi pixel.Vec) {

}

func (s *StateIdle2) Listen(e int, v pixel.Vec) {
	if e == events.ALERT {
		s.ai.SetState(INVESTIGATE, v)
	}
	if e == events.BUSY {
		s.isbusy = true
	}
	if e == events.RELEASED {
		s.isbusy = false
	}
}

func (s *StateIdle2) GetVec() pixel.Vec {
	return pixel.ZV
}

func (s *StateIdle2) IsAlerted() bool {
	return false
}
