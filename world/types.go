package world

import (
	"platformer/common/quadtree"

	tmx "github.com/lafriks/go-tiled"
	"github.com/shinomontaz/pixel"
)

type Background interface {
	Update(dt float64, pos pixel.Vec)
	Draw(t pixel.Target)
}

type tile struct {
	Index   int
	ID      uint32
	R       pixel.Rect
	Tileset *tmx.Tileset
}

func (t tile) Rect() pixel.Rect {
	return t.R
}

func (t tile) Id() int {
	return t.Index
}

type object struct {
	ID     int
	R      pixel.Rect
	Tile   *tmx.LayerTile
	Center pixel.Vec
	T      int
}

func (o object) Rect() pixel.Rect {
	return o.R
}

func (o object) Id() int {
	return o.ID
}

func (o object) Type() int {
	return o.T
}

type groupType int

const (
	GEOM      groupType = 1
	CREATURES groupType = 2
	SCENERY   groupType = 3
	WATER     groupType = 4
	LOOT      groupType = 5
)

// resolveIndices converts quadtree IDs to object slice
func resolvePhysObjs(indices []int, objs []object) []object {
	result := make([]object, 0, len(indices))
	for _, idx := range indices {
		if idx >= 0 && idx < len(objs) {
			result = append(result, objs[idx])
		}
	}
	return result
}

// ensure object satisfies quadtree.Objecter
var _ quadtree.Objecter = tile{}
var _ quadtree.Objecter = object{}
