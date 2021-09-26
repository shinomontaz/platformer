package main

import "github.com/faiface/pixel"

type Attack struct {
	pos    pixel.Vec
	size   float64
	life   float64
	active bool
	power  int
	owner  Objecter

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
	// move projectile by velocity vec, check intersections
	// we not use quad tree or other gamedev datastructs
	for _, obj := range objs {
		for _, a := range atks.attacks {
			if a.owner == obj {
				continue
			}
			if obj.Rect().Contains(a.pos) {
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
