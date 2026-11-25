package ai

import (
	"math"
	"platformer/events"

	"github.com/shinomontaz/pixel"
)

type StateInvestigate struct {
	id      AiState
	w       Worlder
	ai      *Ai
	target  pixel.Vec
	timer   float64
	timeout float64
	isbusy  bool
}

func NewInvestigate(ai *Ai, w Worlder) *StateInvestigate {
	return &StateInvestigate{id: INVESTIGATE, ai: ai, w: w, timeout: 5}
}

func (s *StateInvestigate) Update(dt float64) {
	if s.isbusy {
		return
	}
	if t, ok := s.ai.look(); ok && s.ai.facing(t) && t.gap <= s.ai.maxReach && s.ai.sees(t) {
		s.ai.SetState(CHOOSEATTACK, t.hero.GetPos())
		return
	}

	pos := s.ai.obj.GetPos()

	// Дошли до цели — стоим и ждём, не ходим (иначе осцилляция вокруг target)
	if math.Abs(s.target.X-pos.X) < 10 {
		s.timer += dt
	} else {
		moveDir := 1.0
		if s.target.X < pos.X {
			moveDir = -1.0
		}
		s.ai.step(moveDir)
	}

	if s.timer > s.timeout {
		s.ai.SetState(s.ai.homeState, pixel.ZV)
	}
}

func (s *StateInvestigate) EventAction(e int) {
	if e == events.BUSY {
		s.isbusy = true
	}
	if e == events.RELEASED {
		s.isbusy = false
	}
}

func (s *StateInvestigate) Start(poi pixel.Vec) {
	s.target = poi
	s.timer = 0
}

func (s *StateInvestigate) IsAlerted() bool {
	return false
}
