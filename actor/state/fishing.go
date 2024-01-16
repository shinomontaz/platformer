package state

import (
	"math"
	"platformer/common"
	"platformer/events"

	"github.com/shinomontaz/pixel"
)

type Fishing struct {
	Common
	idleLimit     float64
	counter       float64
	sprite        *pixel.Sprite
	animSpriteNum int
}

func NewFishing(a Actor, an Animater) *Fishing {
	fs := &Fishing{
		Common: Common{
			id:    FISHING,
			a:     a,
			anims: an,
			trs:   a.GetTransition(IDLE),
		},
		sprite:    pixel.NewSprite(nil, pixel.Rect{}),
		idleLimit: common.GetRandFloat() * 2.0, // seconds before idle
	}

	return fs
}

func (s *Fishing) Start() {
	s.a.Inform(events.RELEASED)
	s.counter = 0
	s.animSpriteNum = 0
	s.a.AddSound("fishing")
}

func (s *Fishing) Update(dt float64) {
	s.counter += dt
	s.animSpriteNum = int(math.Floor(s.counter / 0.2))
}

func (s *Fishing) Listen(e int, v *pixel.Vec) {
}

func (s *Fishing) GetSprite() *pixel.Sprite {
	if s.sprite == nil {
		s.sprite = pixel.NewSprite(nil, pixel.Rect{})
	}

	pic, rect := s.anims.GetSprite("fishing", s.animSpriteNum)
	s.sprite.Set(pic, rect)

	return s.sprite
}
