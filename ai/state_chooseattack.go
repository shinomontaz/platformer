package ai

import (
	"platformer/actor"
	"platformer/common"
	"platformer/creatures"
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
}

func NewChooseAttack(ai *Ai, w Worlder) *StateChooseAttack {
	return &StateChooseAttack{
		id:          CHOOSEATTACK,
		ai:          ai,
		w:           w,
		nonseelimit: 1,
	}
}

func (s *StateChooseAttack) Update(dt float64) {
	if s.isbusy {
		return
	}

	hero := creatures.GetHero()

	if !hero.IsGround() {
		return
	}

	heropos := hero.GetPos()
	pos := s.ai.obj.GetPos()
	dir := s.ai.obj.GetDir()
	var isSee bool
	if (heropos.X < pos.X && dir < 0) || (heropos.X > pos.X && dir > 0) {
		isSee = s.w.IsSee(pos, heropos)
		if !isSee {
			s.timer += dt
			if s.timer > s.nonseelimit {
				s.ai.SetState(INVESTIGATE, s.lastpos)
			}
			return // if not see do nothing, or
		} else {
			s.lastpos = heropos
			s.timer = 0
		}
	}

	skills := s.ai.obj.GetSkills()
	sort.Slice(skills, func(i, j int) bool {
		return skills[i].Max < skills[j].Max
	})
	skills2 := s.ai.obj.GetSkills()
	sort.Slice(skills2, func(i, j int) bool {
		return skills2[i].Min < skills2[j].Min
	})

	l := pixel.L(pos, heropos)
	currDist := l.Len()
	var choosed *actor.Skill

	if currDist > skills[len(skills)-1].Max {
		choosed = skills[len(skills)-1]
	} else if currDist < skills2[0].Min {
		choosed = skills2[0]
	} else {
		for _, skill := range skills {
			if currDist < skill.Max && currDist > skill.Min {
				if choosed != nil { // we already choosed some appropriate skill
					var minWeightSkill *actor.Skill
					var maxWeightSkill *actor.Skill
					var w1, w2 int
					if choosed.Weight > skill.Weight {
						w1 = choosed.Weight
						w2 = skill.Weight
						minWeightSkill = skill
						maxWeightSkill = choosed
					} else {
						w1 = skill.Weight
						w2 = choosed.Weight
						minWeightSkill = choosed
						maxWeightSkill = skill
					}
					dice := common.GetRandFloat() * float64(w1+w2)
					if dice > float64(w1) {
						choosed = minWeightSkill
					} else {
						choosed = maxWeightSkill
					}
				} else {
					choosed = skill
				}
			}
		}
	}

	s.ai.attackskill = choosed
	s.ai.SetState(ATTACK, heropos)

}

func (s *StateChooseAttack) Start(poi pixel.Vec) {
	s.lastpos = poi
}

func (s *StateChooseAttack) Listen(e int, v pixel.Vec) {
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
