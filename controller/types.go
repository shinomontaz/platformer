package controller

import "github.com/faiface/pixel"

const (
	E_MOVE = iota
	E_CTRL
	E_SHIFT
	E_ESCAPE
)

type Subscriber interface {
	GetId() int
	Notify(e int, v *pixel.Vec)
}
