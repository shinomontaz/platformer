package particles

import (
	"platformer/common"
	"platformer/common/physics"

	"github.com/shinomontaz/pixel"
)

type particle struct {
	phys     physics.Phys
	active   bool
	color    uint32
	ttl      float64
	pos      pixel.Vec
	force    pixel.Vec
	mass     float64
	size     float64
	rigidity float64
	friction float64
}

func (p *particle) Init() {
	rect := pixel.R(p.pos.X-p.size/2, p.pos.Y-p.size/2, p.pos.X+p.size/2, p.pos.Y+p.size/2)
	p.phys = physics.New(rect,
		physics.WithGravity(grav),
		physics.WithMass(p.mass),
		physics.WithRigidity(p.rigidity),
		physics.WithFriction(p.friction),
	)
	p.phys.SetVel(p.force)
}

func (p *particle) update(dt float64, objs []common.O) {
	if p.ttl <= 0 {
		p.active = false
		return
	}

	p.ttl -= dt
	p.phys.Update(dt, objs)
	p.pos = p.phys.GetRect().Center()
}
