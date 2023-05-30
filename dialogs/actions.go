package dialogs

import "platformer/actor"

func runAction(code int, a *actor.Actor) {
	switch code {
	case 1:
		actionSetInteraction(a)
	case 2:
		actionMakeHostile(a)
	}
}
