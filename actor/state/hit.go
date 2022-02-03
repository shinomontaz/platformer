package state

import (
	"math"

	"github.com/faiface/pixel"
)

type Hit struct {
	Common
	counter       float64
	timeout       float64
	sprite        *pixel.Sprite
	animSpriteNum int
}

func NewHit(a Actor, an Animater) *Hit {
	fs := &Hit{
		Common: Common{
			id:    HIT,
			a:     a,
			anims: an,
			trs:   a.GetTransition(HIT),
			busy:  true,
		},
		timeout: 1,
	}

	return fs
}

func (s *Hit) Start() {
	s.counter = 0
}

func (s *Hit) Notify(e int, v *pixel.Vec) {
	// here we don't care of any controller event
	//	s.checkTransitions(e, v)
}

func (s *Hit) Update(dt float64) {
	s.counter += dt
	s.animSpriteNum = int(math.Floor(s.counter / 0.2))
	if s.counter > s.timeout {
		s.a.SetState(STAND)
	}
}

func (s *Hit) GetSprite() *pixel.Sprite {
	if s.sprite == nil {
		s.sprite = pixel.NewSprite(nil, pixel.Rect{})
	}

	pic, rect := s.anims.GetSprite("hurt", s.animSpriteNum)
	s.sprite.Set(pic, rect)

	return s.sprite
}
