package actor

import (
	"platformer/common"

	"github.com/faiface/pixel"
)

type Animater interface {
	//	GetAnims() map[string]Animation
	//	GetSprite(name string, idx int) *pixel.Sprite
	GetSprite(name string, idx int) (pixel.Picture, pixel.Rect)
	GetLen(name string) int
}

type Animation interface {
	GetSprite(idx int) *pixel.Sprite
}

type Worlder interface {
	GetGravity() float64
	GetQt() *common.Quadtree // temporary solution, we will check collision in the world ?
}

type ActorStater interface {
	GetId() int
	Start()
	Update(dt float64)
	Notify(e int, v *pixel.Vec)
	GetSprite() *pixel.Sprite
}

type AnimStater interface {
	GetId() int
	Start()
	Update(dt float64)
	Notify(e int)
}
