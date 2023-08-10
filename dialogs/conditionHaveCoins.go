package dialogs

import (
	"platformer/actor"
	"platformer/inventory"
)

func conditionHaveCoins(a *actor.Actor) bool {
	return inventory.HaveCoins()
}
