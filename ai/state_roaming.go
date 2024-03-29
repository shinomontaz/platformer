package ai

import (
	"fmt"
	"platformer/bindings"
	"platformer/common"
	"platformer/events"
	"platformer/talks"

	"github.com/shinomontaz/pixel/pixelgl"

	"github.com/shinomontaz/pixel"
)

type StateRoaming struct {
	id         int
	w          Worlder
	ai         *Ai
	isbusy     bool
	isagro     bool
	timeout    float64
	timer      float64
	dir        float64
	groundrate float64
}

func NewRoaming(ai *Ai, w Worlder, isagro bool) *StateRoaming {
	return &StateRoaming{
		id:      ROAMING,
		ai:      ai,
		w:       w,
		isagro:  isagro,
		timeout: 2,
	}
}

func (s *StateRoaming) Update(dt float64) {
	if s.isbusy {
		return
	}

	if s.isagro {
		pos := s.ai.obj.GetPos()

		hero := s.ai.obj.GetEnemy()
		if hero == nil {
			return
		}
		// look for hero
		herohp := hero.GetHp()
		heropos := hero.GetPos()
		dir := s.ai.obj.GetDir()
		if (heropos.X < pos.X && dir < 0) || (heropos.X > pos.X && dir > 0) {
			// check if we see hero
			if s.w.IsSee(pos, heropos) && herohp > 0 {
				talks.AddAlert(pos, 200)
				// 	al := addAlert(pos, force)
				// 	for _, en := range w.enemies {
				// 		alrect := al.GetRect()
				// 		if alrect.Contains(en.GetPos()) {
				// 			a := ai.GetByObj(en)
				// 			if a != nil {
				// 				a.Listen(events.ALERT, alrect.Center())
				// 			}
				// 		}
				// 	}

				s.ai.SetState(CHOOSEATTACK, heropos)
			}
			return
		}
	}

	s.timer += dt
	if s.timer > s.timeout {
		if s.dir != 0 {
			s.dir = 0
		} else {
			s.dir = float64(common.GetRandInt() - 5)
		}
		s.timer = 0
	}

	if s.dir != 0 {
		v := pixel.Vec{s.dir, 0}
		groundrate := s.ai.obj.StepPrediction(events.WALK, v)
		b := pixelgl.KeyUnknown

		if groundrate > 0.8 || groundrate > s.groundrate {
			if s.dir < 0 {
				b = bindings.Active.GetBinding(bindings.KeyAction["left"])
			} else {
				b = bindings.Active.GetBinding(bindings.KeyAction["right"])
			}
			s.ai.obj.KeyAction(b)
		}
		s.groundrate = groundrate
	}
}

func (s *StateRoaming) EventAction(e int) {
	//	fmt.Println("StateRoaming event action", e)
	if e == events.BUSY {
		s.isbusy = true
	}
	if e == events.RELEASED {
		s.isbusy = false
	}
}

func (s *StateRoaming) Start(poi pixel.Vec) {
	fmt.Println("state roaming")
}

func (s *StateRoaming) IsAlerted() bool {
	return false
}
