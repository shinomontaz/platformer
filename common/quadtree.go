package common

import (
	"fmt"

	"github.com/shinomontaz/pixel"
)

const MAX = 10

type Objecter struct {
	R    pixel.Rect
	ID   uint32
	Type int
}

func (o *Objecter) Rect() pixel.Rect {
	return o.R
}

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
				if !qt.nodes[idx].Insert(o) {
					panic(fmt.Sprintf("cannot insert %+v into %+v", o.Rect(), qt.nodes[idx]))
				}
			} else {
				qt.items[i] = o
				i++
			}
		}
		qt.items = qt.items[:i]
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

	verticalMidpoint := qt.b.Min.X + (qt.b.W() / 2)
	horizontalMidpoint := qt.b.Min.Y + (qt.b.H() / 2)

	bottom := (r.Min.Y < horizontalMidpoint) && (r.Max.Y < horizontalMidpoint)
	top := r.Min.Y > horizontalMidpoint

	if r.Min.X < verticalMidpoint && r.Max.X < verticalMidpoint {
		if top {
			idx = 1
		}
		if bottom {
			idx = 2
		}
	} else if r.Min.X > verticalMidpoint {
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
		qt.b.Min.Y+(qt.b.H()/2.0),
		qt.b.Max.X,
		qt.b.Max.Y,
	)))

	qt.nodes = append(qt.nodes, New(qt.level+1, pixel.R( // 1
		qt.b.Min.X,
		qt.b.Min.Y+(qt.b.H()/2.0),
		qt.b.Min.X+(qt.b.W()/2.0),
		qt.b.Max.Y,
	)))

	qt.nodes = append(qt.nodes, New(qt.level+1, pixel.R( // 2
		qt.b.Min.X,
		qt.b.Min.Y,
		qt.b.Min.X+(qt.b.W()/2.0),
		qt.b.Min.Y+(qt.b.H()/2.0),
	)))

	qt.nodes = append(qt.nodes, New(qt.level+1, pixel.R( // 3
		qt.b.Min.X+(qt.b.W()/2.0),
		qt.b.Min.Y,
		qt.b.Max.X,
		qt.b.Min.Y+(qt.b.H()/2.0),
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
