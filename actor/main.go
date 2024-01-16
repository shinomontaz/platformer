package actor

import (
	"errors"
	"fmt"
	"math"
	"platformer/bindings"
	"platformer/common"
	"platformer/events"
	"platformer/objects"
	"platformer/projectiles"
	"platformer/sound"

	"platformer/activities"

	"platformer/actor/state"
	"platformer/actor/statemachine"

	"github.com/shinomontaz/pixel"
	"github.com/shinomontaz/pixel/pixelgl"
)

const (
	NOACTION = iota
	HITTED
	STRIKE
	ALERTED
)

var counter int

type Actor struct {
	id        int
	enemy     *Actor // enemy
	phys      common.Phys
	groundObj common.Objecter

	action int
	state  Stater
	states map[int]Stater

	keycombo int // sum of all actionId ( from pressed keys )/ it is unique due to actionId const definition: powers of 2

	rect pixel.Rect

	animdir   float64
	anim      common.Animater
	sprite    *pixel.Sprite
	dir       float64
	vel       pixel.Vec
	runspeed  float64
	walkspeed float64
	grav      float64
	jumpforce float64
	mass      float64
	isShift   bool

	sm *statemachine.Machine
	w  Worlder

	hp           int
	strength     int
	portrait     *pixel.Sprite
	sounds       map[string]soundeffect
	appliedForce pixel.Vec
	attackrange  float64

	target pixel.Vec

	esbrs []common.EventSubscriber

	skills       []*Skill
	skillmap     map[int]int // map of keyaction sum (as ints) to skills (indices in skills list). we can do sum as unique key duw to skill id definition: powers of 2
	activeSkill  *Skill
	currObjs     []common.Objecter
	phrasesClass string
	score        int

	//	onkill     OnKillHandler
	onkill []OnKillHandler

	onhit            OnKillHandler
	oninteract       OnInteractHandler
	iswater          bool
	iswaterresistant bool
}

var loader *common.Loader

func Init(l *common.Loader) {
	loader = l
	counter = 0
}

func New(w Worlder, anim common.Animater, rect pixel.Rect, opts ...Option) *Actor {
	counter++
	a := &Actor{
		id:        counter,
		anim:      anim,
		rect:      rect,
		dir:       1,
		animdir:   1,
		vel:       pixel.ZV,
		grav:      w.GetGravity(),
		w:         w,
		mass:      1,
		sounds:    make(map[string]soundeffect),
		esbrs:     make([]common.EventSubscriber, 0),
		skills:    make([]*Skill, 0),
		skillmap:  make(map[int]int),
		groundObj: common.Objecter{},
		onkill:    make([]OnKillHandler, 0),
	}

	for _, opt := range opts { // skills inited here as well
		opt(a)
	}

	// we have skills list => we need a skillmap
	for idx, sk := range a.skills {
		sum := 0
		for _, actionId := range sk.Keys {
			sum += actionId
		}
		a.skillmap[sum] = idx
	}

	p := common.NewPhys(rect,
		common.WithGravity(a.grav),
		common.WithMass(a.mass),
	) // TODO does we really need phys to know run and walk speeds?
	a.phys = p

	a.initStates()

	return a
}

func (a *Actor) initStates() {
	sStand := state.New(state.STAND, a, a.anim)
	sWalk := state.New(state.WALK, a, a.anim)
	sRun := state.New(state.RUN, a, a.anim)
	sIdle := state.New(state.IDLE, a, a.anim)
	sJump := state.New(state.JUMP, a, a.anim)
	sFall := state.New(state.FALL, a, a.anim)
	sDead := state.New(state.DEAD, a, a.anim)
	sHit := state.New(state.HIT, a, a.anim)
	sInteract := state.New(state.INTERACT, a, a.anim)
	sResurrect := state.New(state.RESURRECT, a, a.anim)
	sFishing := state.New(state.FISHING, a, a.anim)
	sSwim := state.New(state.SWIM, a, a.anim)

	a.states = map[int]Stater{
		state.STAND:     sStand,
		state.IDLE:      sIdle,
		state.WALK:      sWalk,
		state.RUN:       sRun,
		state.FALL:      sFall,
		state.JUMP:      sJump,
		state.HIT:       sHit,
		state.DEAD:      sDead,
		state.INTERACT:  sInteract,
		state.RESURRECT: sResurrect,
		state.FISHING:   sFishing,
		state.SWIM:      sSwim,
	}

	// apply state machine
	if a.sm != nil {
		for _, st := range a.sm.GetStates() {
			if _, ok := a.states[st]; !ok {
				sState := state.New(st, a, a.anim)
				if sState == nil {
					panic("Fail to create actor state!")
				}
				a.states[st] = sState
			}
		}
	}

	a.SetState(state.STAND)
}

