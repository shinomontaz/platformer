package common

import (
	"image/color"
	"math"

	"github.com/lucasb-eyer/go-colorful"
	"github.com/shinomontaz/pixel"
	"github.com/shinomontaz/pixel/imdraw"
)

const MINSPEED = 4

type Phys struct {
	rect     pixel.Rect
	vel      pixel.Vec
	gravity  float64
	ground   bool
	color    color.Color
	iswater  bool
	rigidity float64
	mass     float64
	currObjs []Objecter
}

func NewPhys(r pixel.Rect, vel pixel.Vec, rigidity, gravity, mass float64) Phys {
	return Phys{
		rect:     r,
		color:    colorful.HappyColor(),
		ground:   false,
		rigidity: rigidity,
		gravity:  gravity,
		vel:      vel,
		mass:     mass,
	}
}

func (p *Phys) GetVel() *pixel.Vec {
	return &p.vel
}

func (p *Phys) Impulse(v pixel.Vec) {
	p.vel = p.vel.Add(v)
}

func (p *Phys) SetMass(m float64) {
	p.mass = m
}

func (p *Phys) IsGround() bool {
	return p.ground
}

func (p *Phys) SetWater(iswater bool) {
	p.iswater = iswater
}

func (p *Phys) Update(dt float64, move *pixel.Vec, objs []Objecter) {
	// do speed update by move vec
	p.currObjs = objs
	if p.ground {
		if move.X != 0 {
			p.vel.X += move.X
		} else {
			p.vel.X /= 1.1
			if math.Abs(p.vel.X) <= MINSPEED {
				p.vel.X = 0
			}
		}
	}

	if !p.ground {
		p.vel.Y -= p.gravity
	}
	if p.iswater {
		p.vel = p.vel.Scaled(0.9)
		p.vel.Y += (1 - p.mass) * p.gravity * 1.1
	}
	if p.ground && move.Y > 0 {
		p.vel.Y = move.Y
	}

	if p.vel.X != 0 || p.vel.Y != 0 {
		p.ground = false
		vec := p.vel.Scaled(dt)
		ground, vel, vec := p.StepPrediction(vec) // p.vel and vec can be updated here
		oldvel := p.vel
		p.vel = vel
		if ground > 0 {
			p.ground = true
			// do a vertical bouncing
			if oldvel.Y < 0 && oldvel.Y < -MINSPEED {
				p.vel.Y = -oldvel.Y * p.rigidity
			}
		}
		p.rect = p.rect.Moved(vec)
	}
}

func (p *Phys) StepPrediction(v pixel.Vec) (float64, pixel.Vec, pixel.Vec) { // in: v - velocity
	ground := 0.0
	vel := p.vel
	broadbox := Broadbox(p.rect, v)
	collisiontimes := []float64{}
	if len(p.currObjs) == 0 {
		return ground, vel, v
	}
	// precise check for each object that can intersects
	for _, obj := range p.currObjs {
		rect := obj.Rect()
		if Isinbox(rect, broadbox) { // return !(p.Max.X < q.Min.X || p.Min.X > q.Max.X || p.Max.Y < q.Min.Y || p.Min.Y > q.Max.Y)
			if rect.Max.Y == p.rect.Min.Y {
				l := math.Max(rect.Min.X, p.rect.Min.X)
				r := math.Min(rect.Max.X, p.rect.Max.X)
				ground = (r - l) / p.rect.W()
				continue
			}
			coltime, n := collide(p.rect, rect, v) // coltime in [0,1], where 1 means no collision
			if coltime == 1 {                      // no collision
				continue
			}

			if n.X == 0 && vel.Y == 0 { // here we step on ground, no collision actually
				l := math.Max(rect.Min.X, p.rect.Min.X)
				r := math.Min(rect.Max.X, p.rect.Max.X)
				ground = (r - l) / p.rect.W()
			} else {
				collisiontimes = append(collisiontimes, coltime)
				if n.Y > 0 {
					l := math.Min(rect.Min.X, p.rect.Min.X)
					r := math.Min(rect.Max.X, p.rect.Max.X)
					ground = (r - l) / p.rect.W()
					vel.Y = 0
				} else if n.Y < 0 {
					vel.Y = -vel.Y * 0.5
				}
				if n.X != 0 && ground == 0 {
					vel.X = -vel.X * 0.5
				}
			}
		}
	}

	if len(collisiontimes) > 0 {
		mintime := collisiontimes[0]
		// find minimal collision time
		for _, ct := range collisiontimes {
			if mintime > ct {
				mintime = ct
			}
		}

		v.X *= mintime
		v.Y *= mintime
	}

	return ground, vel, v
}

