package world

import "github.com/shinomontaz/pixel"

type Background interface {
	Update(dt float64, pos pixel.Vec)
	Draw(t pixel.Target)
}
