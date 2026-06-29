package state

import (
	"math"

	"github.com/shinomontaz/pixel"
)

type Sit struct {
	Common
	time          float64
	animSpriteNum int
	animSpriteMax int
	sprite        *pixel.Sprite
}

func NewSit(a Actor, an Animater) *Sit {

	fs := &Sit{
		Common: Common{
			id:    SIT,
			a:     a,
			anims: an,
			trs:   a.GetTransition(SIT),
		},
		animSpriteMax: an.GetLen("sit"),
	}

	return fs
}

func (s *Sit) Start() {
	s.animSpriteNum = 0
	s.time = 0.0

}

func (s *Sit) Listen(e int, v *pixel.Vec) {
	// here we don't care of any controller event
	if v.Y >= 0 {
		s.a.SetState(STAND)
	}
	s.checkTransitions(e, v)
}

func (s *Sit) Update(dt float64) {
	s.time += dt
	s.animSpriteNum = int(math.Floor(s.time / 0.1))
}

func (s *Sit) GetSprite() *pixel.Sprite {
	if s.sprite == nil {
		s.sprite = pixel.NewSprite(nil, pixel.Rect{})
	}
	pic, rect := s.anims.GetSprite("sit", s.animSpriteNum)
	s.sprite.Set(pic, rect)

	return s.sprite
}