func (p *Phys) Move(v pixel.Vec) {
	p.rect = p.rect.Moved(v)
}

func (p *Phys) GetRect() pixel.Rect {
	return p.rect
}

func (p *Phys) Draw(t pixel.Target) {
	imd := imdraw.New(nil)

	vertices := p.rect.Vertices()

	imd.Color = p.color
	for _, v := range vertices {
		imd.Push(v)
	}
	imd.Rectangle(1)

	imd.Draw(t)
}

// sweptAABB implemented here
//Box b1, Box b2, float& normalx, float& normaly
func collide(p, q pixel.Rect, v pixel.Vec) (float64, pixel.Vec) {
	n := pixel.ZV
	var xInvEntry, yInvEntry, xInvExit, yInvExit float64
	if v.X > 0 {
		xInvEntry = q.Min.X - p.Max.X
		xInvExit = q.Max.X - p.Min.X
	} else {
		xInvEntry = q.Max.X - p.Min.X
		xInvExit = q.Min.X - p.Max.X
	}

	if v.Y > 0 {
		yInvEntry = q.Min.Y - p.Max.Y
		yInvExit = q.Max.Y - p.Min.Y
	} else {
		yInvEntry = q.Max.Y - p.Min.Y
		yInvExit = q.Min.Y - p.Max.Y
	}

	var xEntry, yEntry, xExit, yExit float64

	if v.X == 0 {
		xEntry = math.Inf(-1)
		xExit = math.Inf(1)
	} else {
		xEntry = xInvEntry / v.X
		xExit = xInvExit / v.X
	}

	if v.Y == 0 {
		yEntry = math.Inf(-1)
		yExit = math.Inf(1)
	} else {
		yEntry = yInvEntry / v.Y
		yExit = yInvExit / v.Y
	}

	entryTime := math.Max(xEntry, yEntry)
	exitTime := math.Min(xExit, yExit)

	if entryTime > exitTime || (xEntry < 0 && yEntry < 0) || xEntry > 1 || yEntry > 1 {
		return 1, n
	} else {
		if xEntry > yEntry {
			if xInvEntry < 0 {
				n.X = 1
			} else {
				n.X = -1
			}
		} else {
			if yInvEntry < 0 {
				n.Y = 1
			} else {
				n.Y = -1
			}
		}
	}

	return entryTime, n
}

func Isinbox(p, q pixel.Rect) bool {
	//return !(b1.x + b1.w < b2.x || b1.x > b2.x + b2.w || b1.y + b1.h < b2.y || b1.y > b2.y + b2.h);
	return !(p.Max.X < q.Min.X || p.Min.X > q.Max.X || p.Max.Y < q.Min.Y || p.Min.Y > q.Max.Y)
}

func Broadbox(r pixel.Rect, v pixel.Vec) pixel.Rect {
	var minx, miny, w, h float64
	if v.X > 0 {
		minx = r.Min.X
		w = v.X + r.W()
	} else {
		minx = r.Min.X + v.X
		w = r.W() - v.X
	}
	if v.Y > 0 {
		miny = r.Min.Y
		h = v.Y + r.H()
	} else {
		miny = r.Min.Y + v.Y
		h = r.H() - v.Y
	}

	return pixel.R(minx, miny, minx+w, miny+h)
}
