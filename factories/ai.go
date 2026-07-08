package factories

import (
	"platformer/actor"
	"platformer/ai"
)

type AiOption func(*ai.Ai)

func NewAi(aiType string, obj *actor.Actor, w actor.Worlder) *ai.Ai {
	switch aiType {
	case "npc_roaming":
		return ai.NewActiveNpc(obj, w)
	case "enemy_roaming":
		return ai.NewActiveEnemy(obj, w)
	case "enemy_resurrected":
		return ai.NewActiveEnemy(obj, w)
	case "enemy_agressive":
		return ai.NewAgressiveEnemy(obj, w)
	case "enemy_patrol":
		return ai.NewPatrolEnemy(obj, w, 150)
	case "npc_fishing":
		return ai.NewFishingNpc(obj, w)
	case "npc_swimming":
		return ai.NewSwimmingNpc(obj, w)
	default:
		return ai.NewCalmEnemy(obj, w)
	}
}
