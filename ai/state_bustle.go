package ai

import (
	"fmt"
	"platformer/bindings"
	"platformer/common"
	"platformer/events"

	"github.com/shinomontaz/pixel"
	"github.com/shinomontaz/pixel/pixelgl"
)

// Состояние суеты. Используется между атаками с некоторой вероятностью, особенно между дальними
type StateBustle struct {
	id         int
	w          Worlder
	ai         *Ai
	isbusy     bool
	timeout    float64
	timer      float64
	dir        float64
	groundrate float64
	isswitched bool
	poi        pixel.Vec
}

func NewBustle(ai *Ai, w Worlder) *StateBustle {
	return &StateBustle{
		id: BUSTLE,
		ai: ai,
		w:  w,
	}
}

func (s *StateBustle) Update(dt float64) {
	if s.isbusy {
		return
	}

	hero := s.ai.obj.GetEnemy()
	if hero == nil {
		return
	}
	s.timer += dt
	if s.timer > s.timeout {
		//go to previos state
		pos := s.ai.obj.GetPos()
		herohp := hero.GetHp()
		heropos := hero.GetPos()
		dir := s.ai.obj.GetDir()
		if (heropos.X < pos.X && dir < 0) || (heropos.X > pos.X && dir > 0) {
			// check if we see hero
			if s.w.IsSee(pos, heropos) && herohp > 0 {
				s.ai.SetState(CHOOSEATTACK, heropos)
			}
		} else {
			s.ai.SetState(INVESTIGATE, heropos)
		}
	}

	if s.timer > s.timeout/2 && !s.isswitched {
		if s.dir != 0 {
			s.dir = 0
		} else {
			s.dir = float64(common.GetRandInt() - 5)
		}
		s.isswitched = true
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

func (s *StateBustle) EventAction(e int) {
	if e == events.BUSY {
		s.isbusy = true
	}
	if e == events.RELEASED {
		s.isbusy = false
	}
}

func (s *StateBustle) Start(poi pixel.Vec) {
	s.timeout = 0.5 + float64(common.GetRandInt())/20
	s.dir = float64(common.GetRandInt() - 5)
	fmt.Println("bustle start", s.timeout, s.dir)
	s.isswitched = false
	s.poi = poi
}

func (s *StateBustle) IsAlerted() bool {
	return false
}
