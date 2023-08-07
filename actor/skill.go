package actor

import (
	"platformer/bindings"
	"platformer/config"
	"platformer/events"

	"github.com/shinomontaz/pixel"
)

var ActionSkillMap = map[string]int{
	"melee":     events.MELEE,
	"meleemove": events.MELEEMOVE,
	"ranged":    events.RANGED,
	"cast":      events.CAST,
}

type Skill struct {
	Id     int
	Min    float64
	Max    float64
	Weight int
	Type   string
	Name   string
	Keys   []int
	Event  int
	Ttl    float64
	Speed  float64
	Hitbox pixel.Rect
}

func NewSkill(pr config.Skill) *Skill {
	s := &Skill{
		Min:    pr.Min,
		Max:    pr.Max,
		Weight: pr.Weight,
		Type:   pr.Type,
		Name:   pr.Name,
		Ttl:    0.5,
		Speed:  pr.Speed,
		Keys:   make([]int, 0),
	}

	s.Event = ActionSkillMap[pr.Name]

	for _, ev := range pr.Keys {
		if eventId, ok := bindings.KeyAction[ev]; ok {
			s.Keys = append(s.Keys, eventId)
		}
	}

	if pr.Ttl > 0 {
		s.Ttl = pr.Ttl
	}

	if len(pr.Hitbox) > 0 {
		s.Hitbox = pixel.R(pr.Hitbox[0], pr.Hitbox[1], pr.Hitbox[2], pr.Hitbox[3])
	}

	return s
}
