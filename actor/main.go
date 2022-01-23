package actor

import (
	"math"
	"platformer/events"

	"github.com/faiface/pixel"
)

const (
	WALKING = iota
	RUNNING
	JUMPING
	FALLING
	STANDING
	IDLE
	FIRING
	HURT
	DYING
	DEAD
)

const (
	NOACTION = iota
	HITTED
	STRIKE
)

type Actor struct {
	id   int
	phys Phys

	state  ActorStater
	states map[int]ActorStater

	rect pixel.Rect

	anim      Animater
	sprite    *pixel.Sprite
	dir       float64
	vec       pixel.Vec // delta speed
	vel       pixel.Vec
	runspeed  float64
	walkspeed float64
	isShift   bool
}

func New(w Worlder, anim Animater, rect pixel.Rect, run, walk float64) *Actor {
	grav := w.GetGravity()
	p := NewPhys(rect, run, walk, grav*40, grav)
	p.SetQt(w.GetQt())

	a := &Actor{
		phys:      p,
		anim:      anim,
		rect:      rect,
		dir:       1,
		runspeed:  run,
		walkspeed: walk,
		vel:       pixel.ZV,
	}

	// init states
	sFree := NewFreeState(a, anim)
	sAttack := NewAttackState(a, anim)
	sDead := NewDeadState(a, anim)
	sHit := NewHitState(a, anim)

	a.states = map[int]ActorStater{STATE_FREE: sFree, STATE_ATTACK: sAttack, STATE_DEAD: sDead, STATE_HIT: sHit}
	a.state = sFree
	return a
}

func (a *Actor) Move(vec pixel.Vec) {
	a.phys.Move(vec)
	a.rect = a.phys.GetRect()
}

func (a *Actor) GetId() int {
	return a.id
}

func (a *Actor) Notify(e int, v pixel.Vec) {
	if v.X != 0 {
		if v.X > 0 {
			a.dir = 1
		} else {
			a.dir = -1
		}
		if !a.isShift {
			v.X *= 2
		}
	}

	if e == events.SHIFT {
		a.isShift = !a.isShift
	}

	a.state.Notify(e, &v)
	a.vec = v
}

func (a *Actor) GetPos() pixel.Vec {
	return a.rect.Center()
}

func (a *Actor) Update(dt float64) {
	a.phys.Update(dt, &a.vec)
	newspeed := a.phys.GetVel()
	var event int
	if math.Abs(newspeed.X) <= a.runspeed && math.Abs(newspeed.X) > a.walkspeed {
		event = events.RUN
	} else if (math.Abs(a.vel.X) > a.walkspeed && math.Abs(newspeed.X) <= a.walkspeed) || (a.vel.X == 0 && math.Abs(newspeed.X) > 0 && math.Abs(newspeed.X) <= a.walkspeed) {
		event = events.WALK
	}
	a.state.Notify(event, newspeed)
	a.vel = *newspeed

	a.rect = a.phys.GetRect()
	a.state.Update(dt)
}

// func (h *Actor) selectAttack(move pixel.Vec) string {
// 	if move.Len() == 0 {
// 		return fmt.Sprintf("attack%d", rand.Intn(2)+1)
// 	}
// 	return "attack3"
// }

func (a *Actor) SetState(id int) {
	a.state = a.states[id]
	a.state.Start()
}

func (a *Actor) Draw(t pixel.Target) {
	a.sprite = a.state.GetSprite()
	drawrect := a.rect.ResizedMin(pixel.Vec{a.rect.W() * 1.25, a.rect.H() * 1.25})
	a.sprite.Draw(t, pixel.IM.
		ScaledXY(pixel.ZV, pixel.V(
			drawrect.W()/a.sprite.Frame().W(),
			drawrect.H()/a.sprite.Frame().H(),
		)).
		ScaledXY(pixel.ZV, pixel.V(a.dir, 1)).
		Moved(drawrect.Center()),
	)
	//	a.phys.Draw(t)
}

func (a *Actor) Hit(pos, vec pixel.Vec, power int) {
}
