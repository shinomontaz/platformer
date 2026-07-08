package ai

import (
	"platformer/bindings"
	"platformer/events"
	"platformer/talks"

	"github.com/shinomontaz/pixel"
)

// StatePatrol — монстр патрулирует между spawnPos ± patrolRange.
// При IsSee(hero) → CHOOSEATTACK. При ALERT → INVESTIGATE.
// В отличие от ROAMING, патрулирующий монстр видит героя в обе стороны.
type StatePatrol struct {
	id          AiState
	w           Worlder
	ai          *Ai
	isbusy      bool
	spawnPos    pixel.Vec
	patrolRange float64
	dir         float64 // ±1 направление патруля
	waitTimer   float64
	waitDur     float64 // сколько стоять на точке разворота
	isWaiting   bool
}

func NewPatrol(ai *Ai, w Worlder, patrolRange float64) *StatePatrol {
	return &StatePatrol{
		id:          PATROL,
		ai:          ai,
		w:           w,
		patrolRange: patrolRange,
	}
}

func (s *StatePatrol) Update(dt float64) {
	if s.isbusy {
		return
	}

	pos := s.ai.obj.GetPos()

	// Патрулирующий монстр видит героя в ЛЮБУЮ сторону
	hero := s.ai.obj.GetEnemy()
	if hero != nil {
		herohp := hero.GetHp()
		heropos := hero.GetPos()
		if s.w.IsSee(pos, heropos) && herohp > 0 {
			talks.AddAlert(pos, 200)
			// Разворачиваемся к герою если нужно
			if (heropos.X < pos.X && s.dir > 0) || (heropos.X > pos.X && s.dir < 0) {
				s.dir = -s.dir
			}
			s.ai.SetState(CHOOSEATTACK, heropos)
			return
		}
	}

	// Стоим на точке разворота
	if s.isWaiting {
		s.waitTimer += dt
		if s.waitTimer >= s.waitDur {
			s.isWaiting = false
			s.dir = -s.dir
		}
		return
	}

	// Дошли до края зоны патрулирования?
	if s.dir > 0 && pos.X >= s.spawnPos.X+s.patrolRange {
		s.isWaiting = true
		s.waitTimer = 0
		s.waitDur = 1.0 + float64(int(pos.X*100)%20)/10.0 // 1.0-3.0с псевдослучайно
		return
	}
	if s.dir < 0 && pos.X <= s.spawnPos.X-s.patrolRange {
		s.isWaiting = true
		s.waitTimer = 0
		s.waitDur = 1.0 + float64(int(pos.X*100)%20)/10.0
		return
	}

	// Идём с проверкой что можно шагнуть
	v := pixel.V(s.dir, 0)
	groundrate := s.ai.obj.StepPrediction(events.WALK, v)

	if groundrate > 0.5 {
		b := bindings.Active.GetBinding(bindings.KeyAction["right"])
		if s.dir < 0 {
			b = bindings.Active.GetBinding(bindings.KeyAction["left"])
		}
		s.ai.obj.KeyAction(b)
	} else {
		// Нет земли — разворачиваемся
		s.dir = -s.dir
	}
}

func (s *StatePatrol) EventAction(e int) {
	switch e {
	case events.BUSY:
		s.isbusy = true
	case events.RELEASED:
		s.isbusy = false
	case events.ALERT:
		hero := s.ai.obj.GetEnemy()
		if hero != nil {
			s.ai.SetState(INVESTIGATE, hero.GetPos())
		}
	}
}

func (s *StatePatrol) Start(poi pixel.Vec) {
	s.spawnPos = s.ai.obj.GetPos()
	s.dir = 1.0
	s.isWaiting = false
	s.waitTimer = 0
	s.waitDur = 1.0
}

func (s *StatePatrol) IsAlerted() bool {
	return false
}
