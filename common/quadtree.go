package common

import (
	"fmt"

	"github.com/faiface/pixel"
	"github.com/faiface/pixel/imdraw"
)

const MAX = 10

/*
________
|   |   |
| 1 | 0 |
|___|___|
|   |   |
| 2 | 3 |
|___|___|
*/
type Quadtree struct {
	items []Objecter
	nodes []*Quadtree
	level int
	b     pixel.Rect
}

func New(level int, b pixel.Rect) *Quadtree {
	return &Quadtree{
		level: level,
		b:     b,
		nodes: make([]*Quadtree, 0),
		items: make([]Objecter, 0),
	}
}

func (qt *Quadtree) Insert(o Objecter) bool {
	if !qt.b.Intersects(o.Rect()) {
		fmt.Println("!intersects!", o.Rect(), qt.b)
		return false
	}

	if len(qt.nodes) == 0 {
		if len(qt.items) < MAX {
			qt.items = append(qt.items, o)
			return true
		}

		qt.subdivide()
		// try to move all items by subnodes
		i := 0
		for _, o := range qt.items {
			idx := qt.index(o.Rect())
			if idx > -1 {
				qt.nodes[idx].Insert(o)
			} else {
				qt.items[i] = o
				i++
			}
		}
		qt.items = qt.items[0:i]
	}
	idx := qt.index(o.Rect())
	if idx > -1 {
		qt.nodes[idx].Insert(o)
		return true
	}

	qt.items = append(qt.items, o)
	return true
}

func (qt *Quadtree) index(r pixel.Rect) int { // get quadrand index if object contains in any and -1 otherwise
	idx := -1

	top := (r.Max.Y < (qt.b.Min.Y + qt.b.H()/2.0)) && (r.Min.Y < (qt.b.Min.Y + qt.b.H()/2.0))
	bottom := r.Max.Y > (qt.b.Min.Y + qt.b.H()/2.0)

	if r.Max.X < qt.b.Min.X+(qt.b.W()/2.0) {
		if top {
			idx = 1
		}
		if bottom {
			idx = 2
		}
	} else if r.Max.X > qt.b.Min.X+(qt.b.W()/2.0) {
		if top {
			idx = 0
		}
		if bottom {
			idx = 3
		}
	}

	return idx
}

func (qt *Quadtree) subdivide() {
	qt.nodes = append(qt.nodes, New(qt.level+1, pixel.R( // 0
		qt.b.Min.X+(qt.b.W()/2.0),
		qt.b.Max.Y,
		qt.b.Min.X,
		qt.b.Min.Y+(qt.b.H()/2.0),
	)))

	qt.nodes = append(qt.nodes, New(qt.level+1, pixel.R( // 1
		qt.b.Max.X,
		qt.b.Max.Y,
		qt.b.Min.X+(qt.b.W()/2.0),
		qt.b.Min.Y+(qt.b.H()/2.0),
	)))

	qt.nodes = append(qt.nodes, New(qt.level+1, pixel.R( // 2
		qt.b.Max.X,
		qt.b.Min.Y+(qt.b.H()/2.0),
		qt.b.Min.X+(qt.b.W()/2.0),
		qt.b.Min.Y,
	)))

	qt.nodes = append(qt.nodes, New(qt.level+1, pixel.R( // 3
		qt.b.Min.X+(qt.b.H()/2.0),
		qt.b.Min.X+(qt.b.W()/2.0),
		qt.b.Min.X,
		qt.b.Min.Y,
	)))
}

func (qt *Quadtree) Retrieve(r pixel.Rect) []Objecter {
	res := qt.items
	idx := qt.index(r)
	if len(qt.nodes) > 0 {
		if idx != -1 {
			res = append(res, qt.nodes[idx].Retrieve(r)...)
		} else {
			for _, n := range qt.nodes {
				res = append(res, n.Retrieve(r)...)
			}
		}
	}
	return res
}

func (qt *Quadtree) Print() string {
	str := ""
	str = fmt.Sprintf("level: %d ( %d )", qt.level, len(qt.items))
	for _, n := range qt.nodes {
		str = fmt.Sprintf("%s\n %s", str, n.Print())
	}

	return str
}

func (qt *Quadtree) Clear() {
	qt.items = []Objecter{}
	qt.level = 0
	for _, n := range qt.nodes {
		n.Clear()
	}
}

type Objecter interface {
	Draw(imd *imdraw.IMDraw)
	Rect() pixel.Rect
	Pixels() []uint32
	Hit(pos, vel pixel.Vec, power int) // hit coords, hit velocity, hit strength
	// Name() string
	// Type() int
}
