package state

import (
	"fmt"
	"math"
	"math/rand"

	"github.com/faiface/pixel"
)

type Melee struct {
	Common
	time          float64
	idleLimit     float64
	attackidx     int
	animSpriteNum int
	sprite        *pixel.Sprite
	variants      int
	vel           float64
}

func NewMelee(a Actor, an Animater) *Melee {
	fs := &Melee{
		Common: Common{
			id:    ATTACK,
			a:     a,
			anims: an,
			trs:   a.GetTransition(ATTACK),
		},
		idleLimit: 0.5, // seconds before idle
		variants:  an.GetGroupLen("melee"),
	}

	return fs
}

func (s *Melee) Start() {
	s.time = 0.0
	s.attackidx = 1
	if s.variants > 1 {
		s.attackidx += rand.Intn(s.variants)
	}
	// here add hitbox!
	s.a.Strike()
}

func (s *Melee) Notify(e int, v *pixel.Vec) {
	// here we don't care of any controller event
	s.vel = v.Len()
	s.checkTransitions(e, v)
}

func (s *Melee) Update(dt float64) {
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
	//	s.pc.SetVec(pixel.ZV)
	//	s.pc.SetCmd(STRIKE)
}

func (s *Melee) GetSprite() *pixel.Sprite {
	if s.sprite == nil {
		s.sprite = pixel.NewSprite(nil, pixel.Rect{})
	}
	pic, rect := s.anims.GetGroupSprite("melee", fmt.Sprintf("attack%d", s.attackidx), s.animSpriteNum)
	s.sprite.Set(pic, rect)

	return s.sprite
}
