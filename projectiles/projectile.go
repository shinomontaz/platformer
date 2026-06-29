package projectiles

import (
	"math"
	"platformer/activities"
	"platformer/common"
	"platformer/common/physics"

	"github.com/shinomontaz/pixel"
)

type OnKillHandler func(pos pixel.Vec)

type projectile struct {
	active   bool
	rect     pixel.Rect
	phys     physics.Phys
	pos      pixel.Vec
	force    pixel.Vec
	vel      pixel.Vec
	dir      float64
	mass     float64
	size     float64
	rigidity float64
	friction float64
	time     float64

	strength float64
	owner    common.Actorer

	anim          Animater
	sprite        *pixel.Sprite
	animSpriteNum int

	hb     *activities.HitBox
	onkill OnKillHandler
}

func (p *projectile) Init(g float64) {
	p.rect = pixel.R(p.pos.X-p.size/2, p.pos.Y-p.size/2, p.pos.X+p.size/2, p.pos.Y+p.size/2)
	p.phys = physics.New(p.rect,
		physics.WithGravity(g),
		physics.WithMass(p.mass),
		physics.WithRigidity(p.rigidity),
		physics.WithFriction(p.friction),
	)
	p.active = true
	p.sprite = pixel.NewSprite(nil, pixel.Rect{})
	p.phys.Apply(p.force)
}

func (p *projectile) update(dt float64, objs, spec []common.O) {
	p.phys.Update(dt, objs)

	p.time += dt

	p.vel = p.phys.GetVel()

	p.hb.Move(p.vel.Scaled(dt))
	if ((p.vel.X <= 1 && p.vel.X >= -1) && p.vel.Y < 1) || p.phys.IsGround() {
		p.active = false
		if p.onkill == nil {
			return
		}

		p.onkill(p.pos)
		return
	}

	p.animSpriteNum = int(math.Floor(p.time / 0.1))

}

func (p *projectile) draw(t pixel.Target) {
	pic, sprect := p.anim.GetSprite("projectile", p.animSpriteNum)
	p.sprite.Set(pic, sprect)

	rect := p.phys.GetRect()
	p.sprite.Draw(t, pixel.IM.
		ScaledXY(pixel.ZV, pixel.V(
			rect.W()/p.sprite.Frame().W(),
			rect.H()/p.sprite.Frame().H(),
		)).
		ScaledXY(pixel.ZV, pixel.V(-1*p.dir, 1)).
		Moved(rect.Center()),
	)
}
