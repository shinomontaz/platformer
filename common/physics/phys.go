package physics

import (
	"image/color"
	"math"
	"platformer/common"

	"github.com/lucasb-eyer/go-colorful"
	"github.com/shinomontaz/pixel"
	"github.com/shinomontaz/pixel/imdraw"
)

type Phys struct {
	rect pixel.Rect
	vel  pixel.Vec

	rigidity       float64
	rigidityBottom float64
	gravity        float64
	mass           float64
	friction       float64
	maxvelocity    float64

	isground   bool
	groundObj  common.O
	groundPhys *Phys
	color      color.Color
	mask       []bool

	currObjs []common.O
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

func New(r pixel.Rect, opts ...PhysOption) Phys {
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

func (p *Phys) GetVel() pixel.Vec {
	return p.vel
}

func (p *Phys) SetMass(m float64) {
	p.mass = m
}

func (p *Phys) IsGround() bool {
	return p.isground
}

func (p *Phys) GetGroundObject() common.O {
	return p.groundObj
}

func (p *Phys) SetVel(newspeed pixel.Vec) {
	p.vel = newspeed
}

func (p *Phys) SetGroundPhys(phys *Phys) {
	p.groundPhys = phys
}

func (p *Phys) Update(dt float64, objs []common.O) {
	p.currObjs = objs
	dt = math.Min(dt, 0.1)

	// 1. Трение (для частиц и снарядов)
	if p.friction > 0 && p.isground {
		p.vel.X *= (1.0 - p.friction*dt*60)
		if math.Abs(p.vel.X) < 0.5 {
			p.vel.X = 0
		}
	}

	// 2. Гравитация
	if p.mass > 0 && !p.isground {
		p.vel.Y -= p.gravity * dt
	}

	// 3. На земле — обнулить Y
	if p.isground {
		if p.vel.Y < 0 {
			p.vel.Y = 0
		}
	}

	// 4. Движение и коллизии
	groundVel := pixel.ZV
	if p.groundPhys != nil {
		groundVel = p.groundPhys.GetVel()
		p.vel = p.vel.Add(groundVel)
	}

	if p.vel.X != 0 || p.vel.Y != 0 {
		p.isground = false
		vec := p.vel.Scaled(dt)
		ground, newvel, vec, groundObj := p.StepPrediction(vec)
		if groundObj != nil && (p.groundObj == nil || groundObj.Id() != p.groundObj.Id()) {
			p.groundObj = groundObj
		}
		p.vel = newvel
		if ground > 0.05 {
			p.isground = true
			if p.vel.Y > 0 && p.vel.Y < p.gravity*dt {
				p.vel.Y = 0
			}
		}
		p.rect = p.rect.Moved(vec)
	}

	if groundVel.X != 0 || groundVel.Y != 0 {
		p.vel = p.vel.Sub(groundVel)
	}
}

// in: v - move vector, out - ground rate, new velocity, available move
func (p *Phys) StepPrediction(v pixel.Vec) (float64, pixel.Vec, pixel.Vec, common.O) {
	ground := 0.0
	vel := p.vel
	var gObj common.O
	broadbox := Broadbox(p.rect, v)
	collisiontimes := []float64{}
	if len(p.currObjs) == 0 {
		return ground, vel, v, gObj
	}
	for _, obj := range p.currObjs {
		rect := obj.Rect()
		if Isinbox(rect, broadbox) {
			if math.Abs(rect.Max.Y-p.rect.Min.Y) < 0.001 {
				l := math.Max(rect.Min.X, p.rect.Min.X)
				r := math.Min(rect.Max.X, p.rect.Max.X)
				ground = (r - l) / p.rect.W()
				gObj = obj

				continue
			}
			coltime, n := collide(p.rect, rect, v)
			if coltime == 1 { // нет пересечения
				continue
			}

			if n.X == 0 && vel.Y == 0 {
				l := math.Max(rect.Min.X, p.rect.Min.X)
				r := math.Min(rect.Max.X, p.rect.Max.X)
				ground = (r - l) / p.rect.W()
				gObj = obj
			} else {
				// Шаг вверх: низкая преграда сбоку
				stepHeight := rect.Max.Y - p.rect.Min.Y
				if n.X != 0 && vel.X != 0 && vel.Y <= 0 &&
					rect.Max.Y < p.rect.Max.Y && stepHeight > 0 && stepHeight <= p.rect.H()/2 {
					// Поднять траекторию ровно на высоту ступеньки
					v = pixel.Vec{X: v.X, Y: stepHeight + 0.01}
					continue // без coltime — движение не обрезается
				}

				collisiontimes = append(collisiontimes, coltime)
				if n.Y > 0 {
					l := math.Min(rect.Min.X, p.rect.Min.X)
					r := math.Min(rect.Max.X, p.rect.Max.X)
					ground = (r - l) / p.rect.W()
					gObj = obj
					vel.Y = -vel.Y * p.rigidityBottom
				} else if n.Y < 0 {
					vel.Y = -vel.Y * p.rigidity
				}
				if n.X != 0 {
					if (ground >= 0 || n.Y > 0) && (vel.X != 0) && (rect.Max.Y < p.rect.Max.Y) && (rect.Max.Y-p.rect.Min.Y <= p.rect.H()/4) {
						vel.X *= 0.8
						v = pixel.Vec{X: v.X, Y: v.Y + rect.Max.Y - p.rect.Min.Y + p.rect.H()/4}
					} else {
						vel.X = -vel.X * p.rigidity
						if ground != 0 {
							vel.X *= p.rigidity
						}
					}
				}
			}
		}
	}

	if len(collisiontimes) > 0 {
		mintime := collisiontimes[0]
		for _, ct := range collisiontimes {
			if mintime > ct {
				mintime = ct
			}
		}
		v.X *= mintime
		v.Y *= mintime
	}

	return ground, vel, v, gObj
}

func (p *Phys) Move(v pixel.Vec) {
	p.rect = p.rect.Moved(v)
}

func (p *Phys) SetRect(r pixel.Rect) {
	p.rect = r
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
