package actor

import (
	"fmt"
	"math"
	"math/rand"

	"github.com/faiface/pixel"
)

type AttackState struct {
	CommonState
	time          float64
	idleLimit     float64
	attackidx     int
	animSpriteNum int
	sprite        *pixel.Sprite
}

func NewAttackState(a *Actor, an Animater) *AttackState {
	fs := &AttackState{
		CommonState: CommonState{
			id:    STATE_ATTACK,
			a:     a,
			anims: an,
		},
		idleLimit: 0.5, // seconds before idle
	}

	return fs
}

func (s *AttackState) Start() {
	s.time = 0.0
	s.attackidx = rand.Intn(2) + 1
}

func (s *AttackState) Notify(e int, v *pixel.Vec) {
	// here we don't care of any controller event
}

func (s *AttackState) Update(dt float64) {
	if s.time > s.idleLimit {
		s.a.SetState(STATE_FREE)
		return
	}

	s.time += dt
	s.animSpriteNum = int(math.Floor(s.time / 0.1))
	//	s.pc.SetVec(pixel.ZV)
	//	s.pc.SetCmd(STRIKE)
}

func (s *AttackState) GetSprite() *pixel.Sprite {
	if s.sprite == nil {
		s.sprite = pixel.NewSprite(nil, pixel.Rect{})
	}
	pic, rect := s.anims.GetSprite(fmt.Sprintf("attack%d", s.attackidx), s.animSpriteNum)
	s.sprite.Set(pic, rect)

	return s.sprite
}
