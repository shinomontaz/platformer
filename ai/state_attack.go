package ai

import (
	"log"
	"platformer/bindings"

	"platformer/events"

	"github.com/shinomontaz/pixel"
)

type StateAttack struct {
	id            AiState
	w             Worlder
	ai            *Ai
	lastpos       pixel.Vec
	isbusy        bool
	used          bool
	cooldownTimer float64
}

func NewAttack(ai *Ai, w Worlder) *StateAttack {
	return &StateAttack{id: ATTACK, ai: ai, w: w}
}

func (s *StateAttack) Update(dt float64) {
	if s.isbusy {
		return
	}

	if s.used { // ударили — ждём кулдаун и выбираем заново
		s.cooldownTimer -= dt
		if s.cooldownTimer <= 0 {
			s.ai.SetState(CHOOSEATTACK, s.lastpos)
		}
		return
	}

	sk := s.ai.attackskill
	if sk == nil {
		s.ai.SetState(CHOOSEATTACK, s.lastpos)

		return
	}

	t, ok := s.ai.look()
	if !ok {
		s.ai.SetState(s.ai.homeState, s.lastpos)

		return
	}
	s.lastpos = t.hero.GetPos()
	if !s.ai.canUse(sk, t) { // герой успел сместиться
		s.ai.SetState(CHOOSEATTACK, s.lastpos)
		return
	}

	s.ai.obj.SetDir(sign(t.dx))

	// стреляем только стоя: ждём, пока актёр затормозит и перейдёт в STAND
	if sk.Type == "ranged" {
		if v := s.ai.obj.GetVel(); v.X != 0 || !s.ai.obj.IsGround() {
			return
		}
	}

	s.ai.obj.SetTarget(s.lastpos)
	s.ai.obj.SetSkill(sk)

	if s.ai.debug {
		log.Printf("[ancient] ATTACK use=%s keys=%v", sk.Name, sk.Keys)
	}

	for _, k := range sk.Keys {
		s.ai.obj.KeyAction(bindings.Active.GetBinding(k))
	}
	s.used = true
	s.cooldownTimer = 0.2
	//	s.cooldownTimer = sk.Cooldown
}

func (s *StateAttack) Start(poi pixel.Vec) {
	s.lastpos = poi
	s.used = false
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
