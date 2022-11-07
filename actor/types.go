package actor

import (
	"platformer/common"

	"github.com/shinomontaz/pixel"
)

type Animater interface {
	//	GetAnims() map[string]Animation
	//	GetSprite(name string, idx int) *pixel.Sprite
	GetSprite(name string, idx int) (pixel.Picture, pixel.Rect)
	GetGroupSprite(group, name string, idx int) (pixel.Picture, pixel.Rect)
	GetGroupLen(name string) int
	GetLen(name string) int
}

type Animation interface {
	GetSprite(idx int) *pixel.Sprite
}

type Worlder interface {
	GetGravity() float64
	AddStrike(owner *Actor, r pixel.Rect, power int, speed pixel.Vec)
	AddInteraction(interactor *Actor)
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
