package ai

import (
	"fmt"
	"platformer/actor/state"

	"github.com/shinomontaz/pixel"
)

type StateFishing struct {
	id     int
	w      Worlder
	ai     *Ai
	isbusy bool
}

func NewFishing(ai *Ai, w Worlder, isagro bool) *StateFishing {
	return &StateFishing{
		id: FISHING,
		ai: ai,
		w:  w,
	}
}

func (s *StateFishing) Update(dt float64) {
	if s.isbusy {
		return
	}
}

func (s *StateFishing) Start(poi pixel.Vec) {
	fmt.Println("state fishing ai")
	s.ai.obj.SetState(state.FISHING)
}

func (s *StateFishing) EventAction(e int) {
}

func (s *StateFishing) GetVec() pixel.Vec {
	return pixel.ZV
}

func (s *StateFishing) IsAlerted() bool {
	return false
}
