package state

import "github.com/faiface/pixel"

type Dead struct {
	Common
}

func NewDead(a Actor, an Animater) *Dead {
	fs := &Dead{
		Common: Common{
			id:    DEAD,
			a:     a,
			anims: an,
		},
	}

	return fs
}

func (s *Dead) Start() {
}

func (s *Dead) Notify(e int, v *pixel.Vec) {
	// here we don't care of any controller event
}

func (s *Dead) Update(dt float64) {
}

func (s *Dead) GetSprite() *pixel.Sprite {
	return nil
	//	return s.anims.GetSprite("dead", 0)
}
