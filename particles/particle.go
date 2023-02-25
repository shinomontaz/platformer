package particles

import (
	"image/color"
	"platformer/common"

	"github.com/shinomontaz/pixel"
	"github.com/shinomontaz/pixel/imdraw"
)

type particle struct {
	phys   common.Phys
	active bool
	color  color.Color
	//	color    uint32
	ttl float64
}

func newParticle(pos, force pixel.Vec, mass, size float64, col color.Color) particle {
	rect := pixel.R(pos.X-size/2, pos.Y-size/2, pos.X+size/2, pos.Y+size/2)
	phys := common.NewPhys(rect,
		common.WithGravity(grav),
		common.WithMass(mass),
	)
	phys.Apply(force)
	p := particle{
		active: true,
		color:  col,
		ttl:    2,
		phys:   phys,
	}

	return p
}

func (p *particle) update(dt float64, objs []common.Objecter) {
	if p.ttl <= 0 {
		p.active = false
		return
	}

	p.ttl -= dt
	p.phys.Update(dt, objs)
}

func (p *particle) draw(imd *imdraw.IMDraw) {
	vertices := p.phys.GetRect().Vertices()
	imd.Color = p.color
	for _, v := range vertices {
		imd.Push(v)
	}
	imd.Rectangle(1) // filled rectangle
}

// func (p *particle) draw(t pixel.Target) {
// 	imd := imdraw.New(nil)
// 	r := p.phys.GetRect()
// 	vertices := r.Vertices()
// 	imd.Color = p.color
// 	for _, v := range vertices {
// 		imd.Push(v)
// 	}
// 	imd.Rectangle(1) // filled rectangle
// 	imd.Draw(t)
// }