func (a *Actor) GetTransition(state int) statemachine.Transition {
	if a.sm != nil {
		return a.sm.GetTransition(state)
	}
	return statemachine.Transition{}
}

func (a *Actor) StepPrediction(e int, v pixel.Vec) float64 {
	if a.iswater {
		return 1
	}

	if e == events.WALK {
		v.X *= a.walkspeed / 20
	} else {
		v.X *= a.runspeed / 20
	}

	groundrate, _, _, _ := a.phys.StepPrediction(v)
	return groundrate
}

func (a *Actor) KeyAction(key pixelgl.Button) {
	if a.state.Busy() {
		return
	}
	action := bindings.Active.GetAction(key) // get action id for this key
	a.keycombo += action

	isWalk := false
	if action == bindings.RIGHT {
		a.dir = 1
		isWalk = true
	} else if action == bindings.LEFT {
		a.dir = -1
		isWalk = true
	}

	if isWalk {
		a.appliedForce.X = a.dir * a.walkspeed
		if a.isShift {
			a.appliedForce.X = a.dir * a.runspeed
		}
	}

	if action == bindings.SHIFT {
		a.isShift = true
	}

	if action == bindings.UP && a.state.GetId() != state.JUMP {
		multiplier := 0.5
		if math.Abs(a.vel.X) > a.walkspeed || a.iswater {
			multiplier = 1
		}
		a.appliedForce.Y = a.grav * (a.jumpforce + float64(a.strength) - a.mass) * multiplier
	}

	if action == bindings.DOWN {
		fmt.Println("key down")
		a.appliedForce.Y = -a.grav * a.jumpforce * 10
	}
}

func (a *Actor) Move(v pixel.Vec) {
	a.rect = a.rect.Moved(v)
	a.phys.Move(v)
}

func (a *Actor) GetPos() pixel.Vec {
	return a.rect.Min
}

func (a *Actor) GetRect() pixel.Rect {
	return a.rect
}

func (a *Actor) GetVel() pixel.Vec {
	return a.vel
}

func (a *Actor) SetVel(v pixel.Vec) {
	a.vel = v
	a.phys.SetSpeed(a.vel)
}

func (a *Actor) Update(dt float64, objs []common.Objecter) {
	a.currObjs = objs
	if a.appliedForce.X == 0 && math.Abs(a.vel.X) > a.walkspeed*0.1 && a.state.GetId() != state.MELEEMOVE { // active slowing in case of no active force applied
		a.appliedForce.X = -a.walkspeed
		if a.vel.X < 0 {
			a.appliedForce.X *= -1
		}
	}

	if math.Signbit(a.appliedForce.X) != math.Signbit(a.vel.X) { // to brake (тормозить) harder
		a.appliedForce.X *= 3
	} else {
		if math.Abs(a.appliedForce.X) == a.walkspeed && math.Abs(a.vel.X) >= a.walkspeed*0.9 {
			a.appliedForce.X = 0
		}
		if math.Abs(a.appliedForce.X) == a.runspeed && math.Abs(a.vel.X) >= a.runspeed*0.9 {
			a.appliedForce.X = 0
		}
	}

	a.phys.Apply(a.appliedForce)
	a.phys.Update(dt, objs)
	gdObj := a.phys.GetGroundObject()
	if gdObj.ID != a.groundObj.ID {
		a.groundObj = gdObj
		a.phys.SetGroundPhys(objects.GetPhysById(gdObj.ID))
	}

	newspeed := a.phys.GetVel()
	if math.Abs(newspeed.X) <= a.runspeed && math.Abs(newspeed.X) > a.walkspeed {
		a.action = events.RUN
	} else if math.Abs(newspeed.X) > 0 && math.Abs(newspeed.X) <= a.walkspeed {
		a.action = events.WALK
	}
	a.vel = newspeed

	if idx, ok := a.skillmap[a.keycombo]; ok {
		a.SetSkill(a.skills[idx])
		a.action = a.activeSkill.Event
	} else if a.keycombo == bindings.ENTER {
		a.action = events.INTERACT
	}

	a.state.Listen(a.action, &a.vel)
	a.state.Update(dt)

	a.rect = a.phys.GetRect()
	a.appliedForce = pixel.ZV
	a.keycombo = 0
	a.action = 0
	a.isShift = false
}

