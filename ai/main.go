package ai

import (
	"platformer/events"
)

var list []*Ai

type Ai struct {
	obj    Manageder
	w      Worlder
	state  Stater
	states map[int]Stater
}

func New(obj Manageder, w Worlder) *Ai {
	a := &Ai{
		obj:    obj,
		w:      w,
		states: make(map[int]Stater),
	}
	a.initStates()
	a.Notify(events.SHIFT)
	list = append(list, a)
	return a
}

func (ai *Ai) initStates() {
	sIdle := NewIdle(ai, ai.w)
	ai.states[IDLE] = sIdle

	sAttack := NewAttack(ai, ai.w)
	ai.states[ATTACK] = sAttack

	ai.SetState(IDLE)
}

func (ai *Ai) SetState(state int) {
	ai.state = ai.states[state]
	ai.state.Start()
}

func Update() {
	for _, a := range list {
		a.Update()
	}
}

func (a *Ai) Update() {
	a.state.Update(0.0)
}

func (a *Ai) Notify(e int) {
	a.obj.Notify(e, a.state.GetVec())
}

// func (ai *Ai) SetGround(g bool) {
// 	ai.ground = g
// }

// func (ai *Ai) Update(target pixel.Vec, objs []Objecter) {
// 	ai.cmd = NOACTION
// 	ai.vec = pixel.ZV

// 	pos := ai.pers.rect.Center()
// 	move := 0.0
// 	if math.Abs(pos.X-target.X) > 5 {
// 		if math.Signbit(pos.X - target.X) {
// 			move = 1.0
// 		} else {
// 			move = -1.0
// 		}
// 	}
// 	ai.vec.X += move

// 	if ai.pers.rect.Contains(target) {
// 		ai.cmd = STRIKE
// 	}
// }
