package ai

import (
	"fmt"
	"platformer/bindings"
	"platformer/common/rand"
	"platformer/events"
	"platformer/talks"

	"github.com/shinomontaz/pixel"
)

type StateRoaming struct {
	id            int
	w             Worlder
	ai            *Ai
	isbusy        bool
	isagro        bool
	timeout       float64
	timer         float64
	dir           float64 // ±1 или 0
	phaseTimer    float64
	phaseDuration float64
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
		herohp := hero.GetHp()
		heropos := hero.GetPos()
		dir := s.ai.obj.GetDir()
		if (heropos.X < pos.X && dir < 0) || (heropos.X > pos.X && dir > 0) {
			if s.w.IsSee(pos, heropos) && herohp > 0 {
				talks.AddAlert(pos, 200)
				s.ai.SetState(CHOOSEATTACK, heropos)
			}
			return
		}
	}

	s.timer += dt
	s.phaseTimer += dt

	// Переключаем направление/остановку по таймеру
	if s.timer > s.timeout {
		if s.dir != 0 {
			s.dir = 0
			s.timeout = 1.0 + rand.Float64()*2.0 // стоим 1-3 секунды
		} else {
			if rand.Float64() > 0.5 {
				s.dir = 1
			} else {
				s.dir = -1
			}
			s.timeout = 2.0 + rand.Float64()*3.0 // идём 2-5 секунд
		}
		s.timer = 0
	}

	if s.dir != 0 {
		v := pixel.Vec{s.dir, 0}
		groundrate := s.ai.obj.StepPrediction(events.WALK, v)

		if groundrate > 0.5 {
			b := bindings.Active.GetBinding(bindings.KeyAction["right"])
			if s.dir < 0 {
				b = bindings.Active.GetBinding(bindings.KeyAction["left"])
			}
			s.ai.obj.KeyAction(b)
		} else {
			// Нет земли впереди — разворачиваемся
			s.dir = -s.dir
			s.timer = 0
		}
	}
}

func (s *StateRoaming) EventAction(e int) {
	switch e {
	case events.BUSY:
		s.isbusy = true
	case events.RELEASED:
		s.isbusy = false
	case events.HIT:
		fmt.Println("Hit!")
	}
}

func (s *StateRoaming) Start(poi pixel.Vec) {
	s.dir = 0
	s.timer = 0
	s.phaseTimer = 0
	s.timeout = 1.0 + rand.Float64()*2.0
}

func (s *StateRoaming) IsAlerted() bool {
	return false
}
