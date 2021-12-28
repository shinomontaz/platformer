package actor

import (
	"fmt"

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
	phys Physicer

	state  ActorStater
	states map[int]ActorStater

	counter float64
	rect    pixel.Rect

	anim   Animater
	sprite *pixel.Sprite
	dir    float64
	vec    *pixel.Vec
}

func New(phys Physicer, anim Animater) *Actor {
	a := &Actor{
		phys: phys,
		anim: anim,
		rect: phys.GetRect(),
		dir:  1,
	}

	fmt.Println("On hero created", phys.GetRect())

	sFree := NewFreeState(a, anim)

	sAttack := NewAttackState(a, anim)
	sDead := NewDeadState(a, anim)
	sHit := NewHitState(a, anim)

	a.states = map[int]ActorStater{STATE_FREE: sFree, STATE_ATTACK: sAttack, STATE_DEAD: sDead, STATE_HIT: sHit}
	a.state = sFree
	return a
}

func (a *Actor) GetId() int {
	return a.id
}

func (a *Actor) Notify(e int, v *pixel.Vec) {
	if v != nil && v.X != 0 {
		if v.X > 0 {
			a.dir = 1
		} else {
			a.dir = -1
		}
	}

	a.vec = v
	a.state.Notify(e, v)
}

func (a *Actor) GetPos() pixel.Vec {
	return a.rect.Center()
}

func (a *Actor) Update(dt float64) {
	if a.vec != nil {
		a.phys.Update(dt, *a.vec)
	}

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

	a.sprite.Draw(t, pixel.IM.
		ScaledXY(pixel.ZV, pixel.V(
			a.rect.W()/a.sprite.Frame().W(),
			a.rect.H()/a.sprite.Frame().H(),
		)).
		ScaledXY(pixel.ZV, pixel.V(a.dir, 1)).
		Moved(a.rect.Center()),
	)
}

func (a *Actor) Hit(pos, vec pixel.Vec, power int) {
}
