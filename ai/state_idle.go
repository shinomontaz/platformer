package ai

import (
	"platformer/events"

	"github.com/faiface/pixel"
)

type StateIdle struct {
	id int
	w  Worlder
	ai *Ai
}

func NewIdle(ai *Ai, w Worlder) *StateIdle {
	return &StateIdle{
		id: IDLE,
		ai: ai,
		w:  w,
	}
}

func (s *StateIdle) Update(dt float64) {
	hero := s.w.GetHero()
	pos := s.ai.obj.GetPos()
	dir := s.ai.obj.GetDir()
	if (hero.X < pos.X && dir < 0) || (hero.X > pos.X && dir > 0) {
		// check if we see target
		if s.w.IsSee(pos, hero) {
			s.w.AddAlert(pos, 100)
			s.ai.SetState(ATTACK, hero)
		}
	}
}

func (s *StateIdle) Start(poi pixel.Vec) {

}

func (s *StateIdle) Notify(e int, v pixel.Vec) {
	if e == events.ALERT {
		s.ai.SetState(INVESTIGATE, v)
	}
}

func (s *StateIdle) GetVec() pixel.Vec {
	return pixel.ZV
}

func (s *StateIdle) IsAlerted() bool {
	return false
}
