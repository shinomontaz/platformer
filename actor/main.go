package actor

import (
	"math"
	"platformer/common"
	"platformer/events"
	"platformer/projectiles"
	"platformer/sound"
	"platformer/talks"

	"platformer/activities"

	"platformer/actor/state"
	"platformer/actor/statemachine"

	"github.com/shinomontaz/pixel"
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
	ALERTED
)

var counter int

type Actor struct {
	id   int
	phys common.Phys

	state  Stater
	states map[int]Stater

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

	sbrs []common.Subscriber

	skills       []*Skill
	activeSkill  *Skill
	currObjs     []common.Objecter
	phrasesClass string

	onkill     OnKillHandler
	onhit      OnKillHandler
	oninteract OnInteractHandler
}

var loader *common.Loader

func Init(l *common.Loader) {
	loader = l
	counter = 0
}

func New(w Worlder, anim common.Animater, rect pixel.Rect, opts ...Option) *Actor {
	counter++
	a := &Actor{
		id:      counter,
		anim:    anim,
		rect:    rect,
		dir:     1,
		animdir: 1,
		vel:     pixel.ZV,
		grav:    w.GetGravity(),
		w:       w,
		mass:    1,
		sounds:  make(map[string]soundeffect),
		sbrs:    make([]common.Subscriber, 0),
		skills:  make([]*Skill, 0),
	}

	for _, opt := range opts {
		opt(a)
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
	// basic states

	sStand := state.New(state.STAND, a, a.anim)
	sWalk := state.New(state.WALK, a, a.anim)
	sRun := state.New(state.RUN, a, a.anim)
	sIdle := state.New(state.IDLE, a, a.anim)
	sJump := state.New(state.JUMP, a, a.anim)
	sFall := state.New(state.FALL, a, a.anim)
	sDead := state.New(state.DEAD, a, a.anim)
	sHit := state.New(state.HIT, a, a.anim)
	sInteract := state.New(state.INTERACT, a, a.anim)

	a.states = map[int]Stater{
		state.STAND:    sStand,
		state.IDLE:     sIdle,
		state.WALK:     sWalk,
		state.RUN:      sRun,
		state.FALL:     sFall,
		state.JUMP:     sJump,
		state.HIT:      sHit,
		state.DEAD:     sDead,
		state.INTERACT: sInteract,
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
	if e == events.WALK {
		v.X *= a.walkspeed / 20
	} else {
		v.X *= a.runspeed / 20
	}

	groundrate, _, _ := a.phys.StepPrediction(v)
	return groundrate
}

func (a *Actor) Listen(e int, v pixel.Vec) {
	a.appliedForce = pixel.ZV

	if a.state.Busy() {
		return
	}

	if v.X != 0 {
		if v.X > 0 {
			a.dir = 1
		} else {
			a.dir = -1
		}
		if a.isShift {
			if math.Abs(a.vel.X) < a.runspeed {
				a.appliedForce.X = v.X * a.runspeed
			}
		} else {
			if math.Abs(a.vel.X) < a.walkspeed*19/20 {
				a.appliedForce.X = v.X * a.walkspeed
			}
		}
	}

	if v.Y > 0 {
		multiplier := 0.5
		if math.Abs(a.vel.X) > a.walkspeed {
			multiplier = 1
		}
		a.appliedForce.Y = v.Y * a.grav * (a.jumpforce + multiplier*float64(a.strength) - a.mass)
	}

	if e == events.SHIFT {
		a.isShift = !a.isShift
	}

	a.state.Listen(e, &a.vel)
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

func (a *Actor) Update(dt float64, objs []common.Objecter) {
	a.currObjs = objs
	if a.appliedForce.X == 0 && math.Abs(a.vel.X) > a.walkspeed*0.1 && a.vel.Y == 0 { // active slowing in case of no active force applied
		a.appliedForce.X = -a.walkspeed
		if a.vel.X < 0 {
			a.appliedForce.X *= -1
		}

	}
	a.phys.Apply(a.appliedForce)
	a.phys.Update(dt, objs)
	newspeed := a.phys.GetVel()
	var event int
	if math.Abs(newspeed.X) <= a.runspeed && math.Abs(newspeed.X) > a.walkspeed {
		event = events.RUN
	} else if (math.Abs(a.vel.X) >= a.walkspeed && math.Abs(newspeed.X) <= a.walkspeed) || (a.vel.X == 0 && math.Abs(newspeed.X) > 0 && math.Abs(newspeed.X) <= a.walkspeed) {
		event = events.WALK
	}
	a.state.Listen(event, newspeed)
	a.vel = *newspeed

	a.rect = a.phys.GetRect()
	a.state.Update(dt)
	a.appliedForce = pixel.ZV
}

func (a *Actor) UpdateSpecial(dt float64, objs []common.Objecter) {
	force := pixel.ZV

	for _, o := range objs {
		isIntercects := a.rect.Intersects(o.R)
		if !isIntercects { // no collision
			continue
		}

		part := a.rect.Intersect(o.R)

		if o.Type == common.WATER {
			ratio := part.H() / a.rect.H()
			if ratio == 1 {
				if a.hp > 0 {
					a.Kill()
				}
				ratio += 0.5
			}

			force.Y += a.grav * ratio // archimedus
			a.phys.SetSpeed(a.vel.Scaled(0.9))
		}
	}
	a.phys.Apply(force)
}

func (a *Actor) SetState(id int) {
	a.state = a.states[id]
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
	drawrect := a.rect.ResizedMin(pixel.Vec{a.rect.W() * 1.25, a.rect.H() * 1.25})
	a.sprite.Draw(t, pixel.IM.
		ScaledXY(pixel.ZV, pixel.V(
			drawrect.W()/a.sprite.Frame().W(),
			drawrect.H()/a.sprite.Frame().H(),
		)).
		ScaledXY(pixel.ZV, pixel.V(a.animdir*a.dir, 1)).
		Moved(drawrect.Center()),
	)
	//	a.phys.Draw(t)
}

func (a *Actor) GetSkills() []*Skill {
	return a.skills
}

func (a *Actor) GetSkillName() string {
	if a.activeSkill != nil {
		return a.activeSkill.Name
	}
	return ""
}

func (a *Actor) Strike() {
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

		activities.AddStrike(a, rect, power, pixel.ZV)
	}
}

func (a *Actor) Cast() {
	if a.activeSkill.Type == "spell" {
		//		activities.AddSpell(a, a.target, a.activeSkill.Name, a.currObjs)
		a.w.AddSpell(a, a.target, a.activeSkill.Name, a.currObjs)
	}
	if a.activeSkill.Type == "ranged" {
		//t string, pos, f pixel.Vec, strength float64, owner common.Actorer
		projectiles.AddProjectile(a.activeSkill.Name, a.GetRect().Center(), pixel.V(10000*a.dir, 0), 1, a.dir, a)
	}

}

func (a *Actor) SetSkill(s *Skill) {
	a.activeSkill = s
}

func (a *Actor) SetTarget(t pixel.Vec) {
	a.target = t
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

func (a *Actor) Inform(e int, v pixel.Vec) {
	for _, s := range a.sbrs {
		s.Listen(e, v)
	}
}

func (a *Actor) AddListener(s common.Subscriber) {
	a.sbrs = append(a.sbrs, s)
}

func (a *Actor) Hit(vec pixel.Vec, power int) {
	if _, ok := a.states[state.HIT]; !ok { // cannot hit unhittable
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
	a.Inform(events.ALERT, pixel.Vec{-vec.X, vec.Y})
	if a.onhit != nil {
		a.onhit(a.rect.Center(), v)
	}
}

func (a *Actor) Kill() {
	a.hp = 0
	a.SetState(state.DEAD)
	// a.mass = 0
	// a.phys.SetMass(a.mass)
	a.Inform(events.GAMEVENT_DIE, pixel.ZV)
	a.onkill = nil
	a.onhit = nil
	a.oninteract = nil
}

func (a *Actor) Phrasing() {
	// talks.AddPhrase(a.rect.Min, a.phrasesClass)
	talks.AddPhrase(a.rect.Min, a.phrasesClass)
}

func (a *Actor) OnKill() {
	if a.onkill == nil {
		return
	}
	// create random velocity vec
	rnd := float64(common.GetRandInt() - 5)
	y := 100.0 * rnd
	a.onkill(a.rect.Center(), pixel.V(30*float64(common.GetRandInt()-5), y))
}

func (a *Actor) SetOnKill(okh OnKillHandler) {
	a.onkill = okh
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
	a.oninteract()
}

func (a *Actor) GetId() int {
	return a.id
}
