package main

import (
	"github.com/faiface/pixel"
	"github.com/faiface/pixel/imdraw"
)

type Quadtree interface {
	GetObjects() []Objecter
}

type Objecter interface {
	Draw(imd *imdraw.IMDraw)
	Rect() pixel.Rect
	Hit(pos, vel pixel.Vec, power int) // hit coords, hit velocity, hit strength
	// Name() string
	// Type() int
}
