package dialogs

import (
	"platformer/actor"
	"platformer/talks"
)

func actionSetInteraction(a *actor.Actor) {
	// we need: phrasesClass (string)
	phrasesClass := "grandpa"
	a.SetOnInteract(func(a *actor.Actor) { talks.AddPhrase(a.GetRect().Min, phrasesClass) })
}
