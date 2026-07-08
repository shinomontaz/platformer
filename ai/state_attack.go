package ai

import (
	"platformer/bindings"

	"platformer/events"

	"github.com/shinomontaz/pixel"
)

type StateAttack struct {
	id            AiState
	w             Worlder
	ai            *Ai
	timer         float64
	nonseelimit   float64
	lastpos       pixel.Vec
	vec           pixel.Vec
	isbusy        bool
	cooldownTimer float64
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
	if s.isbusy {
		return
	}

	if s.cooldownTimer > 0 {
		s.cooldownTimer -= dt
		return
	}

	if s.ai.attackskill == nil {
		s.ai.SetState(CHOOSEATTACK, s.lastpos)
		return
	}

	hero := s.ai.obj.GetEnemy()
	if hero == nil {
		return
	}

	herohp := hero.GetHp()
	if herohp <= 0 {
		s.ai.SetState(s.ai.homeState, s.lastpos)
		return
	}
	heropos := hero.GetPos()
	s.lastpos = heropos

	pos := s.ai.obj.GetPos()
	dir := s.ai.obj.GetDir()
	var isSee bool

	// Проверяем видимость только если смотрим в сторону героя
	if (heropos.X < pos.X && dir < 0) || (heropos.X > pos.X && dir > 0) {
		isSee = s.w.IsSee(pos, heropos)
		if !isSee {
			s.timer += dt
			if s.timer > s.nonseelimit {
				s.ai.SetState(INVESTIGATE, s.lastpos)
			}

			return
		}
	} else { // мы не смотрим в сторону героя
		s.timer = 0
		s.ai.SetState(INVESTIGATE, s.lastpos)
		return
	}

	s.vec = pixel.V(-1, 0)
	if s.lastpos.X > pos.X {
		s.vec = pixel.V(1, 0)
	}
	l := pixel.L(pos, heropos)
	currDist := l.Len()
	if currDist < s.ai.attackskill.Min {
		s.ai.SetState(CHOOSEATTACK, s.lastpos)

		return
	}

	if currDist > s.ai.attackskill.Max {
		s.ai.SetState(INVESTIGATE, s.lastpos)

		return
	}

	// we already check that we see target and all distances are ok
	s.ai.obj.SetTarget(s.lastpos)
	s.ai.obj.SetSkill(s.ai.attackskill)
	for _, k := range s.ai.attackskill.Keys {
		b := bindings.Active.GetBinding(k)
		s.ai.obj.KeyAction(b)
	}

	s.cooldownTimer = s.ai.attackskill.Cooldown
}

func (s *StateAttack) Start(poi pixel.Vec) {
	s.lastpos = poi
	s.timer = 0
	s.cooldownTimer = 0
}

func (s *StateAttack) EventAction(e int) {
	if e == events.BUSY {
		s.isbusy = true
	}
	if e == events.RELEASED {
		s.isbusy = false
	}
}

func (s *StateAttack) IsAlerted() bool {
	return true
}
