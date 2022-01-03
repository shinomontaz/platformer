package common

import (
	"github.com/faiface/pixel"
	"github.com/faiface/pixel/imdraw"
)

type Quadtree struct {
	list []Objecter
}

func NewQuadtree() *Quadtree {
	return &Quadtree{}
}

func (qt *Quadtree) Add(o Objecter) {
	qt.list = append(qt.list, o)
}

func (qt *Quadtree) GetObjects() []Objecter {
	return qt.list
}

func (qt *Quadtree) Clear() {

}

func (qt *Quadtree) CanIntersect(rect pixel.Rect) []Objecter {
	// TODO: make
	return qt.list
}

type Objecter interface {
	Draw(imd *imdraw.IMDraw)
	Rect() *pixel.Rect
	Pixels() []uint32
	Hit(pos, vel pixel.Vec, power int) // hit coords, hit velocity, hit strength
	// Name() string
	// Type() int
}
