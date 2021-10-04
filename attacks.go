package main

import "github.com/faiface/pixel"

type Attack struct {
	pos   pixel.Vec
	life  float64
	power int
	owner Objecter

	rect pixel.Rect
	vel  pixel.Vec
	mdt  float64
	anim *Anim
}

type Attacks struct {
	attacks []Attack
}

func NewAttacks() *Attacks {
	return &Attacks{
		attacks: make([]Attack, 0),
	}
}

func (atks *Attacks) Add(a Attack) {
	atks.attacks = append(atks.attacks, a)
}

func (atks *Attacks) Update(dt float64, objs []Objecter) {
	i := 0
	for _, a := range atks.attacks { // update lifetimes of all attacks, remove otdated
		a.mdt += dt
		if a.life == 0 || a.mdt < a.life {
			// update rectange of attack obj in case it have velocity
			if a.vel.Len() > 0 {
				a.rect = a.rect.Moved(a.vel)
			}
			atks.attacks[i] = a
			i++
		}
	}
	atks.attacks = atks.attacks[:i]

	// move projectile by velocity vec, check intersections
	// we not use quad tree or other gamedev datastructs
	for _, obj := range objs {
		for _, a := range atks.attacks {
			if a.owner == obj {
				continue
			}
			if obj.Rect().Intersects(a.rect) {
				obj.Hit(a.pos, a.vel, a.power)
			}
		}
	}

}

func (atks *Attacks) Draw(t pixel.Target) {
	for _, a := range atks.attacks {
		if a.anim != nil {
			// draw attack anim - that is actually projectile
		}
	}
}
