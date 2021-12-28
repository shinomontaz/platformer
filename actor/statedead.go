package actor

import "github.com/faiface/pixel"

type DeadState struct {
	CommonState
}

func NewDeadState(a *Actor, an Animater) *DeadState {
	fs := &DeadState{
		CommonState: CommonState{
			id:    STATE_DEAD,
			a:     a,
			anims: an,
		},
	}

	return fs
}

func (s *DeadState) Start() {
}

func (s *DeadState) Notify(e int, v *pixel.Vec) {
	// here we don't care of any controller event
}

func (s *DeadState) Update(dt float64) {
}

func (s *DeadState) GetSprite() *pixel.Sprite {
	return nil
	//	return s.anims.GetSprite("dead", 0)
}
