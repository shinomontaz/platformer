package main

import (
	"image/color"
	"math"

	"github.com/faiface/pixel"
	"github.com/faiface/pixel/imdraw"
	"github.com/lucasb-eyer/go-colorful"
)

type phys struct {
	rect      pixel.Rect
	vel       pixel.Vec
	runSpeed  float64
	walkSpeed float64
	jumpSpeed float64
	gravity   float64
	ground    bool
	color     color.Color
}

func NewPhys() phys {
	return phys{
		color:  colorful.HappyColor(),
		ground: true,
	}
}

func (p *phys) Intersects(obj Objecter) bool {
	rect := obj.Rect()
	if p.rect.Max.X <= rect.Min.X || p.rect.Min.X >= rect.Max.X {
		return false
	}
	if p.rect.Min.Y > rect.Max.Y || p.rect.Max.Y < rect.Max.Y {
		return false
	}

	return true
}

func (p *phys) update(dt float64, move pixel.Vec, objs []Objecter) {
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
		for _, obj := range objs {
			if !p.Intersects(obj) {
				continue
			}

			// Handle collision
			rect := obj.Rect()

			if p.vel.Y < 0 {
				p.rect = p.rect.Moved(pixel.V(0, rect.Max.Y-p.rect.Min.Y))
				p.ground = true
			} else {
				p.rect = p.rect.Moved(pixel.V(0, rect.Min.Y-p.rect.Max.Y))
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

func (p *phys) draw(t pixel.Target) {
	imd := imdraw.New(nil)

	vertices := p.rect.Vertices()

	imd.Color = p.color
	for _, v := range vertices {
		imd.Push(v)
	}
	imd.Rectangle(1)
	imd.Draw(t)
}
