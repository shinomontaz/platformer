package state

import (
	"math"
	"platformer/events"

	"github.com/shinomontaz/pixel"
)

type Dead struct {
	Common
	counter       float64
	sprite        *pixel.Sprite
	animSpriteNum int
	animLen       int
}

func NewDead(a Actor, an Animater) *Dead {
	fs := &Dead{
		Common: Common{
			id:    DEAD,
			a:     a,
			anims: an,
			busy:  true,
		},
		animLen: an.GetLen("die"),
	}

	return fs
}

func (s *Dead) Start() {
	s.a.Inform(events.BUSY, pixel.ZV)
	s.counter = 0
}

func (s *Dead) Listen(e int, v *pixel.Vec) {
	// here we don't care of any controller event
}

func (s *Dead) Update(dt float64) {
	s.counter += dt
	if s.animSpriteNum < s.animLen-1 {
		s.animSpriteNum = int(math.Floor(s.counter / 0.1))
	}
}

func (s *Dead) GetSprite() *pixel.Sprite {
	if s.sprite == nil {
		s.sprite = pixel.NewSprite(nil, pixel.Rect{})
	}

	pic, rect := s.anims.GetSprite("die", s.animSpriteNum)
	s.sprite.Set(pic, rect)

	return s.sprite
}
