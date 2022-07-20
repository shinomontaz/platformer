package magic

import (
	"platformer/common"

	"platformer/actor"

	"github.com/shinomontaz/pixel"
)

type Speller interface {
	IsFinished() bool
	Draw(t pixel.Target)
	Update(dt float64)
}

type Worlder interface {
	GetGravity() float64
	GetQt() *common.Quadtree // temporary solution, we will check collision in the world ?
	AddStrike(a *actor.Actor, rect pixel.Rect, power int, speed pixel.Vec)
}

type Animater interface {
	GetSprite(name string, idx int) (pixel.Picture, pixel.Rect)
	GetLen(name string) int
}

type Animation interface {
	GetSprite(idx int) *pixel.Sprite
}

type Spellprofile interface {
	GetHitbox() pixel.Rect
	GetSound() string
}
