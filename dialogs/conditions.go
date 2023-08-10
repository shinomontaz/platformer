package dialogs

import (
	"platformer/actor"
)

func checkCondition(code int, a *actor.Actor) bool {
	switch code {
	case 1:
		return conditionHaveCoins(a)
	}
	return false
}
