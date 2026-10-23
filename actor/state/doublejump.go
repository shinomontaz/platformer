package state

import (
	"fmt"
	"math"

	"github.com/shinomontaz/pixel"
)

type DoubleJump struct {
	Common
	sprite          *pixel.Sprite
	animSpriteNum   int
	animSpriteMax   int
	counter         float64
	DoubleJumpLimit float64 // secs to change animation to deep DoubleJump
}

func NewDoubleJump(a Actor, an Animater) *DoubleJump {
	fs := &DoubleJump{
		Common: Common{
			id:    DOUBLEJUMP,
			a:     a,
			anims: an,
			trs:   a.GetTransition(DOUBLEJUMP),
		},
		animSpriteMax: an.GetLen("doublejump"),
		sprite:        pixel.NewSprite(nil, pixel.Rect{}),
	}

	return fs
}

func (s *DoubleJump) Start() {
	fmt.Println("DOUBLEJUMP start!")
	s.DoubleJumpLimit = 1.0
	s.animSpriteNum = 0
	s.counter = 0
}

func (s *DoubleJump) Update(dt float64) {
	s.counter += dt
	s.animSpriteNum = int(math.Floor(s.counter / 0.3))
}

func (s *DoubleJump) Listen(e int, v *pixel.Vec) {
	switch {
	case v.Y < 0 && (!s.iswater || (s.iswater && !s.iswaterresistant)):
		s.a.SetState(FALL)
	}

	s.checkTransitions(e, v)
}

func (s *DoubleJump) GetSprite() *pixel.Sprite {
	pic, rect := s.anims.GetSprite("doublejump", s.animSpriteNum)
	s.sprite.Set(pic, rect)

	return s.sprite
}
