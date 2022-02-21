package magic

import (
	"fmt"
	"math"
	"math/rand"
	"platformer/actor"
	"platformer/animation"
	"platformer/common"

	"github.com/faiface/pixel"
	"github.com/faiface/pixel/text"
	"golang.org/x/image/colornames"
	"golang.org/x/image/font"
)

type Deathstrike struct {
	anim          Animater
	sprite        *pixel.Sprite
	animSpriteNum int
	vel           pixel.Vec
	w             Worlder
	power         int
	owner         *actor.Actor
	target        pixel.Vec
	ttl           float64
	timer         float64
	prepare_time  float64
	rect          pixel.Rect
	qt            *common.Quadtree
	atlas         *text.Atlas
	fnt           font.Face
}

func NewDeathstrike(s string, owner *actor.Actor, target pixel.Vec) *Deathstrike {
	basicrect := pixel.R(0, 0, 48, 64)

	ds := Deathstrike{
		anim:         animation.Get(s),
		w:            w,
		power:        1,
		owner:        owner,
		target:       target,
		ttl:          3.0,
		prepare_time: 1.0,
		rect:         basicrect,
		qt:           w.GetQt(),
	}

	ds.fnt = common.GetFont("regular")
	ds.atlas = text.NewAtlas(ds.fnt, text.ASCII)

	ds.init()

	return &ds
}

func (d *Deathstrike) init() {
	// get target floor
	// move rect to target

	broadbox := pixel.R(d.target.X-d.rect.W()/2, d.target.Y-2*d.rect.H(), d.target.X+d.rect.W()/2, d.target.Y+d.rect.H())
	objs := d.qt.Retrieve(broadbox)

	groundY := objs[0].Rect().Max.Y
	for _, o := range objs {
		y := o.Rect().Max.Y
		x1, x2 := o.Rect().Min.X, o.Rect().Max.X
		if y > groundY && d.target.X <= x2 && d.target.X >= x1 {
			groundY = y
		}
	}

	d.rect = pixel.R(d.target.X-d.rect.W()/2, groundY, d.target.X+d.rect.W()/2, groundY+d.rect.H())

}
func (d *Deathstrike) GetOwner() *actor.Actor {
	return d.owner
}

func (d *Deathstrike) Draw(t pixel.Target) {
	if d.timer < d.prepare_time {
		owrect := d.owner.GetRect()
		rect := pixel.R(owrect.Center().X-20, owrect.Max.Y, owrect.Center().X, owrect.Max.Y+30)
		ch := randSeq([]rune("#$%&@*?arlTVXx"), 2)

		pos := pixel.V(rect.Min.X, rect.Min.Y)
		// draw exclamation sign
		txt := text.New(pos, d.atlas)
		txt.Color = colornames.Grey
		fmt.Fprintln(txt, ch)
		txt.Draw(t, pixel.IM)

		return
	}

	if d.sprite == nil {
		d.sprite = pixel.NewSprite(nil, pixel.Rect{})
	}
	pic, rect := d.anim.GetSprite("spell", d.animSpriteNum)
	d.sprite.Set(pic, rect)

	drawrect := d.rect.ResizedMin(pixel.Vec{d.rect.W(), d.rect.H()})
	d.sprite.Draw(t, pixel.IM.
		ScaledXY(pixel.ZV, pixel.V(
			drawrect.W()/d.sprite.Frame().W(),
			drawrect.H()/d.sprite.Frame().H(),
		)).Moved(drawrect.Center()),
	)
}

func randSeq(letters []rune, n int) string {
	b := make([]rune, n)
	for i := range b {
		b[i] = letters[rand.Intn(len(letters))]
	}
	return string(b)
}

func (d *Deathstrike) Update(dt float64) {
	d.timer += dt
	if d.timer > d.prepare_time {
		d.animSpriteNum = int(math.Floor(d.timer / 0.1))
	}
}

// func (d *Deathstrike) Apply(a Actor) {

// }

func (d *Deathstrike) IsFinished() bool {
	return d.timer > d.ttl
}
