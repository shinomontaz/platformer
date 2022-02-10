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

func NewCommon(obj Manageder, w Worlder) *Ai {
	a := &Ai{
		obj:    obj,
		w:      w,
		states: make(map[int]Stater),
	}

	sIdle := NewIdle(a, a.w)
	a.states[IDLE] = sIdle

	sAttack := NewAttack(a, a.w)
	a.states[ATTACK] = sAttack

	sInvestigate := NewInvestigate(a, a.w)
	a.states[INVESTIGATE] = sInvestigate

	a.SetState(IDLE, pixel.ZV)

	//	a.initStates()
	list[obj] = a

	obj.SetAi(a)
	//	list = append(list, a)
	return a
}

func NewMage(obj Manageder, w Worlder) *Ai {
	a := &Ai{
		obj:    obj,
		w:      w,
		states: make(map[int]Stater),
	}

	sIdle := NewIdle(a, a.w)
	a.states[IDLE] = sIdle

	sAttack := NewCast(a, a.w)
	a.states[CAST] = sAttack

	sInvestigate := NewInvestigate(a, a.w)
	a.states[INVESTIGATE] = sInvestigate

	a.SetState(IDLE, pixel.ZV)

	//	a.initStates()
	list[obj] = a

	obj.SetAi(a)
	//	list = append(list, a)
	return a
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
