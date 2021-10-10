package actor

import (
	"platformer/controller"
)

const (
	STATE_FREE = iota
	STATE_ATTACK
	STATE_HIT
	STATE_DEAD
)

type FreeState struct {
	id int
	a  *Actor
}

func (s *FreeState) Start() {
}

func (s *FreeState) Update(dt float64) {
	// detect animation and make it idle in case of long inactivity
	// make movements
}

func (s *FreeState) Notify(e int) {
	// if event = attack or hit => switch state
	// if event = move => make it move
	if e == controller.E_MOVE {
		// get vector
		// check physics
	}
	if e == controller.E_CTRL {
		// switch state to attacking
		s.a.SetState(STATE_ATTACK)
	}
	if e == controller.E_SHIFT {
		// to free state add modifier
	}

	// set animation?
}

type AttackState struct {
	id        int
	a         *Actor
	time      float64
	timelimit float64
}

func (s *AttackState) Start() {
	s.time = 0.0
}

func (s *AttackState) Update(dt float64) {
	if s.time > s.timelimit {
		s.a.SetState(STATE_FREE)
		return
	}

	s.time += dt
	//	s.pc.SetVec(pixel.ZV)
	//	s.pc.SetCmd(STRIKE)
}

type DeadState struct {
	id int
	a  *Actor
}

func (s *DeadState) Start() {
}

func (s *DeadState) Update(dt float64) {
}

type HitState struct {
	id        int
	a         *Actor
	time      float64
	timelimit float64
}

func (s *HitState) Start() {
}

func (s *HitState) Update(dt float64) {
	if s.time > s.timelimit {
		s.a.SetState(STATE_FREE)
		return
	}

	s.time += dt
}
