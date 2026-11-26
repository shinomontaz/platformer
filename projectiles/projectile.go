package projectiles

import (
	"math"
	"platformer/activities"
	"platformer/common"
	"platformer/common/physics"
	"platformer/common/rand"
	"platformer/particles"

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
	rigidity float64
	friction float64
	time     float64

	strength float64
	owner    common.Actorer

	anim          Animater
	sprite        *pixel.Sprite
	animSpriteNum int

	hb       *activities.HitBox
	onkill   OnKillHandler
	ttl      float64
	trace    float64
	traceAcc float64
	traceCol uint32
}

func (p *projectile) Init(g float64) {
	pic, fr := p.anim.GetSprite("projectile", 0)
	w, h := fr.W(), fr.H()
	p.rect = pixel.R(p.pos.X-w/2, p.pos.Y-h/2, p.pos.X+w/2, p.pos.Y+h/2)

	p.phys = physics.New(p.rect,
		physics.WithGravity(g),
		physics.WithMass(p.mass),
		physics.WithRigidity(p.rigidity),
		physics.WithFriction(p.friction),
	)
	p.active = true
	p.sprite = pixel.NewSprite(nil, pixel.Rect{})
	p.phys.SetVel(p.force)

	p.traceCol = dominantColor(pic, fr)
}

func (p *projectile) update(dt float64, objs, spec []common.O) {
	if p.ttl > 0 && p.time >= p.ttl {
		p.active = false
		if p.onkill != nil {
			p.onkill(p.pos)
		}
		return
	}
	if p.trace > 0 {
		p.traceAcc += dt
		for p.traceAcc >= p.trace {
			p.traceAcc -= p.trace
			pr := p.phys.GetRect()
			// random pos across projectile height:

			particles.AddTrace(pixel.V(pr.Center().X, pr.Min.Y+rand.Float64()*pr.H()), p.traceCol)
		}
	}

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
		ScaledXY(pixel.ZV, pixel.V(p.dir, 1)).
		Moved(rect.Center()),
	)
}

func dominantColor(pic pixel.Picture, fr pixel.Rect) uint32 {
	pd := pixel.PictureDataFromPicture(pic)
	var r, g, b, n int
	for y := int(fr.Min.Y); y < int(fr.Max.Y); y++ {
		for x := int(fr.Min.X); x < int(fr.Max.X); x++ {
			c := pd.Pix[y*pd.Stride+x]
			if c.A < 128 {
				continue
			}
			r += int(c.R)
			g += int(c.G)
			b += int(c.B)
			n++
		}
	}
	if n == 0 {
		return 0x808080FF // пустой кадр — серый fallback
	}
	r = int(uint8(r/n) & 0xE0) // квантование до 32 ступеней/канал
	g = int(uint8(g/n) & 0xE0)
	b = int(uint8(b/n) & 0xE0)
	return uint32(r)<<24 | uint32(g)<<16 | uint32(b)<<8 | 0xFF
}
