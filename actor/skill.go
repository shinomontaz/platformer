package actor

import (
	"platformer/config"
	"platformer/events"
)

type Skill struct {
	Min    float64
	Max    float64
	Weight int
	Type   string
	Name   string
	Event  int
}

func NewSkill(pr config.Skill) *Skill {
	var e int
	switch pr.Type {
	case "melee":
		e = events.CTRL
	case "spell":
		e = events.CAST
	}

	return &Skill{
		Min:    pr.Min,
		Max:    pr.Max,
		Weight: pr.Weight,
		Type:   pr.Type,
		Name:   pr.Name,
		Event:  e,
	}
}
