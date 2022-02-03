package state

import (
	"math/rand"
	"platformer/events"

	"github.com/faiface/pixel"
)

type Stand struct {
	Common
	idleLimit float64
	counter   float64
	sprite    *pixel.Sprite
}

func NewStand(a Actor, an Animater) *Stand {
	fs := &Stand{
		Common: Common{
			id:    STAND,
			a:     a,
			anims: an,
			trs:   a.GetTransition(STAND),
		},
		sprite:    pixel.NewSprite(nil, pixel.Rect{}),
		idleLimit: rand.Float64() * 5.0, // seconds before idle
	}

	return fs
}

func (s *Stand) Start() {
	s.counter = 0
}

func (s *Stand) Update(dt float64) {
	s.counter += dt
	if s.counter > s.idleLimit {
		s.a.SetState(IDLE)
	}
}

func (s *Stand) Notify(e int, v *pixel.Vec) {
	switch e {
	case events.WALK:
		s.a.SetState(WALK)
		return
	case events.RUN:
		s.a.SetState(RUN)
		return
	}

	switch {
	case v.Y > 0:
		s.a.SetState(JUMP)
	case v.Y < 0:
		s.a.SetState(FALL)
	}

	s.checkTransitions(e, v)
}

func (s *Stand) GetSprite() *pixel.Sprite {
	if s.sprite == nil {
		s.sprite = pixel.NewSprite(nil, pixel.Rect{})
	}

	pic, rect := s.anims.GetSprite("idle", 0)
	s.sprite.Set(pic, rect)

	return s.sprite
}
