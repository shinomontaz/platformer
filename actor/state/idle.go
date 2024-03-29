package state

import (
	"math"
	"platformer/common"
	"platformer/events"

	"github.com/shinomontaz/pixel"
)

type Idle struct {
	Common
	idleLimit     float64
	counter       float64
	sprite        *pixel.Sprite
	animSpriteNum int
}

func NewIdle(a Actor, an Animater) *Idle {
	fs := &Idle{
		Common: Common{
			id:    IDLE,
			a:     a,
			anims: an,
			trs:   a.GetTransition(IDLE),
		},
		sprite:    pixel.NewSprite(nil, pixel.Rect{}),
		idleLimit: common.GetRandFloat() * 2.0, // seconds before idle
	}

	return fs
}

func (s *Idle) Start() {
	s.a.Inform(events.RELEASED)
	s.counter = 0
	s.animSpriteNum = 0
	s.a.AddSound("idle")
}

func (s *Idle) Update(dt float64) {
	s.counter += dt
	s.animSpriteNum = int(math.Floor(s.counter / 0.2))
	if s.animSpriteNum > 4 {
		s.a.SetState(STAND)
	}
}

func (s *Idle) Listen(e int, v *pixel.Vec) {
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

func (s *Idle) GetSprite() *pixel.Sprite {
	if s.sprite == nil {
		s.sprite = pixel.NewSprite(nil, pixel.Rect{})
	}

	pic, rect := s.anims.GetSprite("idle", s.animSpriteNum)
	s.sprite.Set(pic, rect)

	return s.sprite
}
