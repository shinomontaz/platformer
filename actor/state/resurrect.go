package state

import (
	"math"
	"platformer/events"

	"github.com/shinomontaz/pixel"
)

type Resurrect struct {
	Common
	counter       float64
	sprite        *pixel.Sprite
	animSpriteNum int
	animLen       int
	lastCount     float64
}

func NewResurrect(a Actor, an Animater) *Resurrect {
	fs := &Resurrect{
		Common: Common{
			id:    RESURRECT,
			a:     a,
			anims: an,
			busy:  true,
		},
		animLen: an.GetLen("die"),
	}
	fs.animSpriteNum = fs.animLen - 1

	return fs
}

func (s *Resurrect) Start() {
	s.a.Inform(events.BUSY)
	s.a.OnKill()
	s.counter = 0
}

func (s *Resurrect) Listen(e int, v *pixel.Vec) {
	// here we don't care of any controller event
}

func (s *Resurrect) Update(dt float64) {
	s.counter += dt

	if s.animSpriteNum > 0 && s.lastCount != math.Floor(s.counter/0.05) {
		s.animSpriteNum--
	} else {
		s.a.SetState(STAND)
	}
}

func (s *Resurrect) GetSprite() *pixel.Sprite {
	if s.sprite == nil {
		s.sprite = pixel.NewSprite(nil, pixel.Rect{})
	}

	pic, rect := s.anims.GetSprite("die", s.animSpriteNum)
	s.sprite.Set(pic, rect)

	return s.sprite
}
