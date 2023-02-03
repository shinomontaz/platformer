package activities

import (
	"platformer/common"

	"github.com/shinomontaz/pixel"
)

var plinteractions []HitBox

func init() {
	plinteractions = make([]HitBox, 0)
}

func AddInteraction(owner common.Actorer, rect pixel.Rect, power int, speed pixel.Vec) {
	b := HitBox{
		rect:   rect,
		ttl:    0.2,
		power:  power,
		owner:  owner,
		hitted: make(map[common.Actorer]struct{}),
		speed:  speed,
	}

	if owner.GetId() == 1 {
		plinteractions = append(plinteractions, b)
	}
}

func UpdateInteractions(dt float64, interactibles []common.Actorer) {
	i := 0
	for _, b := range plinteractions {
		for _, hh := range interactibles {
			if hh == b.owner {
				continue
			}
			if _, ok := b.hitted[hh]; ok {
				continue
			}
			r := hh.GetRect()
			if b.rect.Intersects(r) {
				hh.OnInteract() // active side is always hero
				b.hitted[hh] = struct{}{}
			}
		}
		b.timer += dt
		if b.timer < b.ttl {
			plinteractions[i] = b
			i++
		}
	}

	plinteractions = plinteractions[:i]
}
