package dialogs

import "platformer/actor"

func runAction(code int, d *Dialog, a *actor.Actor) {
	switch code {
	case 1:
		actionSetInteraction(a)
	case 2:
		actionMakeHostile(a)
	case 3:
		actionPayCoin()
	case 4:
		actionVictory(d)
	}
}
