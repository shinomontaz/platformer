package actor

import (
	"platformer/common"

	"github.com/faiface/pixel"
)

type Animater interface {
	//	GetAnims() map[string]Animation
	//	GetSprite(name string, idx int) *pixel.Sprite
	GetSprite(name string, idx int) (pixel.Picture, pixel.Rect)
	GetGroupSprite(group, name string, idx int) (pixel.Picture, pixel.Rect)
	GetGroupLen(name string) int
}

type Animation interface {
	GetSprite(idx int) *pixel.Sprite
}

type Worlder interface {
	GetGravity() float64
	GetQt() *common.Quadtree // temporary solution, we will check collision in the world ?
	AddStrike(owner *Actor, r pixel.Rect, power float64)
}

type Stater interface {
	GetId() int
	Start()
	Update(dt float64)
	Notify(e int, v *pixel.Vec)
	GetSprite() *pixel.Sprite
	Busy() bool
}
