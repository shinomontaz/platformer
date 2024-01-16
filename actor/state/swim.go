package state

import (
	"fmt"
	"math"
	"platformer/common"

	"github.com/shinomontaz/pixel"
)

type Swim struct {
	Common
	sprite        *pixel.Sprite
	animSpriteNum int
	counter       float64
	idleLimit     float64
}

func NewSwim(a Actor, an Animater) *Swim {
	fs := &Swim{
		Common: Common{
			id:               SWIM,
			a:                a,
			anims:            an,
			trs:              a.GetTransition(FALL),
			iswaterresistant: true,
			iswater:          true,
		},
		sprite:    pixel.NewSprite(nil, pixel.Rect{}),
		idleLimit: common.GetRandFloat() * 5.0,
	}

	return fs
}

func (s *Swim) Start() {
	fmt.Println("state swim")
}

func (s *Swim) Update(dt float64) {
	s.counter += dt
	if s.counter > s.idleLimit {
		s.a.SetState(IDLE)
	}
	s.animSpriteNum = int(math.Floor(s.counter / 0.15))
}

func (s *Swim) Listen(e int, v *pixel.Vec) {
	//	s.checkTransitions(e, v)
}

func (s *Swim) GetSprite() *pixel.Sprite {
	pic, rect := s.anims.GetSprite("swim", s.animSpriteNum)
	s.sprite.Set(pic, rect)

	return s.sprite
}
