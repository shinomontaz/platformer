package ai

import (
	"math"
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
				s.ai.SetState(INVESTIGATE, s.lastpos)
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

	if math.Abs(s.lastpos.X-pos.X) < 20 {
		s.vec = pixel.ZV
		s.ai.obj.Notify(events.CTRL, s.vec)
	} else {
		s.ai.obj.Notify(events.WALK, s.vec)
	}
}

func (s *StateAttack) Start(poi pixel.Vec) {
	s.lastpos = poi
}

func (s *StateAttack) Notify(e int, v pixel.Vec) {

}

func (s *StateAttack) IsAlerted() bool {
	return true
}
