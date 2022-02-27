package main

import "github.com/faiface/pixel"

var isfullscreen bool

// get window mode, sound volumes and menu/game mode
// currBounds
func initRuntime() {
	currBounds = pixel.R(0, 0, 600, 500)
	isfullscreen = false
}
