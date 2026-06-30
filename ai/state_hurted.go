package ai

import (
	"platformer/bindings"
	"platformer/events"

	"github.com/shinomontaz/pixel"
)

type StateHurted struct {
	id       int
	w        Worlder
	ai       *Ai
	reaction aiReactionType
	timer    float64
	duration float64 // сколько длится реакция
}

func NewHurted(ai *Ai, w Worlder, reaction aiReactionType) *StateHurted {
	return &StateHurted{
		id:       HURTED,
		ai:       ai,
		w:        w,
		reaction: reaction,
		duration: 0.5,
	}
}

func (s *StateHurted) Update(dt float64) {
	s.timer += dt
	if s.timer >= s.duration {
		// Реакция завершена — возвращаемся
		hero := s.ai.obj.GetEnemy()
		if hero != nil && hero.GetHp() > 0 {
			s.ai.SetState(CHOOSEATTACK, hero.GetPos())
		} else {
			s.ai.SetState(s.ai.homeState, pixel.ZV)
		}
		return
	}

	switch s.reaction {
	case HURTED_FLEE:
		// Отступить от героя
		pos := s.ai.obj.GetPos()
		hero := s.ai.obj.GetEnemy()
		if hero == nil {
			return
		}
		heropos := hero.GetPos()
		fleeDir := 1.0
		if heropos.X > pos.X {
			fleeDir = -1.0 // герой справа — бежим влево
		}
		v := pixel.Vec{fleeDir, 0}
		if s.ai.obj.StepPrediction(events.WALK, v) > 0.5 {
			b := bindings.Active.GetBinding(bindings.KeyAction["right"])
			if fleeDir < 0 {
				b = bindings.Active.GetBinding(bindings.KeyAction["left"])
			}
			s.ai.obj.KeyAction(b)
		}

	case HURTED_COUNTER:
		// Немедленная контратака — ничего не делаем,
		// duration=0 → сразу переходим в CHOOSEATTACK

	case HURTED_RANDOM:
		// Случайное движение (старый BUSTLE)
		// Просто стоим — можно добавить движение

	case HURTED_NONE:
		// Ничего не делаем
	}
}

func (s *StateHurted) EventAction(e int) {
	if e == events.GAMEVENT_HIT {
		// Повторное попадание — сбрасываем таймер
		s.timer = 0
	}
}

func (s *StateHurted) Start(poi pixel.Vec) {
	s.timer = 0
}

func (s *StateHurted) IsAlerted() bool {
	return true
}
