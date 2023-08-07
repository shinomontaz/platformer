package state

import (
	"fmt"
	"math"
	"platformer/common"
	"platformer/events"

	"github.com/shinomontaz/pixel"
)

type Meleemove struct {
	Common
	time          float64
	idleLimit     float64
	attackidx     int
	animSpriteNum int
	sprite        *pixel.Sprite
	vel           pixel.Vec
	striked       bool
}

func NewMeleemove(a Actor, an Animater) *Meleemove {
	fs := &Meleemove{
		Common: Common{
			id:    MELEEMOVE,
			a:     a,
			anims: an,
			trs:   a.GetTransition(MELEEMOVE),
			busy:  true,
		},
	}

	return fs
}

func (s *Meleemove) Start() {
	s.a.Inform(events.BUSY)
	s.time = 0.0
	s.attackidx = 1

	res, err := s.a.GetSkillAttr("ttl")
	if err != nil {
		panic(err)
	}

	s.idleLimit = res.(float64)
	variants := s.anims.GetGroupLen("meleemove")
	if variants > 1 {
		s.attackidx += int(math.Round(common.GetRandFloat() * float64(variants-1)))
	}

	s.striked = false
	s.a.AddSound("melee")
	s.vel = s.a.GetVel()

	res2, err := s.a.GetSkillAttr("speed")
	if err != nil {
		panic(err)
	}

	s.vel.X *= res2.(float64) / math.Abs(s.vel.X)
}

func (s *Meleemove) Listen(e int, v *pixel.Vec) {
	// here we don't care of any controller event
	// s.vel = v.Len()
	// s.checkTransitions(e, v)
}

func (s *Meleemove) Update(dt float64) {
	if s.time > s.idleLimit {
		s.a.SetState(STAND)
		return
	}

	// add speed - we are moving

	s.time += dt
	s.animSpriteNum = int(math.Floor(s.time / 0.1))
	if s.animSpriteNum == 3 && !s.striked {
		s.a.Strike()
		s.striked = true
	}
	s.a.SetVel(s.vel)
}

func (s *Meleemove) GetSprite() *pixel.Sprite {
	if s.sprite == nil {
		s.sprite = pixel.NewSprite(nil, pixel.Rect{})
	}

	pic, rect := s.anims.GetGroupSprite("meleemove", fmt.Sprintf("attack%d", s.attackidx), s.animSpriteNum)
	s.sprite.Set(pic, rect)

	return s.sprite
}
