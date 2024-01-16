package state

import (
	"math"
	"platformer/events"

	"github.com/shinomontaz/pixel"
)

type Roll struct {
	Common
	sprite        *pixel.Sprite
	animSpriteNum int
	counter       float64
	animSpriteMax int
	vel           float64
}

func NewRoll(a Actor, an Animater) *Roll {
	fs := &Roll{
		Common: Common{
			id:    ROLL,
			a:     a,
			anims: an,
			trs:   a.GetTransition(ROLL),
		},
		sprite:        pixel.NewSprite(nil, pixel.Rect{}),
		animSpriteMax: an.GetLen("roll"),
	}

	return fs
}

func (s *Roll) Start() {
	s.animSpriteNum = 0
	s.a.Inform(events.BUSY)
	s.counter = 0
}

func (s *Roll) Update(dt float64) {

	if s.animSpriteNum > s.animSpriteMax {
		if s.vel > 0 {
			s.a.SetState(WALK)
		} else {
			s.a.SetState(STAND)
		}
		return
	}

	s.counter += dt
	s.animSpriteNum = int(math.Floor(s.counter / 0.1))
}

func (s *Roll) Listen(e int, v *pixel.Vec) {
	s.vel = v.Len()
}

func (s *Roll) GetSprite() *pixel.Sprite {
	pic, rect := s.anims.GetSprite("roll", s.animSpriteNum)
	s.sprite.Set(pic, rect)

	return s.sprite
}
