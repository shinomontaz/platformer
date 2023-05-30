package dialogs

import (
	"platformer/actor"
	"platformer/config"
	"platformer/creatures"
	"platformer/factories"
	"platformer/loot"
)

func actionMakeHostile(a *actor.Actor) {
	// delete from npcs
	creatures.DeleteNpc(a)
	//we need: profile to switch (string), world reference
	enemy := factories.NewActor(config.Profiles["bigbloated"], w)
	enemy.Move(a.GetPos())
	enemy.SetOnKill(loot.AddKey)

	ai_type := "enemy_roaming"
	factories.NewAi(ai_type, enemy, w)

	dir := float64(a.GetDir())
	if dir != 0 {
		enemy.SetDir(dir)
	}

	creatures.AddEnemy(enemy)
}
