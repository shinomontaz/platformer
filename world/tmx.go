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