func (a *Actor) UpdateSpecial(dt float64, objs []common.Objecter) {
	force := a.appliedForce

	for _, o := range objs {
		isIntercects := a.rect.Intersects(o.R)
		if !isIntercects { // no collision
			continue
		}

		part := a.rect.Intersect(o.R)

		if o.Type == common.WATER {
			ratio := part.H() / a.rect.H()
			if ratio == 1 && !a.iswater {
				a.iswater = true
				a.state.SetWater(true)
				if !a.iswaterresistant && a.hp > 0 {
					a.Kill()
				} else {
					a.SetState(state.SWIM)
				}
			}

			force.Y += a.grav * (ratio + 0.1) // archimedus
			a.phys.SetSpeed(a.vel.Scaled(0.92))
		}
	}
	a.phys.Apply(force)
}

func (a *Actor) SetState(id int) {
	a.state = a.states[id]
	a.state.SetWaterResistant(a.iswaterresistant)
	a.state.Start()
}

func (a *Actor) GetDir() int {
	return int(a.dir)
}

func (a *Actor) SetDir(d float64) {
	a.dir = d
}

func (a *Actor) Draw(t pixel.Target) {
	a.sprite = a.state.GetSprite()
	drawrect := a.rect
	a.sprite.Draw(t, pixel.IM.
		ScaledXY(pixel.ZV, pixel.V(
			drawrect.W()/a.sprite.Frame().W(),
			drawrect.H()/a.sprite.Frame().H(),
		)).
		ScaledXY(pixel.ZV, pixel.V(a.animdir*a.dir, 1)). //
		Moved(drawrect.Center()),
	)
	//	a.phys.Draw(t)
}

func (a *Actor) GetSkills() []*Skill {
	return a.skills
}

func (a *Actor) GetSkillAttr(attr string) (interface{}, error) {
	if a.activeSkill == nil {
		return "", errors.New("no active skill")
	}

	switch attr {
	case "name":
		return a.activeSkill.Name, nil
	case "ttl":
		return a.activeSkill.Ttl, nil
	case "dir":
		return a.activeSkill.Dir, nil
	case "type":
		return a.activeSkill.Type, nil
	case "min":
		return a.activeSkill.Min, nil
	case "max":
		return a.activeSkill.Max, nil
	case "speed":
		return a.activeSkill.Speed, nil
	default:
		return "", errors.New("unknown attribute")
	}
}

func (a *Actor) Strike(ttl float64) {
	if a.activeSkill == nil || a.activeSkill.Type != "melee" {
		// select one skill of type melee
		for _, s := range a.skills {
			if s.Type == "melee" {
				a.SetSkill(s)
				break
			}
		}
	}
	if a.activeSkill.Type == "melee" {
		power := a.strength
		w := a.activeSkill.Hitbox.W()
		h := a.activeSkill.Hitbox.H()
		minx := a.rect.Center().X + a.activeSkill.Hitbox.Min.X
		miny := a.rect.Center().Y + a.activeSkill.Hitbox.Min.Y
		if a.dir < 0 {
			minx = a.rect.Center().X - a.activeSkill.Hitbox.Min.X - w
		}
		rect := pixel.R(minx, miny, minx+w, miny+h)

		activities.AddStrike(a, rect, power, a.vel, activities.WithTTL(ttl))
	}
}

