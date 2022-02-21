package ai

import (
	"platformer/actor"

	"github.com/faiface/pixel"
)

type Stater interface {
	Update(dt float64)
	Start(poi pixel.Vec)
	IsAlerted() bool
	Listen(e int, v pixel.Vec)
}

// type Subscriber interface {
// 	GetId() int
// 	Listen(e int, v pixel.Vec)
// }

type Worlder interface {
	GetHero() *actor.Actor
	IsSee(from, to pixel.Vec) bool
	AddAlert(place pixel.Vec, raduis float64)
}

type Alerter interface {
	GetRect() pixel.Rect
}
