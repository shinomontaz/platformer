package world

import (
	"platformer/common"
	"platformer/common/quadtree"
	"platformer/loot"
	"platformer/objects"

	tmx "github.com/lafriks/go-tiled"
	"github.com/shinomontaz/pixel"
)

func (w *World) loadMap(name string) (*tmx.Map, error) {
	mapsource, err := w.loader.Open(name)
	if err != nil {
		return nil, err
	}
	defer mapsource.Close()

	return tmx.LoadReader("assets", mapsource)
}

func (w *World) readObjects(tm *tmx.Map) map[groupType]*tmx.ObjectGroup {
	objMap := make(map[groupType]*tmx.ObjectGroup)
	for _, og := range tm.ObjectGroups {
		og := og
		switch og.Name {
		case "geom":
			objMap[GEOM] = og
		case "meta":
			objMap[CREATURES] = og
		case "scenery":
			objMap[SCENERY] = og
		case "water":
			objMap[WATER] = og
		case "loot":
			objMap[LOOT] = og
		}
	}

	return objMap
}

func (w *World) initPhys() {
	w.qtPhys = quadtree.New(1, 10, w.r)
	og := w.objMap[GEOM]
	if og == nil {
		return
	}
	i := 0
	for _, o := range og.Objects {
		min := pixel.V(
			float64(o.X),
			float64(w.r.H())-float64(o.Y)-float64(o.Height),
		)
		max := pixel.Vec{
			X: min.X + float64(o.Width),
			Y: min.Y + float64(o.Height),
		}

		rc := pixel.Rect{
			Min: min,
			Max: max,
		}
		phType := common.GROUND
		if o.Class == "barier" {
			phType = common.BARRIER
		}
		oo := object{
			ID:     i,
			R:      rc,
			Center: rc.Center(),
			T:      phType,
		}
		w.physobjs = append(w.physobjs, oo)
		w.qtPhys.Insert(oo)
		w.drawRectIfNeeded(o, rc)
		i++
	}
}

/*
crtrs := creatures.New()

		list := g.w.GetMetas()
		for _, o := range list {
			if o.Class == "enemy" {
				enemy := factories.NewActor(config.Profiles[o.Name])
				enemy.SetEnemy(g.hero)
				enemy.Move(pixel.V(o.X, o.Y))
				rew := o.Properties.GetString("reward")
				if rew == "coin" { // make OnKill handler
					enemy.SetOnKill(loot.AddCoin)
				} else if rew == "key" {
					enemy.SetOnKill(loot.AddKey)
				}
				ai_type := o.Properties.GetString("ai")
				if ai_type != "" {
					g.ailist = append(g.ailist, factories.NewAi(ai_type, enemy, w))
				}

				dir := o.Properties.GetFloat("dir")
				if dir != 0 {
					enemy.SetDir(dir)
				}

				enemy.SetOnKill(func(pos, vel pixel.Vec) { score.Add(config.Profiles[o.Name].Hp) })
				//			crtrs.AddEnemy(enemy)
			}
			if o.Class == "npc" {
				npc := factories.NewActor(config.Profiles[o.Name])
				npc.Move(pixel.V(o.X, o.Y))
				rew := o.Properties.GetString("reward")
				if rew == "coin" {
					// make OnKill handler
					npc.SetOnKill(loot.AddCoin)
					// rew_count := o.Properties.GetInt("reward_count")
					// if rew_count > 0 {
					// }
				} else if rew == "key" {
					npc.SetOnKill(loot.AddKey)
				}

				interact_type := o.Properties.GetString("interact")
				if interact_type == "phrase" {
					phrasesClass := o.Properties.GetString("phrasesClass")
					// create interact handler
					npc.SetOnInteract(func(a *actor.Actor) { talks.AddPhrase(a.GetRect().Min, phrasesClass) })
				} else if interact_type == "die" {
					npc.SetOnInteract(func(a *actor.Actor) { a.Kill() })
				} else if interact_type == "dialog" {
					dialog_id := o.Properties.GetInt("dialog")
					npc.SetOnInteract(func(a *actor.Actor) {
						fmt.Println("DIALOG!!!")
						dialogs.SetActive(dialog_id, a)
						g.SetState(gamestate.DIALOG)
					})
				}
				ai_type := o.Properties.GetString("ai")
				if ai_type != "" {
					g.ailist = append(g.ailist, factories.NewAi(ai_type, npc, w))
				}

				dir := o.Properties.GetFloat("dir")
				if dir != 0 {
					npc.SetDir(dir)
				}

				crtrs.AddNpc(npc)
			}
		}
	//	crtrs.SetHero(g.hero)
*/
func (w *World) initLoot() {
	og := w.objMap[LOOT]
	if og == nil {
		return
	}
	for _, o := range og.Objects {
		if o.Class == "coin" {
			loot.AddCoin(pixel.V(o.X, o.Y), pixel.ZV)
		}
		if o.Class == "object" {
			ai_type := o.Properties.GetString("ai")
			objects.Add(o.Name, pixel.R(o.X, o.Y-o.Height, o.X+o.Width, o.Y), ai_type)
		}
	}
}

func (w *World) initWater() {
	og := w.objMap[WATER]
	if og == nil {
		return
	}
	for _, o := range og.Objects {
		min := pixel.V(
			float64(o.X),
			float64(w.Height)-float64(o.Y)-float64(o.Height),
		)
		max := pixel.Vec{
			X: min.X + float64(o.Width),
			Y: min.Y + float64(o.Height),
		}

		rc := pixel.Rect{
			Min: min,
			Max: max,
		}

		oo := object{
			ID:     len(w.specobjs),
			R:      rc,
			Center: rc.Center(),
			T:      common.WATER,
		}

		w.specobjs = append(w.specobjs, oo)
		w.qtSpec.Insert(oo)

		w.drawRectIfNeeded(o, rc)
	}
}
