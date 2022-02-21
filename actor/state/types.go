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
	Listen(e int, v *pixel.Vec)
	GetSprite() *pixel.Sprite
	Busy() bool
}

type Actor interface {
	SetState(int)
	GetTransition(int) statemachine.Transition
	Strike()
	Cast()
	AddSound(event string)
	Inform(e int, v pixel.Vec)
}

type Animater interface {
	GetSprite(name string, idx int) (pixel.Picture, pixel.Rect)
	GetGroupSprite(group, name string, idx int) (pixel.Picture, pixel.Rect)
	GetGroupLen(name string) int
	GetLen(name string) int
}
