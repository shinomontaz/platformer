package state

import (
	"math"
	"platformer/events"

	"github.com/faiface/pixel"
)

type Walk struct {
	Common
	sprite        *pixel.Sprite
	animSpriteNum int
	counter       float64
}

func NewWalk(a Actor, an Animater) *Walk {
	fs := &Walk{
		Common: Common{
			id:    WALK,
			a:     a,
			anims: an,
			trs:   a.GetTransition(WALK),
		},
		sprite: pixel.NewSprite(nil, pixel.Rect{}),
	}

	return fs
}

func (s *Walk) Start() {
	s.a.Inform(events.RELEASED, pixel.ZV)
	s.counter = 0
}

func (s *Walk) Update(dt float64) {
	s.counter += dt
	s.animSpriteNum = int(math.Floor(s.counter / 0.15))
}

func (s *Walk) Listen(e int, v *pixel.Vec) {
	switch e {
	case events.RUN:
		s.a.SetState(RUN)
		return
	}

	switch {
	case v.Y > 0:
		s.a.SetState(JUMP)
	case v.Y < 0:
		s.a.SetState(FALL)
	case v.Len() == 0:
		s.a.SetState(STAND)
	}

	s.checkTransitions(e, v)
}

func (s *Walk) GetSprite() *pixel.Sprite {
	pic, rect := s.anims.GetSprite("walk", s.animSpriteNum)
	s.sprite.Set(pic, rect)

	return s.sprite
}
