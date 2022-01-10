package actor

import "github.com/faiface/pixel"

type AttackState struct {
	CommonState
	time      float64
	idleLimit float64
}

func NewAttackState(a *Actor, an Animater) *AttackState {
	fs := &AttackState{
		CommonState: CommonState{
			id:    STATE_ATTACK,
			a:     a,
			anims: an,
		},
		idleLimit: 5.0, // seconds before idle
	}

	return fs
}

func (s *AttackState) Start() {
	s.time = 0.0
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
	//	s.pc.SetVec(pixel.ZV)
	//	s.pc.SetCmd(STRIKE)
}

func (s *AttackState) GetSprite() *pixel.Sprite {
	// return s.animations["attack1"].GetSprite(0)
	// return s.anims.GetSprite("attack1", 0)
	return nil
}
