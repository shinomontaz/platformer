package particles

import (
	"platformer/common"

	"github.com/shinomontaz/pixel"
)

type particle struct {
	phys   common.Phys
	active bool
	//	color  color.Color
	color    uint32
	ttl      float64
	pos      pixel.Vec
	force    pixel.Vec
	mass     float64
	size     float64
	rigidity float64
}

func (p *particle) Init() {
	rect := pixel.R(p.pos.X-p.size/2, p.pos.Y-p.size/2, p.pos.X+p.size/2, p.pos.Y+p.size/2)
	p.phys = common.NewPhys(rect,
		common.WithGravity(grav),
		common.WithMass(p.mass),
		common.WithRigidity(p.rigidity),
	)
	p.phys.Apply(p.force)
}

func (p *particle) update(dt float64, objs []common.Objecter) {
	if p.ttl <= 0 {
		p.active = false
		return
	}

	p.ttl -= dt
	p.phys.Update(dt, objs)
	p.pos = p.phys.GetRect().Center()
}

//func (p *particle) draw(imd *imdraw.IMDraw) {
// vertices := p.phys.GetRect().Vertices()
// imd.Color = p.color
// for _, v := range vertices {
// 	imd.Push(v)
// }
// imd.Rectangle(1) // filled rectangle
//}

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
