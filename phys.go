package main

import (
	"math"

	"github.com/faiface/pixel"
)

type phys struct {
	rect      pixel.Rect
	vel       pixel.Vec
	runSpeed  float64
	walkSpeed float64
	jumpSpeed float64
	gravity   float64
	ground    bool
	//	ctrl      IController
}

func (p *phys) update(dt float64, move pixel.Vec, platforms []*platform) {
	// apply controls
	switch {
	case math.Abs(move.X) == 1:
		switch {
		case move.X < 0:
			p.vel.X = -p.walkSpeed
		case move.X > 0:
			p.vel.X = +p.walkSpeed
		}
	case math.Abs(move.X) == 2:
		switch {
		case move.X < 0:
			p.vel.X = -p.runSpeed
		case move.X > 0:
			p.vel.X = +p.runSpeed
		}
	default:
		p.vel.X = 0
	}

	// apply gravity and velocity
	p.vel.Y -= p.gravity
	p.rect = p.rect.Moved(p.vel.Scaled(dt))

	// check collisions against each platform
	p.ground = false
	if p.vel.Y != 0 {
		for _, pl := range platforms {
			if p.rect.Max.X <= pl.rect.Min.X || p.rect.Min.X >= pl.rect.Max.X {
				continue
			}
			if p.rect.Min.Y > pl.rect.Max.Y || p.rect.Max.Y < pl.rect.Max.Y {
				continue
			}
			if p.vel.Y < 0 {
				p.rect = p.rect.Moved(pixel.V(0, pl.rect.Max.Y-p.rect.Min.Y))
				p.ground = true
			} else {
				p.rect = p.rect.Moved(pixel.V(0, pl.rect.Min.Y-p.rect.Max.Y))
				//				p.vel.Y = -p.vel.Y
			}
			p.vel.Y = 0
		}
	}

	// jump if on the ground and the player wants to jump
	if p.ground && move.Y > 0 {
		p.vel.Y = p.jumpSpeed
	}
}
