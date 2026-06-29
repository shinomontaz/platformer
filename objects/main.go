package objects

import (
	"math"
	"platformer/animation"
	"platformer/common"
	"platformer/common/physics"

	"github.com/shinomontaz/pixel"
)

type object struct {
	anim          common.Animater
	animSpriteNum int
	counter       float64
	rect          pixel.Rect
	phys          physics.Phys
	vel           pixel.Vec
	animdir       float64
	dir           float64
}

var objs []*object

func Init(g float64) {
	objs = make([]*object, 0)
}

func Add(name string, rect pixel.Rect, ai string) {
	o := &object{
		anim:    animation.Get(name),
		rect:    rect,
		dir:     1,
		animdir: 1,
		vel:     pixel.ZV,
	}

	o.phys = physics.New(o.rect,
		physics.WithMass(0),
		physics.WithRigidity(1),
		physics.WithFriction(0),
	)

	force := pixel.ZV
	if ai == "floater" {
		force.X = -5000
	}

	o.phys.Apply(force)
	objs = append(objs, o)
}

func GetPhysById(id uint32) *physics.Phys {
	idx := int(id - 1)
	if idx >= 0 && idx < len(objs)-1 {
		return &objs[idx].phys
	}
	return nil
}

type objWrapper struct {
	id   int
	rect pixel.Rect
}

func (o objWrapper) Rect() pixel.Rect { return o.rect }
func (o objWrapper) Id() int          { return o.id }

func List() []common.O {
	res := make([]common.O, 0, len(objs))
	for i, o := range objs {
		res = append(res, objWrapper{id: i + 1, rect: o.rect})
	}
	return res
}

func Update(dt float64, visiblePhys []common.O) {
	for _, o := range objs {
		o.update(dt, visiblePhys)
	}
}

func Draw(win pixel.Target) {
	for _, o := range objs {
		o.draw(win)
	}
}

func (o *object) update(dt float64, vPhys []common.O) {
	o.phys.Update(dt, vPhys)
	o.rect = o.phys.GetRect()

	o.vel = o.phys.GetVel()
	if !math.Signbit(o.vel.X) {
		o.dir = 1
	} else {
		o.dir = -1
	}

	o.counter += dt
	o.animSpriteNum = int(math.Floor(o.counter / 0.15))
}

func (o *object) draw(t pixel.Target) {
	if o.anim == nil {
		return
	}
	sprite := pixel.NewSprite(nil, pixel.Rect{})
	pic, rect := o.anim.GetSprite("walk", o.animSpriteNum)
	sprite.Set(pic, rect)

	drawrect := o.rect
	sprite.Draw(t, pixel.IM.
		ScaledXY(pixel.ZV, pixel.V(
			drawrect.W()/sprite.Frame().W(),
			drawrect.H()/sprite.Frame().H(),
		)).
		ScaledXY(pixel.ZV, pixel.V(o.animdir*o.dir, 1)).
		Moved(drawrect.Center()),
	)
	//	o.phys.Draw(t)
}
