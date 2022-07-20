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
	Hitbox pixel.Rect
}

func NewSkill(pr config.Skill) *Skill {
	var e int
	switch pr.Type {
	case "melee":
		e = events.CTRL
	case "spell":
		e = events.CAST
	}

	s := &Skill{
		Min:    pr.Min,
		Max:    pr.Max,
		Weight: pr.Weight,
		Type:   pr.Type,
		Name:   pr.Name,
		Event:  e,
	}

	if len(pr.Hitbox) > 0 {
		s.Hitbox = pixel.R(pr.Hitbox[0], pr.Hitbox[1], pr.Hitbox[2], pr.Hitbox[3])
	}

	return s
}
