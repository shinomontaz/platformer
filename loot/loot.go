package loot

import (
	"math"
	"platformer/common"
	"platformer/sound"

	"github.com/shinomontaz/pixel"
)

var counter int

type Loot struct {
	id   int
	phys common.Phys
	mass float64

	rect pixel.Rect

	animdir       float64
	anim          common.Animater
	sprite        *pixel.Sprite
	animSpriteNum int
	counter       float64
	dir           float64
	vec           pixel.Vec // delta speed
	vel           pixel.Vec

	sounds map[string]soundeffect

	sbrs []common.Subscriber
	grav float64
}

func New(anim common.Animater, rect pixel.Rect, opts ...Option) *Loot {
	counter++
	a := &Loot{
		id:      counter,
		anim:    anim,
		rect:    rect,
		mass:    0.1,
		dir:     1,
		animdir: 1,
		vel:     pixel.ZV,
		sounds:  make(map[string]soundeffect),
		sbrs:    make([]common.Subscriber, 0),
		sprite:  pixel.NewSprite(nil, pixel.Rect{}),
	}

	for _, opt := range opts {
		opt(a)
	}

	p := common.NewPhys(rect,
		common.WithGravity(a.grav),
		common.WithMass(a.mass),
		common.WithRigidityBottom(0.5),
	) // TODO does we really need phys to know run and walk speeds?
	a.phys = p

	return a
}

func (a *Loot) Move(v pixel.Vec) {
	a.rect = a.rect.Moved(v)
	a.phys.Move(v)
}

func (a *Loot) GetPos() pixel.Vec {
	return a.rect.Min
}

func (a *Loot) GetRect() pixel.Rect {
	return a.rect
}

func (a *Loot) Update(dt float64, objs []common.Objecter) {
	//	a.phys.Update(dt, &a.vec, objs)
	a.phys.Update(dt, objs)

	a.vec = pixel.ZV
	newspeed := a.phys.GetVel()
	a.vel = *newspeed

	a.rect = a.phys.GetRect()

	a.counter += dt
	a.animSpriteNum = int(math.Floor(a.counter / 0.2))
}

func (a *Loot) GetDir() int {
	return int(a.dir)
}

func (a *Loot) Draw(t pixel.Target) {
	pic, rect := a.anim.GetSprite("idle", a.animSpriteNum)
	a.sprite.Set(pic, rect)
	drawrect := a.rect
	a.sprite.Draw(t, pixel.IM.
		ScaledXY(pixel.ZV, pixel.V(
			drawrect.W()/a.sprite.Frame().W(),
			drawrect.H()/a.sprite.Frame().H(),
		)).
		Moved(drawrect.Center()),
	)
}

func (a *Loot) AddSound(event string) {
	if s, ok := a.sounds[event]; ok {
		// select random sound
		i := int(math.Round(common.GetRandFloat() * float64(len(s.List)-1)))
		sound.AddEffect(s.List[i], a.rect.Center())
	}
}

func (a *Loot) IsGround() bool {
	return a.phys.IsGround()
}

func (a *Loot) Inform(e int, v pixel.Vec) {
	for _, s := range a.sbrs {
		s.Listen(e, v)
	}
}

func (a *Loot) AddListener(s common.Subscriber) {
	a.sbrs = append(a.sbrs, s)
}

func (a *Loot) GetId() int {
	return a.id
}
