package ai

import (
	"platformer/magic"

	"github.com/faiface/pixel"
)

type Stater interface {
	Update(dt float64)
	Start(poi pixel.Vec)
	IsAlerted() bool
	Notify(e int, v pixel.Vec)
}

type Subscriber interface {
	GetId() int
	Notify(e int, v pixel.Vec)
}

type Worlder interface {
	GetHero() Manageder
	IsSee(from, to pixel.Vec) bool
	AddAlert(place pixel.Vec, raduis float64)
}

type Manageder interface {
	GetPos() pixel.Vec
	GetDir() int
	Notify(e int, v pixel.Vec)
	SetAi(a *Ai)
	GetHp() int
	GetAttackrange() float64
	GetMagic() *magic.Magic
}

type Alerter interface {
	GetRect() pixel.Rect
}
