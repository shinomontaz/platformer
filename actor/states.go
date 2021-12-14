package actor

import (
	"platformer/controller"

	"github.com/faiface/pixel"
)

const (
	STATE_FREE = iota
	STATE_ATTACK
	STATE_HIT
	STATE_DEAD
)

type CommonState struct {
	id int
	a  *Actor
	//	animations map[string]Animation
	anims Animater
}

type FreeState struct {
	CommonState
	shift     bool
	idleLimit float64
	counter   float64
}

func (s *CommonState) GetId() int {
	return s.id
}

func NewFreeState(a *Actor, an Animater) *FreeState {
	fs := &FreeState{
		CommonState: CommonState{
			id:    STATE_FREE,
			a:     a,
			anims: an,
		},
		idleLimit: 5.0, // seconds before idle
	}

	return fs
}

func (s *FreeState) Start() {
	s.counter = 0
}

func (s *FreeState) Update(dt float64) {
	// detect animation and make it idle in case of long inactivity
	// make movements
	//	s.state.Update(dt)

	// case STANDING:
	// 	h.frame = h.anims["stand"].frames[0]
	// 	h.sheet = h.anims["stand"].sheet

}

func (s *FreeState) Notify(e int, v *pixel.Vec) {
	// if event = attack or hit => switch state
	// if event = move => make it move
	if e == controller.E_MOVE {
		// get vector
		// check physics
		//		s.state.Notify(controller.E_MOVE)
	}
	if e == controller.E_CTRL {
		// switch state to attacking
		// if s.state.GetId() != ASTATE_JUMP {
		// 	s.a.SetState(STATE_ATTACK)
		// }
	}
	if e == controller.E_SHIFT {
		// to free state add modifier
		//		s.state.Notify(controller.E_SHIFT)
	}

	// set animation?
}

func (s *FreeState) GetSprite() *pixel.Sprite {
	//	fmt.Println("free state get sprite")
	return s.anims.GetSprite("idle", 0)
}

type AttackState struct {
	CommonState
	time      float64
	idleLimit float64
}

func NewAttackState(a *Actor, an Animater) *AttackState {
	fs := &AttackState{
		CommonState: CommonState{
			id:    STATE_ATTACK,
			a:     a,
			anims: an,
		},
		idleLimit: 5.0, // seconds before idle
	}

	return fs
}

func (s *AttackState) Start() {
	s.time = 0.0
}

func (s *AttackState) Notify(e int, v *pixel.Vec) {
	// here we don't care of any controller event
}

func (s *AttackState) Update(dt float64) {
	if s.time > s.idleLimit {
		s.a.SetState(STATE_FREE)
		return
	}

	s.time += dt
	//	s.pc.SetVec(pixel.ZV)
	//	s.pc.SetCmd(STRIKE)
}

func (s *AttackState) GetSprite() *pixel.Sprite {
	//	return s.animations["attack1"].GetSprite(0)
	return s.anims.GetSprite("attack1", 0)
}

type DeadState struct {
	CommonState
}

func NewDeadState(a *Actor, an Animater) *DeadState {
	fs := &DeadState{
		CommonState: CommonState{
			id:    STATE_DEAD,
			a:     a,
			anims: an,
		},
	}

	return fs
}

func (s *DeadState) Start() {
}

func (s *DeadState) Notify(e int, v *pixel.Vec) {
	// here we don't care of any controller event
}

func (s *DeadState) Update(dt float64) {
}

func (s *DeadState) GetSprite() *pixel.Sprite {
	return s.anims.GetSprite("dead", 0)
}

type HitState struct {
	CommonState
	time      float64
	timelimit float64
}

func NewHitState(a *Actor, an Animater) *HitState {
	fs := &HitState{
		CommonState: CommonState{
			id:    STATE_HIT,
			a:     a,
			anims: an,
		},
	}

	return fs
}

func (s *HitState) Start() {
}

func (s *HitState) Notify(e int, v *pixel.Vec) {
	// here we don't care of any controller event
}

func (s *HitState) Update(dt float64) {
	if s.time > s.timelimit {
		s.a.SetState(STATE_FREE)
		return
	}

	s.time += dt
}

func (s *HitState) GetSprite() *pixel.Sprite {
	return s.anims.GetSprite("hurt", 0)
}
