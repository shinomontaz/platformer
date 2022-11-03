package factories

import (
	"platformer/actor"
	"platformer/ai"
)

func NewAi(t string, obj *actor.Actor, w Worlder) {
	switch t {
	case "oldman":
		ai.NewOldman(obj, w)
	default:
		ai.NewCommon(obj, w)
	}
}
