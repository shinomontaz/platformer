package actor

import (
	"math"
	"math/rand"
	"platformer/common"
	"platformer/events"
	"platformer/sound"

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
)

type Actor struct {
	phys Phys

	state  Stater
	states map[int]Stater

	rect pixel.Rect

	animdir   float64
	anim      Animater
	sprite    *pixel.Sprite
	dir       float64
	vec       pixel.Vec // delta speed
	vel       pixel.Vec
	runspeed  float64
	walkspeed float64
	grav      float64
	jumpforce float64
	isShift   bool

	sm *statemachine.Machine
	w  Worlder

	hp          int
	strength    int
	portrait    *pixel.Sprite
	sounds      map[string]soundeffect
	attackrange float64

	target pixel.Vec

	sbrs []common.Subscriber

	skills      []*Skill
	activeSkill *Skill
}

func New(w Worlder, anim Animater, rect pixel.Rect, opts ...Option) *Actor {
	a := &Actor{
		anim:    anim,
		rect:    rect,
		dir:     1,
		animdir: 1,
		vel:     pixel.ZV,
		grav:    w.GetGravity(),
		w:       w,
		sounds:  make(map[string]soundeffect),
		sbrs:    make([]common.Subscriber, 0),
		skills:  make([]*Skill, 0),
	}

	for _, opt := range opts {
		opt(a)
	}

	p := NewPhys(rect, a.runspeed, a.grav) // TODO does we really need phys to know run and walk speeds?
	p.SetQt(w.GetQt())
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

	a.states = map[int]Stater{
		state.STAND: sStand,
		state.IDLE:  sIdle,
		state.WALK:  sWalk,
		state.RUN:   sRun,
		state.FALL:  sFall,
		state.JUMP:  sJump,
		state.HIT:   sHit,
		state.DEAD:  sDead,
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

func (a *Actor) Listen(e int, v pixel.Vec) {
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
				v.X *= a.runspeed / 20
			} else {
				v.X = 0
			}
		} else {
			if math.Abs(a.vel.X) < a.walkspeed*19/20 {
				v.X *= a.walkspeed / 20
			} else {
				v.X = 0
			}
		}
	}

	if v.Y > 0 {
		v.Y *= a.grav * a.jumpforce
	}

	if e == events.SHIFT {
		a.isShift = !a.isShift
	}

	a.state.Listen(e, &a.vel)

	a.vec = v
}

func (a *Actor) Move(v pixel.Vec) {
	//	a.rect = pixel.R(pos.X, pos.Y, pos.X+a.rect.W(), pos.Y+a.rect.Y())
	a.rect = a.rect.Moved(v)
	a.phys.Move(v)
}

func (a *Actor) GetPos() pixel.Vec {
	return a.rect.Center()
}

func (a *Actor) GetRect() pixel.Rect {
	return a.rect
}

func (a *Actor) Update(dt float64) {
	a.phys.Update(dt, &a.vec)
	a.vec = pixel.ZV
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
}

func (a *Actor) SetState(id int) {
	a.state = a.states[id]
	a.state.Start()
}

func (a *Actor) GetDir() int {
	return int(a.dir)
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

		a.w.AddStrike(a, rect, power, pixel.ZV)
	}
}

func (a *Actor) Cast() {
	// TODO: get melee skill and cast spell by it
	//	activeSkill
	if a.activeSkill.Type == "spell" {
		a.w.AddSpell(a, a.target, a.activeSkill.Name)
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
		i := rand.Intn(len(s.List))
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
	return a.phys.ground
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
	vec.X *= a.walkspeed * float64(a.strength+1)
	vec.Y = 100
	a.vec = vec
	a.hp -= power
	if a.hp <= 0 {
		a.SetState(state.DEAD)
		a.Inform(events.DIE, pixel.ZV)
		return
	}
	a.SetState(state.HIT)
	a.Inform(events.ALERT, pixel.Vec{-vec.X, vec.Y})
}
