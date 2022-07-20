package state

import (
	"math"
	"platformer/events"

	"github.com/shinomontaz/pixel"
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
	s.a.Inform(events.BUSY, pixel.ZV)
	s.counter = 0
	// aa := s.a.GetAi()
	// if aa != nil {
	// 	aa.Notify(events.ALERT, pixel.ZV)
	// }
	s.a.AddSound("hit")
}

func (s *Hit) Listen(e int, v *pixel.Vec) {
	// here we don't care of any controller event
	//	s.checkTransitions(e, v)
}

func (s *Hit) Update(dt float64) {
	s.counter += dt
	s.animSpriteNum = int(math.Floor(s.counter / 0.05))
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
