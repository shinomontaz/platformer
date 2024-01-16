package objects

import (
	"math"
	"platformer/animation"
	"platformer/common"

	"github.com/shinomontaz/pixel"
)

type object struct {
	anim          common.Animater
	animSpriteNum int
	counter       float64
	rect          pixel.Rect
	phys          common.Phys
	vel           pixel.Vec
	animdir       float64
	dir           float64
}

// type List struct {
// 	objs []*object
// }

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

	o.phys = common.NewPhys(o.rect,
		common.WithMass(0),
		common.WithRigidity(1),
		common.WithFriction(0),
	)

	force := pixel.Vec{0, 0}
	if ai == "floater" {
		force.X = -5000
	}

	o.phys.Apply(force)
	objs = append(objs, o)
}

func List(r pixel.Rect) []common.Objecter {
	res := make([]common.Objecter, 0)
	for i, o := range objs {
		//		if r.Intersects(o.rect) {
		res = append(res, common.Objecter{ID: uint32(i + 1), R: o.rect, Type: common.OBJECT})
		//		}
	}
	return res
}

func GetPhysById(id uint32) *common.Phys {
	idx := int(id - 1) // REASON: we add objects with i+1 to list of Objecter ( look List method )
	if idx >= 0 && idx < len(objs)-1 {
		return &objs[idx].phys
	}
	return nil
}

func Update(dt float64, qt *common.Quadtree) {
	for _, o := range objs {
		vPhys := qt.Retrieve(o.rect)
		o.update(dt, vPhys)
	}
}

func Draw(win pixel.Target) {
	for _, o := range objs {
		o.draw(win)
	}
}

func (o *object) update(dt float64, vPhys []common.Objecter) {
	o.phys.Update(dt, vPhys)
	o.rect = o.phys.GetRect()

	o.vel = o.phys.GetVel()
	//	fmt.Println(o.vel, math.Signbit(o.vel.X), o.dir)
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
