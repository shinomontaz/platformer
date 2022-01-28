package state

import (
	"math"
	"math/rand"
	"platformer/events"

	"github.com/faiface/pixel"
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
		idleLimit: rand.Float64() * 2.0, // seconds before idle
	}

	return fs
}

func (s *Idle) Start() {
	s.counter = 0
	s.animSpriteNum = 0
}

func (s *Idle) Update(dt float64) {
	s.counter += dt
	s.animSpriteNum = int(math.Floor(s.counter / 0.2))
	if s.counter > s.idleLimit {
		s.a.SetState(STAND)
	}
}

func (s *Idle) Notify(e int, v *pixel.Vec) {
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

	s.checkTransitions(e)
}

func (s *Idle) GetSprite() *pixel.Sprite {
	if s.sprite == nil {
		s.sprite = pixel.NewSprite(nil, pixel.Rect{})
	}

	pic, rect := s.anims.GetSprite("idle", s.animSpriteNum)
	s.sprite.Set(pic, rect)

	return s.sprite
}
