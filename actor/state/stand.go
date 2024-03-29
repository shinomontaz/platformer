package state

import (
	"platformer/common"
	"platformer/events"

	"github.com/shinomontaz/pixel"
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
		idleLimit: common.GetRandFloat() * 5.0, // seconds before idle
	}

	return fs
}

func (s *Stand) Start() {
	s.a.Inform(events.RELEASED)
	s.counter = 0
}

func (s *Stand) Update(dt float64) {
	s.counter += dt
	if s.counter > s.idleLimit {
		s.a.SetState(IDLE)
	}
}

func (s *Stand) Listen(e int, v *pixel.Vec) {
	switch e {
	case events.INTERACT:
		s.a.SetState(INTERACT)
		return
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
	case v.Y < 0 && (!s.iswater || (s.iswater && !s.iswaterresistant)):
		s.a.SetState(FALL)
	case v.X != 0:
		s.a.SetState(WALK)
		return
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
