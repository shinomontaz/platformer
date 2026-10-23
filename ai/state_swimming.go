package ai

import (
	"fmt"
	"math"
	"platformer/common/rand"

	"github.com/shinomontaz/pixel"
)

type StateSwimming struct {
	id        AiState
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
			s.timeout = 1 + rand.Float64()*2 // пауза
		} else {
			if rand.Float64() < 0.5 {
				s.leftright = -1
			} else {
				s.leftright = 1
			}
			s.timeout = 3 + rand.Float64()*4 // плывёт дольше
		}
		s.timer = 0
	}

	// Прямое управление скоростью: медленно плывём + мягкая синусоида по вертикали
	vel := s.ai.obj.GetVel()
	vel.X = s.leftright * 30
	vel.Y = math.Sin(s.totaltime) * 20 // слабее Архимеда — не пересилит выталкивание
	s.ai.obj.SetVel(vel)

	if s.leftright != 0 {
		s.ai.obj.SetDir(s.leftright)
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
