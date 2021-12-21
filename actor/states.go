package actor

import (
	"math"
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
	idleLimit     float64
	counter       float64
	sprite        *pixel.Sprite
	isJumping     bool
	isShift       bool
	animName      string
	animSpriteNum int
	newStateAnim  int
	stateAnim     int
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
		stateAnim:    STANDING,
		newStateAnim: STANDING,
		idleLimit:    5.0, // seconds before idle
	}

	return fs
}

func (s *FreeState) resetAnim() {
	s.counter = 0
}

func (s *FreeState) Start() {
	s.resetAnim()
}

func (s *FreeState) Update(dt float64) {
	s.counter += dt

	if s.counter > s.idleLimit && s.stateAnim == STANDING && (s.newStateAnim == s.stateAnim || s.newStateAnim == STANDING) { // make idle animation
		s.newStateAnim = IDLE
	}

	if s.stateAnim == IDLE && s.animSpriteNum > s.anims.GetLen(s.animName) {
		s.newStateAnim = STANDING
	}

	if s.stateAnim != s.newStateAnim {
		s.stateAnim = s.newStateAnim
		s.resetAnim()
	}

	s.animSpriteNum = int(math.Floor(s.counter / 0.1))

}

func (s *FreeState) Notify(e int, v *pixel.Vec) {
	// if event = attack or hit => switch state
	// if event = move => make it move

	if e == controller.E_CTRL {
		if !s.isJumping {
			s.a.SetState(STATE_ATTACK)
		}
	}
	if e == controller.E_SHIFT {
		s.isShift = !s.isShift
	}

	switch {
	case v.Y > 0:
		s.newStateAnim = JUMPING
	case v.Len() == 0:
		s.newStateAnim = STANDING
	case (v.Len() > 0 && s.isShift):
		s.newStateAnim = WALKING
	case (v.Len() > 0 && !s.isShift):
		s.newStateAnim = RUNNING
	}

	// set animation?
}

func (s *FreeState) GetSprite() *pixel.Sprite {
	if s.sprite == nil {
		s.sprite = pixel.NewSprite(nil, pixel.Rect{})
	}

	switch {
	case s.stateAnim == IDLE:
		s.animName = "idle"
	case s.stateAnim == STANDING:
		s.animName = "idle"
		s.animSpriteNum = 0
	case s.stateAnim == RUNNING:
		s.animName = "run"
	case s.stateAnim == WALKING:
		s.animName = "walk"
	case s.stateAnim == JUMPING:
		s.animName = "jump"
	}

	pic, rect := s.anims.GetSprite(s.animName, s.animSpriteNum)
	s.sprite.Set(pic, rect)

	return s.sprite
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
	//	return s.anims.GetSprite("attack1", 0)
	return nil
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
	return nil
	//	return s.anims.GetSprite("dead", 0)
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
	return nil
	//	return s.anims.GetSprite("hurt", 0)
}
