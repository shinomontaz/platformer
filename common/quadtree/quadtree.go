package quadtree

import (
	"fmt"

	"github.com/shinomontaz/pixel"
)

const MAX = 10

type Objecter interface {
	Rect() pixel.Rect
	Id() int
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
type Qt struct {
	ids   []int
	items []Objecter
	nodes []*Qt
	level int
	max   int
	b     pixel.Rect
	vMid  float64
	hMid  float64
}

func New(level, max int, b pixel.Rect) *Qt {
	return &Qt{
		level: level,
		b:     b,
		max:   max,
		ids:   []int{},
		nodes: []*Qt{},
		items: []Objecter{},
		vMid:  b.Min.X + (b.W() / 2),
		hMid:  b.Min.Y + (b.H() / 2),
	}
}

func (qt *Qt) Insert(o Objecter) bool {
	if !qt.b.Intersects(o.Rect()) {
		return false
	}

	if len(qt.items) < qt.max {
		qt.items = append(qt.items, o)
		qt.ids = append(qt.ids, o.Id())
		return true
	}

	if len(qt.nodes) == 0 {
		qt.subdivide()
		// try to move all items by subnodes
		i := 0
		for _, oo := range qt.items {
			idx := qt.index(oo.Rect())
			if idx > -1 {
				if !qt.nodes[idx].Insert(oo) {
					panic(fmt.Sprintf("cannot insert %+v into %+v", oo.Rect(), qt.nodes[idx]))
				}
			} else {
				qt.items[i] = oo
				qt.ids[i] = oo.Id()
				i++
			}
		}
		qt.items = qt.items[:i]
		qt.ids = qt.ids[:i]
	}

	idx := qt.index(o.Rect())
	if idx > -1 {
		qt.nodes[idx].Insert(o)
		return true
	}

	// corder case: point hit rectangle borders
	qt.items = append(qt.items, o)
	qt.ids = append(qt.ids, o.Id())
	return true
}

func (q *Qt) index(r pixel.Rect) int { // get quadrant index, -1 if crosses boundary
	idx := -1

	bottom := (r.Min.Y < q.hMid) && (r.Max.Y < q.hMid)
	top := r.Min.Y > q.hMid

	if r.Min.X < q.vMid && r.Max.X < q.vMid {
		if top {
			idx = 1
		}
		if bottom {
			idx = 2
		}
	} else if r.Min.X > q.vMid {
		if top {
			idx = 0
		}
		if bottom {
			idx = 3
		}
	}

	return idx
}

func (q *Qt) subdivide() {
	q.nodes = []*Qt{
		New(q.level+1, q.max, pixel.R(q.b.Min.X+(q.b.W()/2.0), q.b.Min.Y+(q.b.H()/2.0), q.b.Max.X, q.b.Max.Y)),
		New(q.level+1, q.max, pixel.R(q.b.Min.X, q.b.Min.Y+(q.b.H()/2.0), q.b.Min.X+(q.b.W()/2.0), q.b.Max.Y)),
		New(q.level+1, q.max, pixel.R(q.b.Min.X, q.b.Min.Y, q.b.Min.X+(q.b.W()/2.0), q.b.Min.Y+(q.b.H()/2.0))),
		New(q.level+1, q.max, pixel.R(q.b.Min.X+(q.b.W()/2.0), q.b.Min.Y, q.b.Max.X, q.b.Min.Y+(q.b.H()/2.0))),
	}
}

func (qt *Qt) Retrieve(r pixel.Rect) []int {
	res := make([]int, len(qt.ids))
	copy(res, qt.ids)

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

func (qt *Qt) Print() string {
	str := fmt.Sprintf("level: %d ( %d )", qt.level, len(qt.items))
	for _, n := range qt.nodes {
		str = fmt.Sprintf("%s\n %s", str, n.Print())
	}

	return str
}

func (qt *Qt) Clear() {
	qt.items = qt.items[:0]
	qt.ids = qt.ids[:0]
	qt.level = 0
	for _, n := range qt.nodes {
		n.Clear()
	}
}
