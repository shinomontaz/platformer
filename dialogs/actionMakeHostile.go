package dialogs

import (
	"platformer/actor"
	"platformer/config"
	"platformer/factories"
	"platformer/loot"
)

func actionMakeHostile(a *actor.Actor) {
	// delete from npcs
	crtrs := w.GetCreatures()
	crtrs.DeleteNpc(a)
	//we need: profile to switch (string), world reference
	enemy := factories.NewActor(config.Profiles["bigbloated"], w)
	enemy.Move(a.GetPos())
	enemy.SetOnKill(loot.AddKey)
	//	enemy.SetState(actor.) resurrect!

	ai_type := "enemy_agressive"
	//	ai_type := "enemy_resurrected"
	factories.NewAi(ai_type, enemy, w)

	dir := float64(a.GetDir())
	if dir != 0 {
		enemy.SetDir(dir)
	}

	crtrs.AddEnemy(enemy)
}
