package state

import (
	"math"

	"github.com/faiface/pixel"
)

type Jump struct {
	Common
	sprite        *pixel.Sprite
	animSpriteNum int
	counter       float64
	jumpLimit     float64 // secs to change animation to deepjump
}

func NewJump(a Actor, an Animater) *Jump {
	fs := &Jump{
		Common: Common{
			id:    JUMP,
			a:     a,
			anims: an,
			trs:   a.GetTransition(JUMP),
		},
		sprite: pixel.NewSprite(nil, pixel.Rect{}),
	}

	return fs
}

func (s *Jump) Start() {
	s.jumpLimit = 1.0
	s.animSpriteNum = 0
	s.counter = 0
	s.a.AddSound("jump")
}

func (s *Jump) Update(dt float64) {
	s.counter += dt
	if s.animSpriteNum < 2 {
		s.animSpriteNum = int(math.Floor(s.counter / 0.2))
	}
}

func (s *Jump) Notify(e int, v *pixel.Vec) {
	// switch e {
	// case events.WALK:
	// 	s.a.SetState(WALK)
	// 	return
	// case events.RUN:
	// 	s.a.SetState(RUN)
	// 	return
	// }

	switch {
	case v.Y < 0:
		s.a.SetState(FALL)
	case v.Len() == 0:
		s.a.SetState(STAND)
	}

	s.checkTransitions(e, v)
}

func (s *Jump) GetSprite() *pixel.Sprite {
	pic, rect := s.anims.GetSprite("jump", s.animSpriteNum)
	s.sprite.Set(pic, rect)

	return s.sprite
}
