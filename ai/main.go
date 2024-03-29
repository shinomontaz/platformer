package ai

import (
	"platformer/actor"

	"github.com/shinomontaz/pixel"
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

func NewCalmEnemy(obj *actor.Actor, w Worlder) *Ai {
	a := &Ai{
		obj:    obj,
		w:      w,
		states: make(map[int]Stater),
		id:     counter,
	}

	sIdle := NewIdle(a, a.w, true)
	a.states[IDLE] = sIdle

	sChooseAttack := NewChooseAttack(a, a.w)
	a.states[CHOOSEATTACK] = sChooseAttack

	sAttack := NewAttack(a, a.w)
	a.states[ATTACK] = sAttack

	sBustle := NewBustle(a, a.w)
	a.states[BUSTLE] = sBustle

	sInvestigate := NewInvestigate(a, a.w)
	a.states[INVESTIGATE] = sInvestigate

	a.SetState(IDLE, pixel.ZV)

	list[obj] = a
	obj.AddEventListener(a)

	counter++
	return a
}

func NewAgressiveEnemy(obj *actor.Actor, w Worlder) *Ai {
	a := &Ai{
		obj:    obj,
		w:      w,
		states: make(map[int]Stater),
		id:     counter,
	}

	sRoaming := NewRoaming(a, a.w, true)
	a.states[ROAMING] = sRoaming

	sChooseAttack := NewChooseAttack(a, a.w)
	a.states[CHOOSEATTACK] = sChooseAttack

	sAttack := NewAttack(a, a.w)
	a.states[ATTACK] = sAttack

	sBustle := NewBustle(a, a.w)
	a.states[BUSTLE] = sBustle

	sIdle := NewIdle(a, a.w, true)
	a.states[IDLE] = sIdle

	sInvestigate := NewInvestigate(a, a.w)
	a.states[INVESTIGATE] = sInvestigate

	hero := obj.GetEnemy()
	if hero != nil {
		a.SetState(INVESTIGATE, hero.GetPos())
	} else {
		a.SetState(IDLE, pixel.ZV)
	}

	list[obj] = a
	obj.AddEventListener(a)

	counter++
	return a
}

func NewActiveEnemy(obj *actor.Actor, w Worlder) *Ai {
	a := &Ai{
		obj:    obj,
		w:      w,
		states: make(map[int]Stater),
		id:     counter,
	}

	sRoaming := NewRoaming(a, a.w, true)
	a.states[ROAMING] = sRoaming

	sChooseAttack := NewChooseAttack(a, a.w)
	a.states[CHOOSEATTACK] = sChooseAttack

	sAttack := NewAttack(a, a.w)
	a.states[ATTACK] = sAttack

	sBustle := NewBustle(a, a.w)
	a.states[BUSTLE] = sBustle

	sIdle := NewIdle(a, a.w, true)
	a.states[IDLE] = sIdle

	sInvestigate := NewInvestigate(a, a.w)
	a.states[INVESTIGATE] = sInvestigate

	a.SetState(ROAMING, pixel.ZV)

	list[obj] = a
	obj.AddEventListener(a)

	counter++
	return a
}

func NewCalmNpc(obj *actor.Actor, w Worlder) *Ai {
	a := &Ai{
		obj:    obj,
		w:      w,
		states: make(map[int]Stater),
		id:     counter,
	}

	sIdle := NewIdle(a, a.w, false)
	a.states[IDLE] = sIdle

	a.SetState(IDLE, pixel.ZV)

	list[obj] = a
	obj.AddEventListener(a)

	counter++
	return a
}

func NewFishingNpc(obj *actor.Actor, w Worlder) *Ai {
	a := &Ai{
		obj:    obj,
		w:      w,
		states: make(map[int]Stater),
		id:     counter,
	}

	sFishing := NewFishing(a, a.w, false)
	a.states[FISHING] = sFishing

	a.SetState(FISHING, pixel.ZV)

	list[obj] = a
	obj.AddEventListener(a)

	counter++
	return a
}

func NewActiveNpc(obj *actor.Actor, w Worlder) *Ai {
	a := &Ai{
		obj:    obj,
		w:      w,
		states: make(map[int]Stater),
		id:     counter,
	}

	sRoaming := NewRoaming(a, a.w, false)
	a.states[ROAMING] = sRoaming

	sInvestigate := NewInvestigate(a, a.w)
	a.states[INVESTIGATE] = sInvestigate

	a.SetState(ROAMING, pixel.ZV)

	list[obj] = a
	obj.AddEventListener(a)

	counter++
	return a
}

func NewSwimmingNpc(obj *actor.Actor, w Worlder) *Ai {
	a := &Ai{
		obj:    obj,
		w:      w,
		states: make(map[int]Stater),
		id:     counter,
	}

	sSwimming := NewSwimming(a, a.w, false)
	a.states[SWIMMING] = sSwimming

	a.SetState(SWIMMING, pixel.ZV)

	list[obj] = a
	obj.AddEventListener(a)

	counter++
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

func (a *Ai) EventAction(e int) {
	a.state.EventAction(e)
}
