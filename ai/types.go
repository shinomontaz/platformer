package ai

import "github.com/faiface/pixel"

type Stater interface {
	Update(dt float64)
	Start()
	GetVec() pixel.Vec
}

type Subscriber interface {
	GetId() int
	Notify(e int, v pixel.Vec)
}

type Worlder interface {
	GetHero() pixel.Vec
	IsSee(from, to pixel.Vec) bool
	AddAlert(place pixel.Vec, raduis float64)
}

type Manageder interface {
	GetPos() pixel.Vec
	GetDir() int
	Notify(e int, v pixel.Vec)
}
