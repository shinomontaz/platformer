package main

import "github.com/faiface/pixel/pixelgl"

type IController interface {
	Update(win *pixelgl.Window)
}
