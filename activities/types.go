package activities

import (
	"platformer/common"

	"github.com/shinomontaz/pixel"
)

type Actor interface {
	Listen(e int, v pixel.Vec)
	Draw(t pixel.Target)
	Update(dt float64, visiblePhys []common.Objecter)
	UpdateSpecial(dt float64, visibleSpec []common.Objecter)
	GetRect() pixel.Rect
	GetPos() pixel.Vec
	GetId() int
	Hit(vec pixel.Vec, power int)
}
