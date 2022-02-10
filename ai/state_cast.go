package ai

import (
	"math"
	"platformer/events"

	"github.com/faiface/pixel"
)

type StateCast struct {
	id          int
	w           Worlder
	ai          *Ai
	timer       float64
	nonseelimit float64
	lastpos     pixel.Vec
	vec         pixel.Vec
}

func NewCast(ai *Ai, w Worlder) *StateCast {
	return &StateCast{
		id:          CAST,
		ai:          ai,
		w:           w,
		nonseelimit: 1,
	}
}

func (s *StateCast) Update(dt float64) {
	hero := s.w.GetHero()
	herohp := hero.GetHp()
	if herohp <= 0 {
		s.ai.SetState(IDLE, s.lastpos)
		return
	}
	heropos := hero.GetPos()
	pos := s.ai.obj.GetPos()
	dir := s.ai.obj.GetDir()
	var isSee bool
	if (heropos.X < pos.X && dir < 0) || (heropos.X > pos.X && dir > 0) {
		isSee = s.w.IsSee(pos, heropos)
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
	l := pixel.L(pos, heropos)
	if l.Len() < s.ai.obj.GetAttackrange()/2 {
		// step out of hero
		if heropos.X < pos.X {
			s.vec = pixel.Vec{1, 0}
		} else {
			s.vec = pixel.Vec{-1, 0}
		}
	}

	if math.Abs(s.lastpos.X-pos.X) < s.ai.obj.GetAttackrange() && isSee {
		m := s.ai.obj.GetMagic()
		m.SetSpell("basic")
		m.SetTarget(s.w.GetHero())

		s.ai.obj.Notify(events.CAST, pixel.ZV)
	} else {
		s.ai.obj.Notify(events.WALK, s.vec)
	}
}

func (s *StateCast) Start(poi pixel.Vec) {
	s.lastpos = poi
}

func (s *StateCast) Notify(e int, v pixel.Vec) {

}

func (s *StateCast) IsAlerted() bool {
	return true
}
