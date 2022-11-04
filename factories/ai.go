package factories

import (
	"platformer/actor"
	"platformer/ai"
)

func NewAi(t string, obj *actor.Actor, w Worlder) {
	switch t {
	case "oldman":
		ai.NewActiveNpc(obj, w)
	default:
		ai.NewCalmEnemy(obj, w)
	}
}
