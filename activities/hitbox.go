package activities

import (
	"platformer/common"

	"github.com/shinomontaz/pixel"
	"github.com/shinomontaz/pixel/imdraw"
	"golang.org/x/image/colornames"
)

type HitBox struct {
	rect   pixel.Rect
	ttl    float64
	timer  float64
	power  int
	owner  common.Actorer
	hitted map[common.Actorer]struct{}
	dir    pixel.Line
	speed  pixel.Vec
}

var enboxes []*HitBox
var plboxes []*HitBox

func init() {
	enboxes = make([]*HitBox, 0)
	plboxes = make([]*HitBox, 0)
}

func AddStrike(owner common.Actorer, rect pixel.Rect, power int, speed pixel.Vec, opts ...HitBoxOption) *HitBox {
	center := rect.Center()
	from := pixel.V(rect.Min.X, center.Y)
	to := pixel.V(rect.Max.X, center.Y)

	if owner.GetPos().X > center.X {
		to = pixel.V(rect.Min.X, center.Y)
		from = pixel.V(rect.Max.X, center.Y)
	}
	b := &HitBox{
		rect:   rect,
		ttl:    0.2,
		power:  power,
		owner:  owner,
		dir:    pixel.L(from, to),
		hitted: make(map[common.Actorer]struct{}),
		speed:  speed,
	}

	for _, o := range opts {
		o(b)
	}

	if owner.GetId() == 1 {
		plboxes = append(plboxes, b)
	} else {
		enboxes = append(enboxes, b)
	}

	return b
}

func UpdateStrikes(dt float64, enemies []common.Actorer, player common.Actorer) {
	updatePlStrikes(dt, enemies)
	updateEnStrikes(dt, []common.Actorer{player})
}

func updatePlStrikes(dt float64, hittable []common.Actorer) {
	i := 0
	for _, b := range plboxes {
		for _, hh := range hittable {
			if hh == b.owner {
				continue
			}
			if _, ok := b.hitted[hh]; ok {
				continue
			}
			r := hh.GetRect()
			if b.rect.Intersects(r) {
				vec := pixel.ZV // TODO: detect hit vector
				vec.X = -1
				if r.Center().X > b.dir.A.X {
					vec.X = 1
				}
				hh.Hit(vec, b.power)
				b.hitted[hh] = struct{}{}
			}
		}
		b.timer += dt
		if b.ttl < 0 || b.timer < b.ttl {
			plboxes[i] = b
			i++
		}
	}

	plboxes = plboxes[:i]
}

func updateEnStrikes(dt float64, hittable []common.Actorer) {
	i := 0
	for _, b := range enboxes {
		for _, hh := range hittable {
			if hh == b.owner {
				continue
			}
			if _, ok := b.hitted[hh]; ok {
				continue
			}
			r := hh.GetRect()
			if b.rect.Intersects(r) {
				vec := pixel.ZV // TODO: detect hit vector
				vec.X = -float64(b.power)
				if b.dir.A.X > b.dir.B.X {
					vec.X = float64(b.power)
				}
				hh.Hit(vec, b.power)
				b.hitted[hh] = struct{}{}
			}
		}
		b.timer += dt
		if b.ttl < 0 || b.timer < b.ttl {
			enboxes[i] = b
			i++
		}
	}

	enboxes = enboxes[:i]
}

func DrawStrikes(t pixel.Target) {
	imd := imdraw.New(nil)
	for i := range enboxes {
		vertices := enboxes[i].rect.Vertices()
		imd.Color = colornames.Red
		for _, v := range vertices {
			imd.Push(v)
		}
		imd.Rectangle(1)
	}
	for i := range plboxes {
		vertices := plboxes[i].rect.Vertices()
		imd.Color = colornames.Red
		for _, v := range vertices {
			imd.Push(v)
		}
		imd.Rectangle(1)
	}
	imd.Draw(t)
}

func (hb *HitBox) Move(vec pixel.Vec) {
	hb.rect = hb.rect.Moved(vec)
}
