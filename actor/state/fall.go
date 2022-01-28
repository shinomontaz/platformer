package state

import (
	"platformer/events"

	"github.com/faiface/pixel"
)

type Fall struct {
	Common
	sprite        *pixel.Sprite
	animSpriteNum int
	counter       float64
	fallLimit     float64 // secs to change animation to deepfall
}

func NewFall(a Actor, an Animater) *Fall {
	fs := &Fall{
		Common: Common{
			id:    FALL,
			a:     a,
			anims: an,
			trs:   a.GetTransition(FALL),
		},
		sprite: pixel.NewSprite(nil, pixel.Rect{}),
	}

	return fs
}

func (s *Fall) Start() {
	s.fallLimit = 1.0
	s.animSpriteNum = 3
}

func (s *Fall) Update(dt float64) {
	s.counter += dt
	if s.counter > s.fallLimit {
		s.animSpriteNum = 4
	}
}

func (s *Fall) Notify(e int, v *pixel.Vec) {
	if v.Y == 0 {
		switch e {
		case events.WALK:
			s.a.SetState(WALK)
			return
		case events.RUN:
			s.a.SetState(RUN)
			return
		}
	}

	switch {
	case v.Y > 0:
		s.a.SetState(JUMP)
	case v.Len() == 0:
		s.a.SetState(STAND)
	}

	s.checkTransitions(e)
}

func (s *Fall) GetSprite() *pixel.Sprite {
	pic, rect := s.anims.GetSprite("jump", s.animSpriteNum)
	s.sprite.Set(pic, rect)

	return s.sprite
}
