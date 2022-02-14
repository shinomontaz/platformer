package ai

import (
	"platformer/events"

	"github.com/faiface/pixel"
)

type StateInvestigate struct {
	id      int
	w       Worlder
	ai      *Ai
	target  pixel.Vec
	timer   float64
	timeout float64
}

func NewInvestigate(ai *Ai, w Worlder) *StateInvestigate {
	return &StateInvestigate{
		id:      INVESTIGATE,
		ai:      ai,
		w:       w,
		timeout: 5,
	}
}

func (s *StateInvestigate) Update(dt float64) {
	//	fmt.Println("*StateInvestigate Update")
	pos := s.ai.obj.GetPos()

	hero := s.w.GetHero()
	// look for hero
	herohp := hero.GetHp()
	heropos := hero.GetPos()
	dir := s.ai.obj.GetDir()
	if (heropos.X < pos.X && dir < 0) || (heropos.X > pos.X && dir > 0) {
		// check if we see hero
		if s.w.IsSee(pos, heropos) && herohp > 0 {
			s.ai.SetState(CHOOSEATTACK, heropos)
			//			s.ai.SetState(ATTACK, heropos)
		}
		return
	}

	vec := pixel.Vec{-1, 0}
	if s.target.X > pos.X {
		vec = pixel.Vec{1, 0}
	}

	if s.target.X-pos.X < 5 {
		s.timer += dt
	}
	if s.timer > s.timeout {
		s.ai.SetState(IDLE, pixel.ZV)
	} else {
		s.ai.obj.Notify(events.WALK, vec)
	}
}

func (s *StateInvestigate) Notify(e int, v pixel.Vec) {
}

func (s *StateInvestigate) Start(poi pixel.Vec) {
	s.target = poi
	s.timer = 0
}

func (s *StateInvestigate) IsAlerted() bool {
	return false
}
