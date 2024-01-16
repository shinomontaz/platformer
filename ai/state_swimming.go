package ai

import (
	"fmt"
	"math"
	"platformer/bindings"
	"platformer/common"

	"github.com/shinomontaz/pixel/pixelgl"

	"github.com/shinomontaz/pixel"
)

type StateSwimming struct {
	id        int
	w         Worlder
	ai        *Ai
	timeout   float64
	timer     float64
	totaltime float64
	leftright float64
	updown    float64
}

func NewSwimming(ai *Ai, w Worlder, isagro bool) *StateSwimming {
	return &StateSwimming{
		id:      SWIMMING,
		ai:      ai,
		w:       w,
		timeout: 2,
	}
}

func (s *StateSwimming) Update(dt float64) {
	s.timer += dt
	s.totaltime += dt
	if s.timer > s.timeout {
		if s.leftright != 0 {
			s.leftright = 0
		} else {
			s.leftright = float64(common.GetRandInt() - 5)
		}
		s.timer = 0
	}

	s.updown = math.Sin(s.totaltime)

	if s.leftright != 0 {
		b := pixelgl.KeyUnknown

		if s.leftright < 0 {
			b = bindings.Active.GetBinding(bindings.KeyAction["left"])
		} else {
			b = bindings.Active.GetBinding(bindings.KeyAction["right"])
		}
		s.ai.obj.KeyAction(b)
	}

	if s.updown != 0 {
		b := pixelgl.KeyUnknown
		if s.updown < 0 {
			b = bindings.Active.GetBinding(bindings.KeyAction["down"])
		} else {
			b = bindings.Active.GetBinding(bindings.KeyAction["up"])
		}
		s.ai.obj.KeyAction(b)
	}
}

func (s *StateSwimming) EventAction(e int) {
}

func (s *StateSwimming) Start(poi pixel.Vec) {
	fmt.Println("state swimming")
}

func (s *StateSwimming) IsAlerted() bool {
	return false
}
