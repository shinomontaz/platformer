package menu

import (
	"github.com/faiface/pixel"
)

type Menuer interface {
	GetId() int
	Start()
	Update(dt float64)
	Draw(t pixel.Target)

	Listen(e int, v pixel.Vec)
}
