package actor

import (
	"image/color"
	"math"

	"github.com/faiface/pixel"
	"github.com/faiface/pixel/imdraw"
	"github.com/lucasb-eyer/go-colorful"

	"platformer/common"
)

type Phys struct {
	rect     pixel.Rect
	vel      pixel.Vec
	maxspeed float64
	gravity  float64
	ground   bool
	color    color.Color
	qt       *common.Quadtree
}

func NewPhys(r pixel.Rect, run, gravity float64) Phys {
	return Phys{
		rect:     r,
		color:    colorful.HappyColor(),
		ground:   false,
		maxspeed: run,
		gravity:  gravity,
	}
}

func (p *Phys) SetQt(qt *common.Quadtree) {
	p.qt = qt
}

func (p *Phys) GetVel() *pixel.Vec {
	return &p.vel
}

func (p *Phys) Update(dt float64, move *pixel.Vec) {
	// do speed update by move vec
	if p.ground {
		if move.X != 0 {
			p.vel.X += move.X
		} else {
			p.vel.X /= 1.1
			if math.Abs(p.vel.X) <= p.maxspeed/20 {
				p.vel.X = 0
			}
		}
	}

	if !p.ground {
		p.vel.Y -= p.gravity
	}
	if p.ground && move.Y > 0 {
		p.vel.Y = move.Y
	}

	if p.vel.X != 0 || p.vel.Y != 0 {
		vec := p.vel.Scaled(dt)
		p.collide(&vec) // p.vel and vec can be updated here
		p.rect = p.rect.Moved(vec)
	}
}

func (p *Phys) collide(v *pixel.Vec) {
	p.ground = false
	broadbox := Broadbox(p.rect, *v)
	objs := p.qt.Retrieve(broadbox)
	collisiontimes := []float64{}
	if len(objs) > 0 { // precise check for each object that can intersects
		for _, obj := range objs {
			rect := obj.Rect()
			if Isinbox(rect, broadbox) {
				if rect.Max.Y == p.rect.Min.Y {
					p.ground = true
					continue
				}
				coltime, n := DoCollision(p.rect, rect, *v) // coltime in [0,1], where 1 means no collision
				if coltime == 1 {                           // no collision
					continue
				}

				if n.X == 0 && p.vel.Y == 0 { // here we step on ground, no collision actually
					p.ground = true
				} else {
					collisiontimes = append(collisiontimes, coltime)
					if n.Y > 0 {
						p.ground = true
						p.vel.Y = 0
					} else if n.Y < 0 {
						p.vel.Y = -p.vel.Y * 0.5
					}
					if n.X != 0 && !p.ground {
						p.vel.X = -p.vel.X * 0.5
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
func DoCollision(p, q pixel.Rect, v pixel.Vec) (float64, pixel.Vec) {
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
