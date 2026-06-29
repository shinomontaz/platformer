package actor

import (
	"platformer/common"

	"github.com/shinomontaz/pixel"
)

type Worlder interface {
	GetGravity() float64
	IsSee(from, to pixel.Vec) bool
	AddSpell(owner *Actor, t pixel.Vec, spell string, objs []common.O)
	GetCenter() pixel.Vec
}

type Stater interface {
	GetId() int
	Start()
	Update(dt float64)
	Listen(e int, v *pixel.Vec)
	GetSprite() *pixel.Sprite
	Busy() bool
	SetWater(bool)
	SetWaterResistant(bool)
}

type soundeffect struct {
	List []string
}

type OnKillHandler func(pos, vel pixel.Vec)
type OnInteractHandler func(a *Actor)

type object interface {
	Rect() pixel.Rect
	Id() int
}
