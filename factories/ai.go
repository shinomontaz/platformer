package factories

import (
	"platformer/actor"
	"platformer/ai"
)

func NewAi(t string, obj *actor.Actor, w Worlder) {
	switch t {
	case "npc_roaming":
		ai.NewActiveNpc(obj, w)
	case "enemy_roaming":
		ai.NewActiveEnemy(obj, w)
	default:
		ai.NewCalmEnemy(obj, w)
	}
}
