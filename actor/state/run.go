package state

import (
	"math"
	"platformer/events"

	"github.com/shinomontaz/pixel"
)

type Run struct {
	Common
	sprite        *pixel.Sprite
	animSpriteNum int
	counter       float64
}

func NewRun(a Actor, an Animater) *Run {
	fs := &Run{
		Common: Common{
			id:    RUN,
			a:     a,
			anims: an,
			trs:   a.GetTransition(RUN),
		},
		sprite: pixel.NewSprite(nil, pixel.Rect{}),
	}

	return fs
}

func (s *Run) Start() {
	s.a.Inform(events.RELEASED, pixel.ZV)
	s.counter = 0
}

func (s *Run) Update(dt float64) {
	s.counter += dt
	s.animSpriteNum = int(math.Floor(s.counter / 0.1))
}

func (s *Run) Listen(e int, v *pixel.Vec) {
	switch e {
	case events.WALK:
		s.a.SetState(WALK)
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

func (s *Run) GetSprite() *pixel.Sprite {
	pic, rect := s.anims.GetSprite("run", s.animSpriteNum)
	s.sprite.Set(pic, rect)

	return s.sprite
}
