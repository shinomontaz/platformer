package ai

import (
	"math"
	"platformer/actor"
	"platformer/bindings"
	"platformer/events"

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
	states      map[AiState]Stater
	id          int
	attackskill *actor.Skill
	homeState   AiState // состояние куда возвращаться после INVESTIGATE (IDLE, PATROL, ROAMING)
	name        string

	maxReach float64 // максимальный зазор, на котором есть хоть один навык
	minGap   float64 // минимальный Min среди дальних навыков
	hasMelee bool

	debug     bool
	debugTick float64
}

func parseReaction(s string) ReactionType {
	switch s {
	case "flee":
		return HURTED_FLEE
	case "counter":
		return HURTED_COUNTER
	case "random":
		return HURTED_RANDOM
	default:
		return HURTED_NONE
	}
}

func NewCalmEnemy(obj *actor.Actor, w Worlder) *Ai {
	a := &Ai{
		obj:       obj,
		w:         w,
		states:    make(map[AiState]Stater),
		id:        counter,
		homeState: IDLE,
	}

	a.initReach()

	sIdle := NewIdle(a, a.w, true)
	a.states[IDLE] = sIdle

	sChooseAttack := NewChooseAttack(a, a.w)
	a.states[CHOOSEATTACK] = sChooseAttack

	sAttack := NewAttack(a, a.w)
	a.states[ATTACK] = sAttack

	reaction := parseReaction(obj.GetHurtReaction())
	sHurted := NewHurted(a, a.w, reaction)
	a.states[HURTED] = sHurted

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
		obj:       obj,
		w:         w,
		states:    make(map[AiState]Stater),
		id:        counter,
		homeState: IDLE,
	}

	a.initReach()

	sRoaming := NewRoaming(a, a.w, true)
	a.states[ROAMING] = sRoaming

	sChooseAttack := NewChooseAttack(a, a.w)
	a.states[CHOOSEATTACK] = sChooseAttack

	sAttack := NewAttack(a, a.w)
	a.states[ATTACK] = sAttack

	sIdle := NewIdle(a, a.w, true)
	a.states[IDLE] = sIdle

	reaction := parseReaction(obj.GetHurtReaction())
	sHurted := NewHurted(a, a.w, reaction)
	a.states[HURTED] = sHurted

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
		obj:       obj,
		w:         w,
		states:    make(map[AiState]Stater),
		id:        counter,
		homeState: ROAMING,
	}

	a.initReach()

	sRoaming := NewRoaming(a, a.w, true)
	a.states[ROAMING] = sRoaming

	sChooseAttack := NewChooseAttack(a, a.w)
	a.states[CHOOSEATTACK] = sChooseAttack

	sAttack := NewAttack(a, a.w)
	a.states[ATTACK] = sAttack

	sIdle := NewIdle(a, a.w, true)
	a.states[IDLE] = sIdle

	reaction := parseReaction(obj.GetHurtReaction())
	sHurted := NewHurted(a, a.w, reaction)
	a.states[HURTED] = sHurted

	sInvestigate := NewInvestigate(a, a.w)
	a.states[INVESTIGATE] = sInvestigate

	a.SetState(ROAMING, pixel.ZV)

	list[obj] = a
	obj.AddEventListener(a)

	counter++
	return a
}

func NewPatrolEnemy(obj *actor.Actor, w Worlder, patrolRange float64) *Ai {
	a := &Ai{
		obj:       obj,
		w:         w,
		states:    make(map[AiState]Stater),
		id:        counter,
		homeState: PATROL,
	}

	a.initReach()

	sPatrol := NewPatrol(a, a.w, patrolRange)
	a.states[PATROL] = sPatrol

	sChooseAttack := NewChooseAttack(a, a.w)
	a.states[CHOOSEATTACK] = sChooseAttack

	sAttack := NewAttack(a, a.w)
	a.states[ATTACK] = sAttack

	reaction := parseReaction(obj.GetHurtReaction())
	sHurted := NewHurted(a, a.w, reaction)
	a.states[HURTED] = sHurted

	sInvestigate := NewInvestigate(a, a.w)
	a.states[INVESTIGATE] = sInvestigate

	a.SetState(PATROL, pixel.ZV)

	list[obj] = a
	obj.AddEventListener(a)

	counter++
	return a
}

