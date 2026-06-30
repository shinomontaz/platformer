package ai

import (
	"math"
	"platformer/actor"
	"platformer/bindings"
	"platformer/common/rand"
	"platformer/events"
	"sort"

	"github.com/shinomontaz/pixel"
)

type StateChooseAttack struct {
	id          int
	w           Worlder
	ai          *Ai
	timer       float64
	nonseelimit float64
	lastpos     pixel.Vec
	isbusy      bool
	skills      []*actor.Skill
	skills2     []*actor.Skill

	afterTimer  float64 // таймер после атаки
	cooldownDur float64 // длительность паузы
}

func NewChooseAttack(ai *Ai, w Worlder) *StateChooseAttack {
	cha := &StateChooseAttack{
		id:          CHOOSEATTACK,
		ai:          ai,
		w:           w,
		nonseelimit: 1,
	}

	skills := ai.obj.GetSkills()
	cha.skills = make([]*actor.Skill, 0, len(skills)) // скилы по возрастанию макс дистанции
	cha.skills2 = make([]*actor.Skill, 0, len(skills))

	for _, sk := range skills {
		cha.skills = append(cha.skills, sk)
		cha.skills2 = append(cha.skills2, sk)
	}
	sort.Slice(cha.skills, func(i, j int) bool {
		return cha.skills[i].Max < cha.skills[j].Max
	})
	sort.Slice(cha.skills2, func(i, j int) bool {
		return cha.skills2[i].Min < cha.skills2[j].Min
	})

	return cha
}

func (s *StateChooseAttack) Update(dt float64) {
	if s.isbusy {
		return
	}

	if s.afterTimer > 0 {
		s.afterTimer -= dt

		return
	}

	hero := s.ai.obj.GetEnemy()
	if hero == nil {
		return
	}

	heropos := hero.GetPos()
	pos := s.ai.obj.GetPos()
	dir := s.ai.obj.GetDir()

	// Проверяем видимость только если смотрим в сторону героя
	if (heropos.X < pos.X && dir < 0) || (heropos.X > pos.X && dir > 0) {
		isSee := s.w.IsSee(pos, heropos)
		if !isSee {
			s.timer += dt
			if s.timer > s.nonseelimit {
				s.ai.SetState(INVESTIGATE, s.lastpos)
			}
			return
		}
		s.lastpos = heropos
		s.timer = 0
	}

	l := pos.X - heropos.X
	currDist := math.Abs(l)
	maxRange := s.skills[len(s.skills)-1].Max
	minRange := s.skills2[0].Min

	// Слишком далеко — преследуем
	if currDist > maxRange {
		s.ai.SetState(INVESTIGATE, heropos)

		return
	}

	// Слишком близко — атакуем ближайшим скиллом
	if currDist < minRange {
		s.ai.attackskill = s.skills2[0]
		s.ai.obj.SetSkill(s.ai.attackskill)
		s.ai.SetState(ATTACK, heropos)

		return
	}

	// В диапазоне — ищем подходящий скилл
	var choosed *actor.Skill
	for _, skill := range s.skills {
		skillDir := float64(skill.GetDir())
		if skillDir*l > 0 {
			continue // скилл направлен в другую сторону
		}
		if currDist <= skill.Max && currDist >= skill.Min {
			if choosed != nil {
				// Выбираем по весу
				var w1, w2 int
				var minW, maxW *actor.Skill
				if choosed.Weight > skill.Weight {
					w1, w2 = choosed.Weight, skill.Weight
					minW, maxW = skill, choosed
				} else {
					w1, w2 = skill.Weight, choosed.Weight
					minW, maxW = choosed, skill
				}
				dice := rand.Float64() * float64(w1+w2)
				if dice > float64(w1) {
					choosed = minW
				} else {
					choosed = maxW
				}
			} else {
				choosed = skill
			}
		}
	}

	if choosed != nil {
		// Нашли скилл — атакуем
		s.ai.attackskill = choosed
		s.ai.obj.SetSkill(s.ai.attackskill)
		s.ai.SetState(ATTACK, heropos)

		return
	}

	// Скилл не найден в диапазоне — идём ближе к герою
	// (дистанция между min и max скиллов, но ни один не подходит по направлению)
	moveDir := -1.0
	if heropos.X > pos.X {
		moveDir = 1.0
	}
	v := pixel.V(moveDir, 0)
	groundrate := s.ai.obj.StepPrediction(events.WALK, v)
	if groundrate > 0.5 {
		b := bindings.Active.GetBinding(bindings.KeyAction["right"])
		if moveDir < 0 {
			b = bindings.Active.GetBinding(bindings.KeyAction["left"])
		}
		s.ai.obj.KeyAction(b)
	}
}

func (s *StateChooseAttack) Start(poi pixel.Vec) {
	s.lastpos = poi
}

func (s *StateChooseAttack) EventAction(e int) {
	if e == events.BUSY {
		s.isbusy = true
	}
	if e == events.RELEASED {
		s.isbusy = false
	}
}

func (s *StateChooseAttack) IsAlerted() bool {
	return true
}
