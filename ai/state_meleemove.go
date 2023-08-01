package ai

import (
	"fmt"
	"platformer/creatures"
	"platformer/events"

	"github.com/shinomontaz/pixel"
)

type StateMeleeMove struct {
	id          int
	w           Worlder
	ai          *Ai
	timer       float64
	nonseelimit float64
	timelimit   float64
	lastpos     pixel.Vec
	vec         pixel.Vec
	isbusy      bool
}

func NewMeleeMove(ai *Ai, w Worlder) *StateMeleeMove {
	return &StateMeleeMove{
		id:          MELEEMOVE,
		ai:          ai,
		w:           w,
		nonseelimit: 1,
		timelimit:   3,
	}
}

func (s *StateMeleeMove) Update(dt float64) {
	if s.isbusy {
		return
	}

	if s.ai.attackskill == nil {
		s.ai.SetState(CHOOSEATTACK, s.lastpos)
		return
	}

	hero := creatures.GetHero()
	herohp := hero.GetHp()
	if herohp <= 0 {
		s.ai.SetState(IDLE, s.lastpos)
		return
	}

	if s.timer > s.timelimit {
		s.ai.SetState(INVESTIGATE, s.lastpos)
		return
	}

	// we already check that we see target and all distances are ok
	s.ai.obj.SetSkill(s.ai.attackskill)
	//	s.ai.obj.Listen(events.WALK, s.vec)
	//	fmt.Println("meleemove")
	// for _, k := range s.ai.attackskill.Keys {
	// 	s.ai.obj.KeyEvent(k)
	// }
}

func (s *StateMeleeMove) Start(poi pixel.Vec) {
	fmt.Println("state meleemove start")
	s.isbusy = true
	s.lastpos = poi
	s.timer = 0

	pos := s.ai.obj.GetPos()
	s.vec = pixel.Vec{-1, 0}
	if s.lastpos.X > pos.X {
		s.vec = pixel.Vec{1, 0}
	}
}

func (s *StateMeleeMove) Listen(e int, v pixel.Vec) {
	if e == events.BUSY {
		s.isbusy = true
	}
	if e == events.RELEASED {
		s.isbusy = false
	}
}

func (s *StateMeleeMove) IsAlerted() bool {
	return true
}
