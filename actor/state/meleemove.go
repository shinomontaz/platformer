package state

import (
	"math"

	"github.com/faiface/pixel"
)

type Meleemove struct {
	Common
	time          float64
	idleLimit     float64
	attackidx     int
	animSpriteNum int
	sprite        *pixel.Sprite
	vel           float64
}

func NewMeleemove(a Actor, an Animater) *Meleemove {
	fs := &Meleemove{
		Common: Common{
			id:    ATTACK,
			a:     a,
			anims: an,
			trs:   a.GetTransition(ATTACK),
		},
		idleLimit: 0.5, // seconds before idle
	}

	return fs
}

func (s *Meleemove) Start() {
	s.time = 0.0
	s.attackidx = 3
}

func (s *Meleemove) Notify(e int, v *pixel.Vec) {
	// here we don't care of any controller event
	s.vel = v.Len()
	s.checkTransitions(e, v)
}

func (s *Meleemove) Update(dt float64) {
	if s.time > s.idleLimit {
		// TODO: return to specific "free" where actual state will be detected
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

func (s *Meleemove) GetSprite() *pixel.Sprite {
	if s.sprite == nil {
		s.sprite = pixel.NewSprite(nil, pixel.Rect{})
	}
	pic, rect := s.anims.GetSprite("meleemove", s.animSpriteNum)
	s.sprite.Set(pic, rect)

	return s.sprite
}
