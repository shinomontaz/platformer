package factories

import (
	"platformer/actor"
	"platformer/ai"
)

func NewAi(t string, obj *actor.Actor, w Worlder) {
	switch t {
	// case "deceased":
	// 	ai.NewMage(obj, w)
	default:
		ai.NewCommon(obj, w)
	}
}
