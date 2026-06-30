package ai

import (
	"math"
	"platformer/actor"
	"platformer/bindings"
	"platformer/events"
	"sort"

	"github.com/shinomontaz/pixel"
)

type StateInvestigate struct {
	id      int
	w       Worlder
	ai      *Ai
	target  pixel.Vec
	timer   float64
	timeout float64
	isbusy  bool
	skills  []*actor.Skill
	skills2 []*actor.Skill
}

func NewInvestigate(ai *Ai, w Worlder) *StateInvestigate {
	in := &StateInvestigate{
		id:      INVESTIGATE,
		ai:      ai,
		w:       w,
		timeout: 5,
	}

	in.skills = ai.obj.GetSkills()
	sort.Slice(in.skills, func(i, j int) bool {
		return in.skills[i].Max < in.skills[j].Max
	})
	in.skills2 = ai.obj.GetSkills()
	sort.Slice(in.skills2, func(i, j int) bool {
		return in.skills2[i].Min < in.skills2[j].Min
	})

	return in
}

func (s *StateInvestigate) Update(dt float64) {
	if s.isbusy {
		return
	}
	pos := s.ai.obj.GetPos()

	hero := s.ai.obj.GetEnemy()
	if hero == nil {
		return
	}
	herohp := hero.GetHp()
	heropos := hero.GetPos()
	dir := s.ai.obj.GetDir()
	l := pixel.L(pos, heropos)
	currDist := l.Len()

	if ((heropos.X < pos.X && dir < 0) || (heropos.X > pos.X && dir > 0)) && currDist < s.skills[len(s.skills)-1].Max {
		if s.w.IsSee(pos, heropos) && herohp > 0 {
			s.ai.SetState(CHOOSEATTACK, heropos)
			return
		}
	}

	// Дошли до цели — стоим и ждём, не ходим (иначе осцилляция вокруг target)
	if math.Abs(s.target.X-pos.X) < 10 {
		s.timer += dt
	} else {
		// Идём к цели с проверкой безопасности
		moveDir := 1.0
		if s.target.X < pos.X {
			moveDir = -1.0
		}

		v := pixel.V(moveDir, 0)
		groundrate := s.ai.obj.StepPrediction(events.WALK, v)

		if groundrate > 0.5 {
			b := bindings.Active.GetBinding(bindings.KeyAction["right"])
			if moveDir < 0 {
				b = bindings.Active.GetBinding(bindings.KeyAction["left"])
			}
			s.ai.obj.KeyAction(b)
		}
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
