package ai

import (
	"log"
	"platformer/actor"
	"platformer/common/rand"
	"platformer/events"

	"github.com/shinomontaz/pixel"
)

type StateChooseAttack struct {
	id          AiState
	w           Worlder
	ai          *Ai
	timer       float64
	nonseelimit float64
	lastpos     pixel.Vec
	isbusy      bool
	cand        []*actor.Skill
}

func NewChooseAttack(ai *Ai, w Worlder) *StateChooseAttack {
	return &StateChooseAttack{
		id:          CHOOSEATTACK,
		ai:          ai,
		w:           w,
		nonseelimit: 1,
	}
}

func (s *StateChooseAttack) Update(dt float64) {
	if s.isbusy {
		return
	}
	t, ok := s.ai.look()
	if !ok {
		s.ai.SetState(s.ai.homeState, pixel.ZV)
		return
	}

	if !s.ai.facing(t) {
		s.ai.obj.SetDir(sign(t.dx)) // уже в бою — разворачиваемся к герою
	}
	if !s.ai.sees(t) {
		s.timer += dt
		if s.timer > s.nonseelimit {
			s.ai.SetState(INVESTIGATE, s.lastpos)
		}
		return
	}
	s.timer = 0
	s.lastpos = t.hero.GetPos()

	if sk := s.pick(t); sk != nil {
		s.ai.attackskill = sk
		s.ai.obj.SetSkill(sk)
		s.ai.SetState(ATTACK, s.lastpos)
		return
	}

	if t.gap > s.ai.maxReach {
		s.ai.SetState(INVESTIGATE, s.lastpos)
		return
	}

	// досягаемо, но ни один навык не подходит: стрелку — отойти, остальным — подойти
	dir := sign(t.dx)
	if !s.ai.hasMelee && t.gap < s.ai.minGap {
		dir = -dir
	}
	s.ai.step(dir)
}

// взвешенный случайный выбор среди подходящих навыков
func (s *StateChooseAttack) pick(t target) *actor.Skill {
	s.cand = s.cand[:0]
	total := 0
	for _, sk := range s.ai.obj.GetSkills() {
		ok := sk.Weight > 0 && s.ai.canUse(sk, t)
		if s.ai.debug && s.ai.debugTick <= 0 {
			log.Printf("[ancient] skill=%s type=%s w=%d min=%.0f max=%.0f | gap=%.1f dy=%.1f -> %v",
				sk.Name, sk.Type, sk.Weight, sk.Min, sk.Max, t.gap, t.dy, ok)
		}
		if !ok {
			continue
		}
		s.cand = append(s.cand, sk)
		total += sk.Weight
	}

	if s.ai.debug && s.ai.debugTick <= 0 {
		log.Printf("[ancient] candidates=%d", len(s.cand))
		s.ai.debugTick = 0.5
	}

	if total == 0 {
		return nil
	}
	r := rand.Float64() * float64(total)
	for _, sk := range s.cand {
		r -= float64(sk.Weight)
		if r < 0 {
			return sk
		}
	}
	return s.cand[len(s.cand)-1]
}

func (s *StateChooseAttack) Start(poi pixel.Vec) {
	s.lastpos = poi
}

func (s *StateChooseAttack) EventAction(e int) {
	if e == events.BUSY {
		s.isbusy = true
	}
	if e == events.RELEASED {
		s.isbusy = false
	}
}

func (s *StateChooseAttack) IsAlerted() bool {
	return true
}
