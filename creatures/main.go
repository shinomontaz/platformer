package creatures

import (
	"github.com/shinomontaz/pixel"

	"platformer/activities"
	"platformer/actor"
	"platformer/common"
)

var npcs []common.Actorer
var enemies []common.Actorer
var hero common.Actorer

func reset_lists() {
	enemies = make([]common.Actorer, 0)
	npcs = make([]common.Actorer, 0)
}

func Init() {
	if len(enemies) > 0 {
		reset_lists()
	}
}

func AddEnemy(enemy common.Actorer) {
	enemies = append(enemies, enemy)
}

func AddNpc(npc *actor.Actor) {
	npcs = append(npcs, npc)
}

func GetHero() common.Actorer {
	return hero
}

func SetHero(h common.Actorer) {
	hero = h
}

func Update(dt float64, visiblePhys, visibleSpec []common.Objecter) {

	if hero != nil {
		hero.Update(dt, visiblePhys)
		hero.UpdateSpecial(dt, visibleSpec)
	}

	for _, en := range enemies {
		en.Update(dt, visiblePhys)
	}
	for _, npc := range npcs {
		npc.Update(dt, visiblePhys)
		npc.UpdateSpecial(dt, visibleSpec)
	}
	activities.UpdateStrikes(dt, enemies, hero)
	// TODO:
	activities.UpdateInteractions(dt, npcs)
}

func List() []common.Actorer {
	return append(npcs, enemies...)
}

func Draw(t pixel.Target) {
	if hero != nil {
		hero.Draw(t)
	}

	for _, e := range enemies {
		e.Draw(t)
	}
	for _, n := range npcs {
		n.Draw(t)
	}
}
