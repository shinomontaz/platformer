package ai

import (
	"platformer/events"

	"github.com/faiface/pixel"
)

type StateAttack struct {
	id          int
	w           Worlder
	ai          *Ai
	timer       float64
	nonseelimit float64
	lastpos     pixel.Vec
	vec         pixel.Vec
}

func NewAttack(ai *Ai, w Worlder) *StateAttack {
	return &StateAttack{
		id:          ATTACK,
		ai:          ai,
		w:           w,
		nonseelimit: 1,
	}
}

func (s *StateAttack) Update(dt float64) {
	heropos := s.w.GetHero()
	pos := s.ai.obj.GetPos()
	dir := s.ai.obj.GetDir()
	if (heropos.X < pos.X && dir < 0) || (heropos.X > pos.X && dir > 0) {
		isSee := s.w.IsSee(pos, heropos)
		if !isSee {
			s.timer += dt
			if s.timer > s.nonseelimit {
				s.ai.SetState(IDLE)
			}
		} else {
			s.lastpos = heropos
			s.timer = 0
		}
	}

	s.vec = pixel.Vec{-1, 0}
	if s.lastpos.X > pos.X {
		s.vec = pixel.Vec{1, 0}
	}

	s.ai.Notify(events.WALK)

}

func (s *StateAttack) Start() {

}

func (s *StateAttack) GetVec() pixel.Vec {
	return s.vec
}
