package common

import "github.com/faiface/pixel"

type Subscriber interface {
	GetId() int
	Listen(e int, v pixel.Vec)
}
