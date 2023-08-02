package state

import (
	"math"
	"platformer/events"

	"github.com/shinomontaz/pixel"
)

type DeadSunk struct {
	Common
	counter       float64
	sprite        *pixel.Sprite
	animSpriteNum int
	animLen       int
}

func NewDeadSunk(a Actor, an Animater) *DeadSunk {
	fs := &DeadSunk{
		Common: Common{
			id:    DEADSUNK,
			a:     a,
			anims: an,
			busy:  true,
		},
		animLen: an.GetLen("die2"),
	}

	return fs
}

func (s *DeadSunk) Start() {
	s.a.Inform(events.BUSY)
	s.a.OnKill()
	s.counter = 0
}

func (s *DeadSunk) Listen(e int, v *pixel.Vec) {
	// here we don't care of any controller event
}

func (s *DeadSunk) Update(dt float64) {
	s.counter += dt
	if s.animSpriteNum < s.animLen-1 {
		s.animSpriteNum = int(math.Floor(s.counter / 0.1))
	}
}

func (s *DeadSunk) GetSprite() *pixel.Sprite {
	if s.sprite == nil {
		s.sprite = pixel.NewSprite(nil, pixel.Rect{})
	}

	pic, rect := s.anims.GetSprite("die2", s.animSpriteNum)
	s.sprite.Set(pic, rect)

	return s.sprite
}
