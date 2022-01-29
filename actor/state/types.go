package state

import (
	"platformer/actor/statemachine"

	"github.com/faiface/pixel"
)

const (
	STAND = iota
	IDLE
	WALK
	RUN
	JUMP
	FALL
	HIT
	DEAD
	ATTACK
	MELEE
	MELEEMOVE
	RANGED
	CAST
)

type Stater interface {
	GetId() int
	Start()
	Update(dt float64)
	Notify(e int, v *pixel.Vec)
	GetSprite() *pixel.Sprite
}

type Actor interface {
	SetState(int)
	GetTransition(int) statemachine.Transition
}

type Animater interface {
	GetSprite(name string, idx int) (pixel.Picture, pixel.Rect)
	GetGroupSprite(group, name string, idx int) (pixel.Picture, pixel.Rect)
	GetGroupLen(name string) int
}
