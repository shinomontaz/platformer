package ai

import (
	"platformer/events"

	"github.com/shinomontaz/pixel"
)

type StateIdle struct {
	id     int
	w      Worlder
	ai     *Ai
	isbusy bool
}

func NewIdle(ai *Ai, w Worlder) *StateIdle {
	return &StateIdle{
		id: IDLE,
		ai: ai,
		w:  w,
	}
}

func (s *StateIdle) Update(dt float64) {
	if s.isbusy {
		return
	}

	hero := s.w.GetHero()
	herohp := hero.GetHp()
	heropos := hero.GetPos()
	pos := s.ai.obj.GetPos()
	dir := s.ai.obj.GetDir()
	if (heropos.X < pos.X && dir < 0) || (heropos.X > pos.X && dir > 0) {
		// check if we see target
		if s.w.IsSee(pos, heropos) && herohp > 0 {
			s.w.AddAlert(pos, 200)
			s.ai.SetState(ATTACK, heropos)
		}
	}
}

func (s *StateIdle) Start(poi pixel.Vec) {

}

func (s *StateIdle) Listen(e int, v pixel.Vec) {
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

func (s *StateIdle) GetVec() pixel.Vec {
	return pixel.ZV
}

func (s *StateIdle) IsAlerted() bool {
	return false
}
