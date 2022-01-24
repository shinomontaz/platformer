package ai

import "github.com/faiface/pixel"

type Subscriber interface {
	GetId() int
	Notify(e int, v pixel.Vec)
}
