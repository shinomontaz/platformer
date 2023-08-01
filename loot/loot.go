package loot

import (
	"math"
	"platformer/common"
	"platformer/sound"

	"github.com/shinomontaz/pixel"
)

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
	appliedForce  pixel.Vec

	sounds map[string]soundeffect

	grav float64

	portrait string // path to image for invetory
}

func New(id int, anim common.Animater, rect pixel.Rect, opts ...Option) *Loot {
	a := &Loot{
		id:      id,
		anim:    anim,
		rect:    rect,
		mass:    1,
		dir:     1,
		animdir: 1,
		vel:     pixel.ZV,
		sounds:  make(map[string]soundeffect),
		sprite:  pixel.NewSprite(nil, pixel.Rect{}),
	}

	for _, opt := range opts {
		opt(a)
	}

	p := common.NewPhys(rect,

		common.WithGravity(a.grav),
		common.WithMass(a.mass),
		common.WithRigidityBottom(0.3),
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
	a.phys.Apply(a.appliedForce)
	a.phys.Update(dt, objs)

	a.vec = pixel.ZV
	newspeed := a.phys.GetVel()
	a.vel = newspeed

	a.rect = a.phys.GetRect()

	a.counter += dt
	a.animSpriteNum = int(math.Floor(a.counter / 0.2))
	a.appliedForce = pixel.ZV
}

func (a *Loot) GetDir() int {
	return int(a.dir)
}

func (l *Loot) GetImagePath() string {
	return l.portrait
}

func (l *Loot) Draw(t pixel.Target) {
	pic, rect := l.anim.GetSprite("idle", l.animSpriteNum)
	l.sprite.Set(pic, rect)
	drawrect := l.rect
	l.sprite.Draw(t, pixel.IM.
		ScaledXY(pixel.ZV, pixel.V(
			drawrect.W()/l.sprite.Frame().W(),
			drawrect.H()/l.sprite.Frame().H(),
		)).
		Moved(drawrect.Center()),
	)
}

func (l *Loot) AddSound(event string) {
	if s, ok := l.sounds[event]; ok {
		// select random sound
		i := int(math.Round(common.GetRandFloat() * float64(len(s.List)-1)))
		sound.AddEffect(s.List[i], l.rect.Center())
	}
}

func (a *Loot) IsGround() bool {
	return a.phys.IsGround()
}

func (a *Loot) GetId() int {
	return a.id
}
