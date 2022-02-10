package magic

import (
	"github.com/faiface/pixel"
)

type Actor interface {
	GetPos() pixel.Vec
	GetDir() int
	Notify(e int, v pixel.Vec)
	GetHp() int
}
