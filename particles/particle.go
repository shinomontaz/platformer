package particles

import (
	"platformer/common"

	"github.com/faiface/pixel"
	"github.com/shinomontaz/pixel"
)

type particle struct {
	active   bool
	rect     pixel.Rect
	mass     float64
	color    uint32
	ttl      float64
	rigidity float64
	gravity  float64
	pos      pixel.Vec
	prev     pixel.Vec
	vel      pixel.Vec
	force    pixel.Vec
}

func (p *particle) update(dt float64, objs []common.Objecter) {
	if p.ttl <= 0 {
		p.active = false
		return
	}

	p.prev = p.pos
	p.life -= dt
	//	vec := p.vel.Scaled(dt).ScaledXY(p.force).Scaled(p.mass)

	//	p.phys.Update(dt, &vel, objs)

	p.vel = p.vel.Scaled(dt).ScaledXY(p.force).Scaled(p.mass)
	ground_rate, newvel, move := common.StepPrediction(dt, p.rect, objs) // (dt float64, rect pixel.Rect, vel pixel.Vec, currObjs []common.Objecter), out - ground rate, new velocity, available move

	// new position
	// Broadbox
	// detect collision
	// find new pos, new velocity with regard of rigidity
	// change force with gravity
	// vec = new vec to move

	p.force.Y -= dt * p.gravity
	if p.force.X > 0 {
		p.force.X -= dt * p.gravity * p.mass
	} else {
		p.force.X = 0
	}
}
