package actor

import "github.com/faiface/pixel"

type HitState struct {
	CommonState
	time      float64
	timelimit float64
}

func NewHitState(a *Actor, an Animater) *HitState {
	fs := &HitState{
		CommonState: CommonState{
			id:    STATE_HIT,
			a:     a,
			anims: an,
		},
	}

	return fs
}

func (s *HitState) Start() {
}

func (s *HitState) Notify(e int, v *pixel.Vec) {
	// here we don't care of any controller event
}

func (s *HitState) Update(dt float64) {
	if s.time > s.timelimit {
		s.a.SetState(STATE_FREE)
		return
	}

	s.time += dt
}

func (s *HitState) GetSprite() *pixel.Sprite {
	return nil
	//	return s.anims.GetSprite("hurt", 0)
}
