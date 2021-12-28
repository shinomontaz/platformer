package actor

import (
	"math"
	"platformer/controller"

	"github.com/faiface/pixel"
)

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
		if s.animSpriteNum > 0 {
			s.animSpriteNum = 1
		}
	case s.stateAnim == FALLING:
		s.animName = "jump"
		s.animSpriteNum = 4
	}

	pic, rect := s.anims.GetSprite(s.animName, s.animSpriteNum)
	s.sprite.Set(pic, rect)

	return s.sprite
}
