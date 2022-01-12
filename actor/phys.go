package actor

import (
	"image/color"
	"math"

	"github.com/faiface/pixel"
	"github.com/lucasb-eyer/go-colorful"

	"platformer/common"
)

type Phys struct {
	rect      pixel.Rect
	vel       pixel.Vec
	runSpeed  float64
	walkSpeed float64
	jumpSpeed float64
	gravity   float64
	ground    bool
	color     color.Color
	qt        *common.Quadtree
}

func NewPhys(r pixel.Rect, run, walk, jump, gravity float64) Phys {
	return Phys{
		rect:      r,
		color:     colorful.HappyColor(),
		ground:    true,
		runSpeed:  run,
		walkSpeed: walk,
		jumpSpeed: jump,
		gravity:   gravity,
	}
}

func (p *Phys) SetQt(qt *common.Quadtree) {
	p.qt = qt
}

func (p *Phys) intersects(obj common.Objecter) bool {
	rect := obj.Rect()
	if p.rect.Max.X <= rect.Min.X || p.rect.Min.X >= rect.Max.X {
		return false
	}
	if p.rect.Min.Y > rect.Max.Y || p.rect.Max.Y < rect.Max.Y {
		return false
	}

	return true
}

func (p *Phys) Update(dt float64, move pixel.Vec) {
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

	if !p.ground {
		p.vel.Y -= p.gravity
	}
	vec := p.vel.Scaled(dt)
	if p.canmove(&vec) {
		p.rect = p.rect.Moved(vec)
	} else {
		p.vel = pixel.ZV
	}

	// jump if on the ground and the player wants to jump
	if p.ground && move.Y > 0 {
		p.vel.Y = p.jumpSpeed
	}

}

func (p *Phys) canmove(v *pixel.Vec) bool {
	res := true
	p.ground = false
	if p.vel.Y != 0 {
		objs := p.qt.Retrieve(p.rect)
		if len(objs) > 0 { // precise check for each object that can intersects
			for _, obj := range objs {

				// Handle collision
				rect := obj.Rect()
				if !p.rect.Intersects(rect) {
					continue
				}

				// if p.intersectTop(rect) {
				// 	// do top intersection
				// }
				// if p.intersectBottom(rect) {
				// 	// do ground intersection
				// }
				// if p.intersectBottom(rect) {
				// 	// do ground intersection
				// }

				if p.vel.Y < 0 {
					p.rect = p.rect.Moved(pixel.V(0, rect.Max.Y-p.rect.Min.Y))
					p.ground = true
				} else if p.vel.Y == 0 {
					//					p.rect = p.rect.Moved(pixel.V(0, rect.Min.Y-p.rect.Max.Y))
					res = false
				}
				if !res {
					return res
				}
			}
		}
	}

	return res
}

func (p *Phys) GetVel() pixel.Vec {
	return p.vel
}

func (p *Phys) Move(v pixel.Vec) {
	p.rect = p.rect.Moved(v)
}

func (p *Phys) GetRect() pixel.Rect {
	return p.rect
}

func (p *Phys) Draw(t pixel.Target) {
	// imd := imdraw.New(nil)

	// vertices := p.rect.Vertices()

	// imd.Color = p.color
	// for _, v := range vertices {
	// 	imd.Push(v)
	// }
	// imd.Rectangle(1)

	// imd.Draw(t)
}
