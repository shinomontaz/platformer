package actor

import (
	"platformer/common"

	"github.com/shinomontaz/pixel"
)

type Worlder interface {
	GetGravity() float64
	// AddStrike(owner *Actor, r pixel.Rect, power int, speed pixel.Vec)
	AddSpell(owner *Actor, t pixel.Vec, spell string, objs []common.Objecter)
}

type Stater interface {
	GetId() int
	Start()
	Update(dt float64)
	Listen(e int, v *pixel.Vec)
	GetSprite() *pixel.Sprite
	Busy() bool
}

type soundeffect struct {
	List []string
}

type OnKillHandler func(pos, vel pixel.Vec)
type OnInteractHandler func()
