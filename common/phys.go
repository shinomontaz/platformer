package common

import (
	"image/color"
	"math"

	"github.com/lucasb-eyer/go-colorful"
	"github.com/shinomontaz/pixel"
	"github.com/shinomontaz/pixel/imdraw"
)

type Phys struct {
	rect  pixel.Rect
	vel   pixel.Vec
	force pixel.Vec

	rigidity       float64
	rigidityBottom float64
	gravity        float64
	mass           float64
	friction       float64
	maxvelocity    float64

	isground bool
	color    color.Color

	currObjs []Objecter
}

type PhysOption func(p *Phys)

func WithGravity(g float64) PhysOption {
	return func(p *Phys) {
		p.gravity = g
	}
}

func WithRigidity(r float64) PhysOption {
	return func(p *Phys) {
		p.rigidity = r
	}
}

func WithRigidityBottom(r float64) PhysOption {
	return func(p *Phys) {
		p.rigidityBottom = r
	}
}

func WithFriction(f float64) PhysOption {
	return func(p *Phys) {
		p.friction = f
	}
}

func WithMass(m float64) PhysOption {
	return func(p *Phys) {
		p.mass = m
	}
}

func WithMaxVel(max float64) PhysOption {
	return func(p *Phys) {
		p.maxvelocity = max
	}
}

func NewPhys(r pixel.Rect, opts ...PhysOption) Phys {
	p := Phys{
		rect:           r,
		color:          colorful.HappyColor(),
		rigidity:       0.5,
		rigidityBottom: 0,
		friction:       0.01,
		vel:            pixel.ZV,
		mass:           1,
	}

	for _, opt := range opts {
		opt(&p)
	}

	return p
}

func (p *Phys) GetVel() *pixel.Vec {
	return &p.vel
}

func (p *Phys) SetMass(m float64) {
	p.mass = m
}

func (p *Phys) IsGround() bool {
	return p.isground
}

func (p *Phys) Apply(force pixel.Vec) {
	p.force = p.force.Add(force)
}

func (p *Phys) SetSpeed(newspeed pixel.Vec) {
	p.vel = newspeed
}

func (p *Phys) Update(dt float64, objs []Objecter) {
	p.currObjs = objs

	if p.vel.X != 0 && math.Abs(p.vel.X) >= dt*p.friction*p.gravity && p.isground {
		sign := -1.0
		if p.vel.X > 0 {
			sign = 1.0
		}
		p.force.X -= sign * p.friction * p.mass * p.gravity // Friction force apply
	}

	if math.Abs(p.vel.X) <= dt*p.friction*p.mass*p.gravity { // static friction force check: no mass in this equation
		p.vel.X = 0
	}

	if p.mass > 0 && !p.isground {
		p.force.Y -= p.gravity // Gravity force to the forces
	}

	if p.isground {
		if p.vel.Y < 0 {
			p.vel.Y = 0
		}
		if p.force.Y < 0 {
			p.force.Y = 0
		}
	}

	dt = math.Min(dt, 1.0)
	p.vel = p.vel.Add(p.force.Scaled(dt))

	if p.vel.X != 0 || p.vel.Y != 0 {
		p.isground = false
		vec := p.vel.Scaled(dt)
		ground, newvel, vec := p.StepPrediction(vec) // p.vel and vec can be updated here
		p.vel = newvel
		if ground > 0 {
			p.isground = true
		}
		p.rect = p.rect.Moved(vec)
	}

	p.force = pixel.ZV
}

// in: v - move vector, out - ground rate, new velocity, available move
func (p *Phys) StepPrediction(v pixel.Vec) (float64, pixel.Vec, pixel.Vec) {
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
				if n.Y > 0 { // hit floor
					l := math.Min(rect.Min.X, p.rect.Min.X)
					r := math.Min(rect.Max.X, p.rect.Max.X)
					ground = (r - l) / p.rect.W()
					//					vel.Y = 0
					vel.Y = -vel.Y * p.rigidityBottom
				} else if n.Y < 0 { // hit ceiling
					vel.Y = -vel.Y * p.rigidity
				}
				if n.X != 0 {
					vel.X = -vel.X * p.rigidity // horisontal bouncing
					if ground != 0 {
						vel.X *= p.rigidity
					}
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
