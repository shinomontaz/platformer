package ai

import (
	"platformer/bindings"
	"platformer/common"
	"platformer/events"

	"github.com/shinomontaz/pixel"
)

type StateAttack struct {
	id          int
	w           Worlder
	ai          *Ai
	timer       float64
	nonseelimit float64
	counter     int
	lastpos     pixel.Vec
	vec         pixel.Vec
	isbusy      bool
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
		s.ai.SetState(IDLE, s.lastpos)
		return
	}
	heropos := hero.GetPos()
	s.lastpos = heropos

	if s.counter > 0 { // here we made decision to switch to state BUSTLE with some probability - after strike
		coeff := 0.25
		if s.ai.attackskill.Type == "melee" {
			coeff = 0.5
		} else {
			coeff = 0
		}
		dice := float64(s.counter) * common.GetRandFloat()
		if dice > coeff {
			s.ai.SetState(BUSTLE, s.lastpos)
			s.counter = 0
			return
		}
	}

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
			return
		} else {
			//			s.lastpos = heropos
			s.timer = 0
		}
	}

	s.vec = pixel.Vec{-1, 0}
	if s.lastpos.X > pos.X {
		s.vec = pixel.Vec{1, 0}
	}
	l := pixel.L(pos, heropos)
	currDist := l.Len()
	if currDist < s.ai.attackskill.Min || currDist > s.ai.attackskill.Max {
		// make decision - to step out or choose another attack skill
		dice := common.GetRandFloat()
		if dice > 0.5 {
			s.ai.SetState(CHOOSEATTACK, s.lastpos)
			return
		}
		s.ai.SetState(INVESTIGATE, s.lastpos)
		return
	}

	// we already check that we see target and all distances are ok
	s.ai.obj.SetTarget(s.lastpos)
	s.ai.obj.SetSkill(s.ai.attackskill)
	for _, k := range s.ai.attackskill.Keys {
		b := bindings.Active.GetBinding(k)
		s.ai.obj.KeyAction(b)
		//		fmt.Println("s.ai.obj.KeyAction(b)", b)
	}

	s.counter++
}

func (s *StateAttack) Start(poi pixel.Vec) {
	s.lastpos = poi
	s.timer = 0
	s.counter = 0
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
