package ai

import (
	"platformer/bindings"
	"platformer/common/rand"
	"platformer/events"

	"github.com/shinomontaz/pixel"
)

// Состояние суеты. Монстр ходит туда-обратно вокруг poi с фиксированной амплитудой.
// Используется между атаками, особенно после дальних атак.
type StateBustle struct {
	id            int
	w             Worlder
	ai            *Ai
	isbusy        bool
	timeout       float64
	timer         float64
	phase         int // 0 = вправо, 1 = влево
	phaseTimer    float64
	phaseDuration float64 // длительность одного шага (секунды)
	poi           pixel.Vec
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
	s.phaseTimer += dt

	// Переключаем направление каждые phaseDuration секунд
	if s.phaseTimer > s.phaseDuration {
		s.phase = 1 - s.phase
		s.phaseTimer = 0
	}

	// Фиксированное направление: +1 или -1
	dir := 1.0
	if s.phase == 1 {
		dir = -1.0
	}

	// StepPrediction: безопасно ли идти? Не падаем в воду/пропасть
	v := pixel.Vec{dir, 0}
	groundrate := s.ai.obj.StepPrediction(events.WALK, v)

	if groundrate > 0.5 {
		b := bindings.Active.GetBinding(bindings.KeyAction["right"])
		if dir < 0 {
			b = bindings.Active.GetBinding(bindings.KeyAction["left"])
		}
		s.ai.obj.KeyAction(b)
	} else {
		// Нет земли впереди — разворачиваемся
		s.phase = 1 - s.phase
		s.phaseTimer = 0
	}

	// Таймаут: проверяем видим ли героя
	if s.timer > s.timeout {
		pos := s.ai.obj.GetPos()
		herohp := hero.GetHp()
		heropos := hero.GetPos()
		facedir := s.ai.obj.GetDir()
		if (heropos.X < pos.X && facedir < 0) || (heropos.X > pos.X && facedir > 0) {
			if s.w.IsSee(pos, heropos) && herohp > 0 {
				s.ai.SetState(CHOOSEATTACK, heropos)
				return
			}
		}
		s.ai.SetState(INVESTIGATE, heropos)
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
	s.timeout = 2.0 + rand.Float64()*2.0 // 2-4 секунды
	s.phaseDuration = 0.8                // фиксированная длительность шага
	s.phase = 0
	s.phaseTimer = 0
	s.timer = 0
	s.poi = poi
}

func (s *StateBustle) IsAlerted() bool {
	return false
}