func (a *Actor) UseSkill() {
	if a.activeSkill == nil {
		return
	}
	if a.activeSkill.Type == "spell" {
		a.w.AddSpell(a, a.target, a.activeSkill.Name, a.currObjs)
	}
	if a.activeSkill.Type == "ranged" {
		//t string, strength float64, owner common.Actorer
		projectiles.AddProjectile(a.activeSkill.Name, 1, a.dir, a)
	}
}

func (a *Actor) SetSkill(s *Skill) {
	a.activeSkill = s
}

func (a *Actor) SetTarget(t pixel.Vec) {
	a.target = t
}

func (a *Actor) SetEnemy(en *Actor) {
	a.enemy = en
}

func (a *Actor) GetEnemy() *Actor {
	return a.enemy
}

func (a *Actor) AddSound(event string) {
	if s, ok := a.sounds[event]; ok {
		// select random sound
		i := int(math.Round(common.GetRandFloat() * float64(len(s.List)-1)))
		sound.AddEffect(s.List[i], a.rect.Center())
	}
}

func (a *Actor) GetHp() int {
	return a.hp
}

func (a *Actor) GetPortrait() *pixel.Sprite {
	return a.portrait
}

func (a *Actor) IsGround() bool {
	return a.phys.IsGround()
}

func (a *Actor) Inform(e int) {
	for _, s := range a.esbrs {
		s.EventAction(e)
	}
}

func (a *Actor) AddEventListener(s common.EventSubscriber) {
	a.esbrs = append(a.esbrs, s)
}

func (a *Actor) Hit(vec pixel.Vec, power int) {
	if _, ok := a.states[state.HIT]; !ok { // cannot hit unhittable
		return
	}

	if a.state.GetId() == state.ROLL { // cannot hit in rolling state
		return
	}

	v := pixel.V(vec.X*(20*math.Max(a.walkspeed+float64(power)*10-a.mass, 0)), a.grav*(8+float64(power)*10-a.mass))
	a.phys.Apply(v)

	a.hp -= power
	if a.hp <= 0 {
		a.Kill()
		return
	}
	a.SetState(state.HIT)
	if a.onhit != nil {
		a.onhit(a.rect.Center(), v)
	}
}

func (a *Actor) Kill() {
	a.hp = 0
	if a.iswater {
		a.SetState(state.DEADSUNK)
	} else {
		a.SetState(state.DEAD)
	}

	// a.mass = 0
	// a.phys.SetMass(a.mass)
	a.Inform(events.GAMEVENT_DIE)
	a.onkill = nil
	a.onhit = nil
	a.oninteract = nil
}

// func (a *Actor) Phrasing() {
// 	talks.AddPhrase(a.rect.Min, a.phrasesClass)
// }

func (a *Actor) OnKill() {
	if len(a.onkill) == 0 { // == nil {
		return
	}
	// create random velocity vec
	rnd := float64(common.GetRandInt() - 5)
	y := 100.0 * rnd
	for _, okh := range a.onkill {
		okh(a.rect.Center(), pixel.V(30*float64(common.GetRandInt()-5), y))
	}
}

func (a *Actor) SetOnKill(okh OnKillHandler) {
	//	a.onkill = okh
	a.onkill = append(a.onkill, okh)
}

func (a *Actor) Interact() {
	activities.AddInteraction(a, a.rect, 1, pixel.ZV) // owner common.Actorer, rect pixel.Rect, power int, speed pixel.Vec)
}

func (a *Actor) SetOnInteract(oih OnInteractHandler) {
	a.oninteract = oih
}

func (a *Actor) OnInteract() {
	// now only one interaction posible: phrase
	if a.oninteract == nil {
		return
	}
	//	a.oninteract()
	a.oninteract(a)
}

func (a *Actor) GetId() int {
	return a.id
}
