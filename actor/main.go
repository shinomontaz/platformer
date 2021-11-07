package actor

import (
	"github.com/faiface/pixel"
)

const (
	WALKING = iota
	RUNNING
	JUMPING
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

	anim Animater
	// sheet  pixel.Picture
	// frame  pixel.Rect
	sprite *pixel.Sprite
	// anims  map[string]Animater
	dir float64
}

func New(phys Physicer, anim Animater) *Actor {
	a := &Actor{
		phys: phys,
		anim: anim,
	}

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
	a.state.Notify(e, v)
}

func (a *Actor) GetPos() pixel.Vec {
	return a.rect.Center()
}

func (a *Actor) Update(dt float64) {
	a.state.Update(dt)
	//	a.counter += dt

	// var newState int
	// switch {
	// case !h.phys.ground:
	// 	newState = JUMPING
	// case cmd == STRIKE:
	// 	newState = FIRING
	// case cmd == HITTED:
	// 	newState = HURT
	// case h.phys.vel.Len() == 0:
	// 	newState = STANDING
	// case h.phys.vel.Len() == h.phys.walkSpeed:
	// 	newState = WALKING
	// case h.phys.vel.Len() == h.phys.runSpeed:
	// 	newState = RUNNING
	// }

	// if h.state == IDLE { // make idle animation
	// 	newState = IDLE
	// }
	// if h.counter > 5.0 && h.state == STANDING && (newState == h.state || newState == STANDING) { // make idle animation
	// 	newState = IDLE
	// }

	// // make state transision
	// if h.state != newState {
	// 	h.state = newState
	// 	h.counter = 0
	// }

	// switch a.state.GetId() {
	// case STATE_FREE:
	// 	i := int(math.Floor(h.counter / 0.1)) // MAGIC CONST!

	// 	if i > len(h.anims["idle"].frames) {
	// 		h.state = STANDING
	// 	}

	// 	h.frame = h.anims["idle"].frames[i%len(h.anims["idle"].frames)]
	// 	h.sheet = h.anims["idle"].sheet
	// case STANDING:
	// 	h.frame = h.anims["stand"].frames[0]
	// 	h.sheet = h.anims["stand"].sheet
	// case WALKING:
	// 	i := int(math.Floor(h.counter / 0.1)) // MAGIC CONST!
	// 	h.frame = h.anims["walk"].frames[i%len(h.anims["walk"].frames)]
	// 	h.sheet = h.anims["walk"].sheet
	// case RUNNING:
	// 	i := int(math.Floor(h.counter / 0.1)) // MAGIC CONST!
	// 	h.frame = h.anims["run"].frames[i%len(h.anims["run"].frames)]
	// 	h.sheet = h.anims["run"].sheet
	// case JUMPING:
	// 	speed := h.phys.vel.Y
	// 	i := int((-speed/h.phys.jumpSpeed + 1) / 2 * float64(len(h.anims["jump"].frames)))
	// 	if i < 0 {
	// 		i = 0
	// 	}
	// 	if i >= len(h.anims["jump"].frames) {
	// 		i = len(h.anims["jump"].frames) - 1
	// 	}
	// 	h.frame = h.anims["jump"].frames[i]
	// 	h.sheet = h.anims["jump"].sheet
	// case STATE_HIT:
	// 	h.frame = h.anims["hurt"].frames[1] // only second frame we get!
	// 	h.sheet = h.anims["hurt"].sheet
	// case STATE_ATTACK:
	// 	i := int(math.Floor(h.counter / 0.1)) // h.counter stands for frame rate change in animation
	// 	// if i == 0 || i%len(h.anims[h.attackCode].frames) == 0 {
	// 	// 	h.attackCode = h.selectAttack(h.phys.vel)
	// 	// }
	// 	h.frame = h.anims[h.attackCode].frames[i%len(h.anims[h.attackCode].frames)]
	// 	h.sheet = h.anims[h.attackCode].sheet
	// }

	// if h.phys.vel.X != 0 {
	// 	if h.phys.vel.X > 0 {
	// 		h.dir = +1
	// 	} else {
	// 		h.dir = -1
	// 	}
	// }

	// h.pos = h.phys.rect.Min
	// h.rect = pixel.R(h.phys.rect.Center().X-h.phys.rect.W()/2, h.phys.rect.Center().Y-h.phys.rect.H()/2, h.phys.rect.Center().X-h.phys.rect.W()/2+h.rect.W(), h.phys.rect.Center().Y-h.phys.rect.H()/2+h.rect.H())
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
	// if a.sprite == nil {
	// 	a.sprite = pixel.NewSprite(nil, pixel.Rect{})
	// }
	// draw the correct frame with the correct position and direction
	a.sprite = a.state.GetSprite()
	//	a.sprite.Set(a.sheet, a.frame)
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
