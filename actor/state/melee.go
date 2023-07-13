package state

import (
	"fmt"
	"math"
	"platformer/common"
	"platformer/events"

	"github.com/shinomontaz/pixel"
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
	striked       bool
	skillname     string
}

func NewMelee(a Actor, an Animater) *Melee {

	fs := &Melee{
		Common: Common{
			id:    MELEE,
			a:     a,
			anims: an,
			trs:   a.GetTransition(MELEE),
		},
		idleLimit: 0.5, // seconds before idle
	}

	return fs
}

func (s *Melee) Start() {
	s.a.Inform(events.BUSY, pixel.ZV)
	s.time = 0.0
	s.attackidx = 1
	s.striked = false

	skillname, err := s.a.GetSkillAttr("name")
	if err != nil {
		panic(err)
	}
	s.skillname = skillname.(string)
	if s.skillname == "" {
		s.skillname = "melee"
	}
	s.variants = s.anims.GetGroupLen(s.skillname)

	if s.variants > 1 {
		s.attackidx += int(math.Round(common.GetRandFloat() * float64(s.variants-1)))
	}
	// here add hitbox!
	s.a.AddSound("melee")
}

func (s *Melee) Listen(e int, v *pixel.Vec) {
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
	if s.animSpriteNum == 3 && !s.striked {
		s.a.Strike()
		s.striked = true
	}
}

func (s *Melee) GetSprite() *pixel.Sprite {
	if s.sprite == nil {
		s.sprite = pixel.NewSprite(nil, pixel.Rect{})
	}
	pic, rect := s.anims.GetGroupSprite(s.skillname, fmt.Sprintf("attack%d", s.attackidx), s.animSpriteNum)
	s.sprite.Set(pic, rect)

	return s.sprite
}
