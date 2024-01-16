package state

import (
	"platformer/actor/statemachine"

	"github.com/shinomontaz/pixel"
)

const (
	STAND = iota
	SNEER
	IDLE
	WALK
	RUN
	JUMP
	FALL
	HIT
	DEAD
	DEADSUNK
	ATTACK
	MELEE
	MELEEMOVE
	RANGED
	CAST
	INTERACT
	RESURRECT
	FISHING
	SWIM
	ROLL
)

type Stater interface {
	GetId() int
	Start()
	Update(dt float64)
	Listen(e int, v *pixel.Vec)
	GetSprite() *pixel.Sprite
	Busy() bool
	SetWater(bool)
	SetWaterResistant(bool)
}

type Actor interface {
	SetState(int)
	GetTransition(int) statemachine.Transition
	Strike(ttl float64)
	UseSkill()
	Interact()
	AddSound(event string)
	Inform(e int)
	GetSkillAttr(attr string) (interface{}, error)
	OnKill()
	GetDir() int
	SetVel(v pixel.Vec)
}

type Animater interface {
	GetSprite(name string, idx int) (pixel.Picture, pixel.Rect)
	GetGroupSprite(group, name string, idx int) (pixel.Picture, pixel.Rect)
	GetGroupLen(name string) int
	GetLen(name string) int
}