func NewCalmNpc(obj *actor.Actor, w Worlder) *Ai {
	a := &Ai{
		obj:    obj,
		w:      w,
		states: make(map[AiState]Stater),
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
		states: make(map[AiState]Stater),
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
		states: make(map[AiState]Stater),
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
		states: make(map[AiState]Stater),
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

func (ai *Ai) SetState(state AiState, poi pixel.Vec) {
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
	if e == events.HIT {
		if _, ok := a.states[HURTED]; ok {
			a.SetState(HURTED, pixel.ZV)

			return
		}
	}

	a.state.EventAction(e)
}

type target struct {
	hero   *actor.Actor
	hr     pixel.Rect
	dx, dy float64 // от центра к центру; dx > 0 — герой справа
	gap    float64 // зазор между краями по X, 0 при перекрытии
}

func sign(v float64) float64 {
	if v < 0 {
		return -1
	}
	return 1
}

func isMelee(t string) bool {
	return t == "melee" || t == "meleedown" || t == "meleejump"
}

func (a *Ai) initReach() {
	halfW := a.obj.GetRect().W() / 2
	for _, sk := range a.obj.GetSkills() {
		if sk.Name == "ancientbolt" {
			a.debug = true
		}

		if sk.Weight <= 0 {
			continue
		}
		if isMelee(sk.Type) {
			a.hasMelee = true
			a.maxReach = math.Max(a.maxReach, sk.Hitbox.Max.X-halfW)
			continue
		}
		a.maxReach = math.Max(a.maxReach, sk.Max)
		if a.minGap == 0 || sk.Min < a.minGap {
			a.minGap = sk.Min
		}
	}
}

func (a *Ai) look() (target, bool) {
	hero := a.obj.GetEnemy()
	if hero == nil || hero.GetHp() <= 0 {
		return target{}, false
	}
	r, hr := a.obj.GetRect(), hero.GetRect()
	dx := hr.Center().X - r.Center().X
	return target{
		hero: hero,
		hr:   hr,
		dx:   dx,
		dy:   hr.Center().Y - r.Center().Y,
		gap:  math.Max(0, math.Abs(dx)-(r.W()+hr.W())/2),
	}, true
}

func (a *Ai) facing(t target) bool {
	return t.dx*float64(a.obj.GetDir()) > 0
}

func (a *Ai) sees(t target) bool {
	return a.w.IsSee(a.obj.GetRect().Center(), t.hr.Center())
}

func (a *Ai) canUse(sk *actor.Skill, t target) bool {
	if d := float64(sk.GetDir()); d*t.dx < 0 {
		return false // направленный навык смотрит не туда
	}
	h := a.obj.GetRect().H()
	switch {
	case isMelee(sk.Type):
		return a.obj.StrikeRect(sk, sign(t.dx)).Intersects(t.hr)
	case sk.Type == "ranged":
		return t.gap >= sk.Min && t.gap <= sk.Max && math.Abs(t.dy) <= h
	default: // meleemove и прочее
		return t.gap >= sk.Min && t.gap <= sk.Max && math.Abs(t.dy) <= h/2
	}
}

func (a *Ai) step(dir float64) {
	if a.obj.StepPrediction(events.WALK, pixel.V(dir, 0)) <= 0.5 {
		return
	}
	key := bindings.KeyAction["right"]
	if dir < 0 {
		key = bindings.KeyAction["left"]
	}
	a.obj.KeyAction(bindings.Active.GetBinding(key))
}
