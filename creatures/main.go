package creatures

import (
	"github.com/shinomontaz/pixel"

	"platformer/activities"
	"platformer/actor"
	"platformer/common"
)

var npcs []*actor.Actor
var enemies []*actor.Actor
var hero *actor.Actor

func init() {
	reset_lists()
}

func reset_lists() {
	enemies = make([]*actor.Actor, 0)
	npcs = make([]*actor.Actor, 0)
}

func Init(w Worlder, h *actor.Actor) {
	if len(enemies) > 0 {
		reset_lists()
	}

	// list := w.GetMetas()
	// for _, o := range list {
	// 	if o.Class == "enemy" {
	// 		addEnemy(o, w)
	// 	}
	// 	if o.Class == "npc" {
	// 		addNpc(o, w)
	// 	}
	// }

	hero = h
}

func AddEnemy(enemy *actor.Actor) {
	// enemy := factories.NewActor(config.Profiles[meta.Name], w)
	// enemy.Move(pixel.V(meta.X, meta.Y))
	// factories.NewAi(config.Profiles[meta.Name].Type, enemy, w)
	enemies = append(enemies, enemy)
}

//func AddNpc(meta *tmx.Object, w Worlder) {

func AddNpc(npc *actor.Actor) {
	// npc := factories.NewActor(config.Profiles[meta.Name], w)
	// npc.Move(pixel.V(meta.X, meta.Y))
	// factories.NewAi(config.Profiles[meta.Name].Type, npc, w)
	npcs = append(npcs, npc)
}

func GetHero() *actor.Actor {
	return hero
}

func SetHero(h *actor.Actor) {
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
	activities.UpdateStrikes(dt, hero)
	//	updateSpells(dt, enemies, hero)
}

//func (w *World) AddInteraction(interactor Actor) {
//	AddStrike(owner, r, power, speed)
//}

func List() []*actor.Actor {
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
