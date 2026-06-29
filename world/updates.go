package world

import (
	"platformer/common"
	"platformer/objects"

	"github.com/go-gl/mathgl/mgl32"
)

func (w *World) updateTiles() {
	w.visibleTiles = w.visibleTiles[:0]
	retrieved := w.qtTile.Retrieve(w.viewport)
	for _, idx := range retrieved {
		w.visibleTiles = append(w.visibleTiles, w.tiles[idx])
	}
}

func (w *World) updateObjects() {
	w.visibleObjs = w.visibleObjs[:0]
	retrieved := w.qtObjs.Retrieve(w.viewport)
	for _, idx := range retrieved {
		w.visibleObjs = append(w.visibleObjs, w.objects[idx])
	}
}

func (w *World) updateSpecs() {
	w.visibleSpec = w.visibleSpec[:0]
	for _, idx := range w.qtSpec.Retrieve(w.viewport) {
		w.visibleSpec = append(w.visibleSpec, w.specobjs[idx])
	}
}

func (w *World) updatePhys() {
	w.uObjects = w.uObjects[:0]
	w.visiblePhys = w.visiblePhys[:0]
	c := w.viewport.Center()

	for _, p := range w.qtPhys.Retrieve(w.viewport) {
		o := w.physobjs[p]
		w.visiblePhys = append(w.visiblePhys, o)

		if o.T != common.GROUND { // Для shader
			continue
		}
		uObject := mgl32.Vec4{float32(-c.X + o.R.Min.X), float32(-c.Y + o.R.Min.Y - 150), float32(-c.X + o.R.Max.X), float32(-c.Y + o.R.Max.Y - 150)}
		w.uObjects = append(w.uObjects, uObject)
	}

	objList := objects.List()
	w.visiblePhys = append(w.visiblePhys, objList...)

	w.uNumObjects = int32(len(w.uObjects))
}
