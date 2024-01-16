package state

import (
	"math"
	"platformer/events"

	"github.com/shinomontaz/pixel"
)

type Cast struct {
	Common
	time          float64
	idleLimit     float64
	animSpriteNum int
	sprite        *pixel.Sprite
	variants      int
	vel           float64
}

func NewCast(a Actor, an Animater) *Cast {
	fs := &Cast{
		Common: Common{
			id:    CAST,
			a:     a,
			anims: an,
			trs:   a.GetTransition(CAST),
		},
		idleLimit: 1.5, // seconds before state transition
		variants:  an.GetLen("cast"),
	}

	return fs
}

func (s *Cast) Start() {
	s.a.Inform(events.BUSY)
	s.time = 0.0
	s.a.UseSkill()
}

func (s *Cast) Listen(e int, v *pixel.Vec) {
	// here we don't care of any controller event
	s.vel = v.Len()
	s.checkTransitions(e, v)
}

func (s *Cast) Update(dt float64) {
	if s.time > s.idleLimit {
		if s.vel > 0 {
			s.a.SetState(WALK)
		} else {
			s.a.SetState(STAND)
		}
		return
	}

	s.time += dt
	s.animSpriteNum = int(math.Floor(s.time / 0.1))
}

func (s *Cast) GetSprite() *pixel.Sprite {
	if s.sprite == nil {
		s.sprite = pixel.NewSprite(nil, pixel.Rect{})
	}
	pic, rect := s.anims.GetSprite("cast", s.animSpriteNum)
	s.sprite.Set(pic, rect)

	return s.sprite
}
