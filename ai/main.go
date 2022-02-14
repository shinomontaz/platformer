package ai

import (
	"platformer/actor"

	"github.com/faiface/pixel"
)

var list map[*actor.Actor]*Ai
var counter int

func init() {
	list = make(map[*actor.Actor]*Ai)
}

type Ai struct {
	obj         *actor.Actor
	w           Worlder
	state       Stater
	states      map[int]Stater
	id          int
	attackskill *actor.Skill
}

func NewCommon(obj *actor.Actor, w Worlder) *Ai {
	a := &Ai{
		obj:    obj,
		w:      w,
		states: make(map[int]Stater),
		id:     counter,
	}

	sIdle := NewIdle(a, a.w)
	a.states[IDLE] = sIdle

	sChooseAttack := NewChooseAttack(a, a.w)
	a.states[CHOOSEATTACK] = sChooseAttack

	sAttack := NewAttack(a, a.w)
	a.states[ATTACK] = sAttack

	sInvestigate := NewInvestigate(a, a.w)
	a.states[INVESTIGATE] = sInvestigate

	a.SetState(IDLE, pixel.ZV)

	list[obj] = a
	obj.Subscribe(a)

	counter++
	return a
}

func NewMage(obj *actor.Actor, w Worlder) *Ai {
	a := &Ai{
		obj:    obj,
		w:      w,
		states: make(map[int]Stater),
		id:     counter,
	}

	sIdle := NewIdle(a, a.w)
	a.states[IDLE] = sIdle

	sInvestigate := NewInvestigate(a, a.w)
	a.states[INVESTIGATE] = sInvestigate

	a.SetState(IDLE, pixel.ZV)

	//	a.initStates()
	list[obj] = a
	obj.Subscribe(a)

	counter++
	//	list = append(list, a)
	return a
}

func GetByObj(obj *actor.Actor) *Ai {
	if a, ok := list[obj]; ok {
		return a
	}
	return nil
}

func (ai *Ai) GetId() int {
	return ai.id
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
