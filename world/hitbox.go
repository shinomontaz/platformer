package world

import (
	"platformer/actor"
	"platformer/ai"

	"github.com/faiface/pixel"
)

type HitBox struct {
	rect   pixel.Rect
	ttl    float64
	timer  float64
	power  int
	owner  *actor.Actor
	hitted map[*actor.Actor]struct{}
	dir    pixel.Line
}

var enboxes []*HitBox
var plboxes []*HitBox

func init() {
	enboxes = make([]*HitBox, 0)
	plboxes = make([]*HitBox, 0)
}

func AddStrike(owner *actor.Actor, rect pixel.Rect, power int) *HitBox {
	center := rect.Center()
	from := pixel.V(rect.Min.X, center.Y)
	to := pixel.V(rect.Max.X, center.Y)

	if owner.GetPos().X > center.X {
		to = pixel.V(rect.Min.X, center.Y)
		from = pixel.V(rect.Max.X, center.Y)
	}
	b := &HitBox{
		rect:   rect,
		ttl:    0.5,
		power:  power,
		owner:  owner,
		dir:    pixel.L(from, to),
		hitted: make(map[*actor.Actor]struct{}),
	}

	if ai.GetByObj(owner) != nil {
		enboxes = append(enboxes, b)
	} else {
		plboxes = append(plboxes, b)
	}
	return b
}

func updateStrikes(dt float64, enemies []*actor.Actor, player *actor.Actor) {
	updatePlStrikes(dt, enemies)
	updateEnStrikes(dt, []*actor.Actor{player})
}

func updatePlStrikes(dt float64, hittable []*actor.Actor) {
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
		if b.timer < b.ttl {
			plboxes[i] = b
			i++
		}
	}

	plboxes = plboxes[:i]
}

func updateEnStrikes(dt float64, hittable []*actor.Actor) {
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
		if b.timer < b.ttl {
			enboxes[i] = b
			i++
		}
	}

	enboxes = enboxes[:i]
}
