package state

import (
	"fmt"
	"math"
	"platformer/common"
	"platformer/events"

	"github.com/shinomontaz/pixel"
)

type Attack struct {
	Common
	time          float64
	idleLimit     float64
	attackidx     int
	animSpriteNum int
	sprite        *pixel.Sprite
	variants      int
}

func NewAttack(a Actor, an Animater) *Attack {
	fs := &Attack{
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

func (s *Attack) Start() {
	s.time = 0.0
	s.attackidx = 1
	if s.variants > 1 {
		s.attackidx += int(math.Round(common.GetRandFloat() * float64(s.variants-1)))
	}
}

func (s *Attack) Listen(e int, v *pixel.Vec) {
	s.a.Inform(events.BUSY, pixel.ZV)
	// here we don't care of any controller event
	s.checkTransitions(e, v)
}

func (s *Attack) Update(dt float64) {
	if s.time > s.idleLimit {
		// TODO: return to specific "free" where actual state will be detected
		s.a.SetState(STAND)
		return
	}

	s.time += dt
	s.animSpriteNum = int(math.Floor(s.time / 0.1))
	//	s.pc.SetVec(pixel.ZV)
	//	s.pc.SetCmd(STRIKE)
}

func (s *Attack) GetSprite() *pixel.Sprite {
	if s.sprite == nil {
		s.sprite = pixel.NewSprite(nil, pixel.Rect{})
	}
	pic, rect := s.anims.GetGroupSprite("melee", fmt.Sprintf("attack%d", s.attackidx), s.animSpriteNum)
	s.sprite.Set(pic, rect)

	return s.sprite
}
