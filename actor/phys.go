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
		ground:    false,
		runSpeed:  run,
		walkSpeed: walk,
		jumpSpeed: jump,
		gravity:   gravity,
	}
}

func (p *Phys) SetQt(qt *common.Quadtree) {
	p.qt = qt
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
	if p.ground && move.Y > 0 {
		p.vel.Y = p.jumpSpeed
	}

	if p.vel.X != 0 || p.vel.Y != 0 {
		vec := p.vel.Scaled(dt)
		p.canmove(&vec)
		p.rect = p.rect.Moved(vec)
	}

	// jump if on the ground and the player wants to jump

}

func (p *Phys) canmove(v *pixel.Vec) {
	res := true
	p.ground = false
	moved := p.rect.Moved(*v)
	objs := p.qt.Retrieve(moved)
	if len(objs) > 0 { // precise check for each object that can intersects
		for _, obj := range objs {

			// Handle collision
			rect := obj.Rect()
			if !p.rect.Intersects(rect) {
				continue
			}

			if v.Y > 0 {
				top := p.intersectTop(rect, v)
				if top > 0 {
					v.Y -= top
					res = false
					moved = moved.Moved(pixel.Vec{0, -top})
				}
			} else {
				bottom := p.intersectBottom(rect, v)
				if bottom > 0 {
					v.Y += bottom
					moved = moved.Moved(pixel.Vec{0, bottom})
					p.ground = true
					res = false
					p.vel.Y = 0
				}
			}
			if v.X < 0 {
				left := p.intersectLeft(rect, v)
				if left > 0 {
					v.X += left
					res = false
					panic("!")
				}
			} else {
				right := p.intersectRight(rect, v)
				if right > 0 {
					v.X -= right
				}
				res = false
			}

			if !res {
				return
			}
		}
	}
}

func (p *Phys) intersectTop(r pixel.Rect, v *pixel.Vec) float64 {
	moved := p.rect.Moved(*v)

	if moved.Min.Y > r.Max.Y || (moved.Max.X <= r.Min.X || moved.Min.X >= r.Max.X) {
		return 0
	}

	val := moved.Max.Y - r.Min.Y
	if val < 0 {
		val = 0
	}
	return val
}

func (p *Phys) intersectBottom(r pixel.Rect, v *pixel.Vec) float64 {
	moved := p.rect.Moved(*v)

	if moved.Max.Y < r.Min.Y || (moved.Max.X <= r.Min.X || moved.Min.X >= r.Max.X) {
		return 0
	}

	val := r.Max.Y - moved.Min.Y
	if val < 0 {
		val = 0
	}
	return val
}

func (p *Phys) intersectLeft(r pixel.Rect, v *pixel.Vec) float64 {
	moved := p.rect.Moved(*v)

	if moved.Max.X < r.Min.X || (moved.Max.Y <= r.Min.Y || moved.Min.Y >= r.Max.Y) {
		return 0
	}

	val := r.Max.X - moved.Min.X
	if val < 0 {
		val = 0
	}
	return val
}

func (p *Phys) intersectRight(r pixel.Rect, v *pixel.Vec) float64 {
	moved := p.rect.Moved(*v)

	if moved.Min.X > r.Max.X || !(moved.Max.Y < r.Min.Y || moved.Min.Y > r.Max.Y) {
		return 0
	}
	val := r.Min.X - moved.Max.X
	if val < 0 {
		val = 0
	}
	return val
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
