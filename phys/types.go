package phys

import (
	"github.com/faiface/pixel"
)

type Quadtree interface {
	GetObjects() []Objecter
}

type Objecter interface {
	Rect() *pixel.Rect
	Hit(pos, vel pixel.Vec, power int) // hit coords, hit velocity, hit strength
	// Name() string
	// Type() int
}
