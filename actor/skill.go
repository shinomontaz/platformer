package actor

import (
	"platformer/config"
	"platformer/events"

	"github.com/shinomontaz/pixel"
)

type Skill struct {
	Min    float64
	Max    float64
	Weight int
	Type   string
	Name   string
	Event  int
	Ttl    float64
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
	}

	switch pr.Type {
	case "melee":
		s.Event = events.CTRL
	case "spell":
		s.Event = events.CAST
	case "ranged":
		s.Event = events.RANGED
	}

	if pr.Ttl > 0 {
		s.Ttl = pr.Ttl
	}

	if len(pr.Hitbox) > 0 {
		s.Hitbox = pixel.R(pr.Hitbox[0], pr.Hitbox[1], pr.Hitbox[2], pr.Hitbox[3])
	}

	return s
}
