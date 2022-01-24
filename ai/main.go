package ai

import (
	"github.com/faiface/pixel"
)

var list []*Ai

type Ai struct {
	vec  pixel.Vec
	sbrs map[int]Subscriber
}

func New() *Ai {
	a := &Ai{
		sbrs: make(map[int]Subscriber),
	}
	list = append(list, a)
	return a
}

func Update() {
	for _, a := range list {
		a.Update()
	}
}

func (a *Ai) Update() {

}

func (a *Ai) Subscribe(s Subscriber) {
	a.sbrs[s.GetId()] = s
}

func (a *Ai) Notify(e int) {
	for _, s := range a.sbrs {
		s.Notify(e, a.vec)
	}
}

// func (ai *Ai) SetGround(g bool) {
// 	ai.ground = g
// }

// func (ai *Ai) Update(target pixel.Vec, objs []Objecter) {
// 	ai.cmd = NOACTION
// 	ai.vec = pixel.ZV

// 	pos := ai.pers.rect.Center()
// 	move := 0.0
// 	if math.Abs(pos.X-target.X) > 5 {
// 		if math.Signbit(pos.X - target.X) {
// 			move = 1.0
// 		} else {
// 			move = -1.0
// 		}
// 	}
// 	ai.vec.X += move

// 	if ai.pers.rect.Contains(target) {
// 		ai.cmd = STRIKE
// 	}
// }
