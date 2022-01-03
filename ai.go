package main

import (
	"github.com/faiface/pixel"
)

type Ai struct {
	vec    pixel.Vec
	cmd    int // command action
	ground bool
	//	attack *Attack
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
