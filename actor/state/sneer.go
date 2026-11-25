package state

import (
	"math"
	"platformer/events"

	"github.com/shinomontaz/pixel"
)

type Sneer struct {
	Common
	sprite        *pixel.Sprite
	animSpriteNum int
	time          float64
	limit         float64
}

func NewSneer(a Actor, an Animater) *Sneer {
	return &Sneer{
		Common: Common{
			id:    SNEER,
			a:     a,
			anims: an,
			trs:   a.GetTransition(SNEER),
		},
		sprite: pixel.NewSprite(nil, pixel.Rect{}),
	}
}

func (s *Sneer) Start() {
	s.time = 0
	s.animSpriteNum = 0
	s.limit = 1.0
	if c, err := s.a.GetSkillAttr("cooldown"); err == nil {
		if v, ok := c.(float64); ok && v > 0 {
			s.limit = v
		}
	}
}

func (s *Sneer) Busy() bool { return true }

func (s *Sneer) Listen(e int, v *pixel.Vec) {
	s.a.Inform(events.BUSY)
}

func (s *Sneer) Update(dt float64) {
	if s.time >= s.limit {
		s.a.SetState(STAND)
		return
	}
	s.time += dt
	s.animSpriteNum = int(math.Floor(s.time / 0.2))
	if n := s.anims.GetLen("sneer"); n > 0 {
		s.animSpriteNum %= n // cooldown длиннее анимации — крутим по кругу
	}
}

func (s *Sneer) GetSprite() *pixel.Sprite {
	pic, rect := s.anims.GetSprite("sneer", s.animSpriteNum)
	s.sprite.Set(pic, rect)
	return s.sprite
}
