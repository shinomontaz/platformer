package common

import "github.com/shinomontaz/pixel"

type Subscriber interface {
	//	GetId() int
	Listen(e int, v pixel.Vec)
}

const (
	GROUND = iota
	BARRIER
	WATER
)
