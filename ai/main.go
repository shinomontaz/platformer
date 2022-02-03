package ai

import (
	"github.com/faiface/pixel"
)

var list map[Manageder]*Ai

func init() {
	list = make(map[Manageder]*Ai)
}

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
	list[obj] = a

	//	list = append(list, a)
	return a
}

func (ai *Ai) initStates() {
	sIdle := NewIdle(ai, ai.w)
	ai.states[IDLE] = sIdle

	sAttack := NewAttack(ai, ai.w)
	ai.states[ATTACK] = sAttack

	// sInactive := NewInactive(ai, ai.w)
	// ai.states[INACTIVE] = sInactive

	sInvestigate := NewInvestigate(ai, ai.w)
	ai.states[INVESTIGATE] = sInvestigate

	ai.SetState(IDLE, pixel.ZV)
}

func GetByObj(obj Manageder) *Ai {
	if a, ok := list[obj]; ok {
		return a
	}
	return nil
}

func (ai *Ai) SetState(state int, poi pixel.Vec) {
	ai.state = ai.states[state]
	ai.state.Start(poi)
}

func (ai *Ai) IsAlerted() bool {
	return ai.state.IsAlerted()
}

func Update(dt float64) {
	for _, a := range list {
		a.Update(dt)
	}
}

func (a *Ai) Update(dt float64) {
	a.state.Update(dt)
}

func (a *Ai) GetPos() pixel.Vec {
	return a.obj.GetPos()
}

func (a *Ai) Notify(e int, v pixel.Vec) {
	a.state.Notify(e, v)
}
