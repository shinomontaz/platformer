package state

import (
	"math"
	"platformer/common"
	"platformer/events"

	"github.com/shinomontaz/pixel"
)

type Sneer struct {
	Common
	idleLimit     float64
	counter       float64
	sprite        *pixel.Sprite
	animSpriteNum int
}

func NewSneer(a Actor, an Animater) *Sneer {
	fs := &Sneer{
		Common: Common{
			id:    SNEER,
			a:     a,
			anims: an,
			trs:   a.GetTransition(SNEER),
		},
		sprite:    pixel.NewSprite(nil, pixel.Rect{}),
		idleLimit: common.GetRandFloat() * 2.0, // seconds before idle
	}

	return fs
}

func (s *Sneer) Start() {
	s.a.Inform(events.RELEASED)
	s.counter = 0
	s.animSpriteNum = 0
	s.a.AddSound("sneer")
}

func (s *Sneer) Update(dt float64) {
	s.counter += dt
	s.animSpriteNum = int(math.Floor(s.counter / 0.2))
	if s.animSpriteNum > 4 {
		s.a.SetState(STAND) // TODO: next
	}
}

func (s *Sneer) Listen(e int, v *pixel.Vec) {
}

func (s *Sneer) GetSprite() *pixel.Sprite {
	if s.sprite == nil {
		s.sprite = pixel.NewSprite(nil, pixel.Rect{})
	}

	pic, rect := s.anims.GetSprite("sneer", s.animSpriteNum)
	s.sprite.Set(pic, rect)

	return s.sprite
}
