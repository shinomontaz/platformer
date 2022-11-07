package state

import (
	"math"
	"platformer/events"

	"github.com/shinomontaz/pixel"
)

type Interact struct {
	Common
	time          float64
	idleLimit     float64
	attackidx     int
	animSpriteNum int
	sprite        *pixel.Sprite
	vel           float64
	interacted    bool
}

func NewInteract(a Actor, an Animater) *Interact {

	fs := &Interact{
		Common: Common{
			id:    INTERACT,
			a:     a,
			anims: an,
			trs:   a.GetTransition(INTERACT),
		},
		idleLimit: 0.5, // seconds before idle
		sprite:    pixel.NewSprite(nil, pixel.Rect{}),
	}

	return fs
}

func (s *Interact) Start() {
	s.a.Inform(events.BUSY, pixel.ZV)
	s.time = 0.0
	s.attackidx = 1
	s.interacted = false

	s.a.AddSound("Interact")
}

func (s *Interact) Listen(e int, v *pixel.Vec) {
	// here we don't care of any controller event
	s.vel = v.Len()
	s.checkTransitions(e, v)
}

func (s *Interact) Update(dt float64) {
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
	if s.animSpriteNum == 3 && !s.interacted {
		s.interacted = true
		s.a.Interact()
	}
}

func (s *Interact) GetSprite() *pixel.Sprite {
	pic, rect := s.anims.GetSprite("interact", s.animSpriteNum)
	s.sprite.Set(pic, rect)

	return s.sprite
}
