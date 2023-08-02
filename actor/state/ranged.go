package state

import (
	"math"
	"platformer/events"

	"github.com/shinomontaz/pixel"
)

type Ranged struct {
	Common
	time          float64
	idleLimit     float64
	animSpriteNum int
	sprite        *pixel.Sprite
	variants      int
	vel           float64
}

func NewRanged(a Actor, an Animater) *Ranged {
	fs := &Ranged{
		Common: Common{
			id:    RANGED,
			a:     a,
			anims: an,
			trs:   a.GetTransition(RANGED),
		},
		idleLimit: 3, // seconds before state transition
		variants:  an.GetLen("cast"),
	}

	return fs
}

func (s *Ranged) Start() {
	s.a.Inform(events.BUSY)
	s.time = 0.0
	s.a.Cast()
}

func (s *Ranged) Listen(e int, v *pixel.Vec) {
	// here we don't care of any controller event
	s.vel = v.Len()
	s.checkTransitions(e, v)
}

func (s *Ranged) Update(dt float64) {
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

func (s *Ranged) GetSprite() *pixel.Sprite {
	if s.sprite == nil {
		s.sprite = pixel.NewSprite(nil, pixel.Rect{})
	}
	pic, rect := s.anims.GetSprite("cast", s.animSpriteNum)
	s.sprite.Set(pic, rect)

	return s.sprite
}
