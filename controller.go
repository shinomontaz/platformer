package main

import (
	"github.com/faiface/pixel"
	"github.com/faiface/pixel/pixelgl"
)

type PlayerController struct {
	vec    pixel.Vec
	cmd    int // command action
	ground bool
}

func (pc *PlayerController) SetGround(g bool) {
	pc.ground = g
}

func (pc *PlayerController) Update(win *pixelgl.Window) {
	pc.cmd = NOACTION
	pc.vec = pixel.ZV

	if win.Pressed(pixelgl.KeyLeftShift) {
		if win.Pressed(pixelgl.KeyLeft) {
			pc.vec.X--
		} else if win.Pressed(pixelgl.KeyRight) {
			pc.vec.X++
		}
	} else if !win.Pressed(pixelgl.KeyLeftShift) { // running
		if win.Pressed(pixelgl.KeyLeft) {
			pc.vec.X -= 2.0
		} else if win.Pressed(pixelgl.KeyRight) {
			pc.vec.X += 2.0
		}
	}

	if !pc.ground {
		return
	}

	if win.Pressed(pixelgl.KeyLeftControl) {
		// attacking!
		pc.cmd = STRIKE
		return
	}

	if win.Pressed(pixelgl.KeyUp) {
		pc.vec.Y++
	}
}
