package ai

import (
	"math/rand"
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
			return
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
	currDist := l.Len()
	if currDist < s.ai.attackskill.Min || currDist > s.ai.attackskill.Max {
		// make decision - to step out or choose another attack skill
		dice := rand.Float64()
		if dice > 0.5 {
			s.ai.SetState(CHOOSEATTACK, s.lastpos)
			return
		}

		if heropos.X < pos.X {
			s.vec = pixel.Vec{-1, 0}
		} else {
			s.vec = pixel.Vec{1, 0}
		}
		s.ai.obj.Listen(events.WALK, s.vec)
		return
	}

	// we already check that we see target and all distances are ok
	s.vec = pixel.ZV
	s.ai.obj.SetTarget(heropos)
	s.ai.obj.SetSkill(s.ai.attackskill)
	//	fmt.Println("attack state notify: ", s.ai.attackskill.Event, s.vec)
	s.ai.obj.Listen(s.ai.attackskill.Event, s.vec)
}

func (s *StateAttack) Start(poi pixel.Vec) {
	s.lastpos = poi
}

func (s *StateAttack) Listen(e int, v pixel.Vec) {
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
