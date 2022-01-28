package state

import "github.com/faiface/pixel"

type Hit struct {
	Common
}

func NewHit(a Actor, an Animater) *Hit {
	fs := &Hit{
		Common: Common{
			id:    HIT,
			a:     a,
			anims: an,
			trs:   a.GetTransition(HIT),
		},
	}

	return fs
}

func (s *Hit) Start() {
}

func (s *Hit) Notify(e int, v *pixel.Vec) {
	// here we don't care of any controller event
	s.checkTransitions(e)
}

func (s *Hit) Update(dt float64) {
}

func (s *Hit) GetSprite() *pixel.Sprite {
	return nil
}
