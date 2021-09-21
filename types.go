package main

import (
	"github.com/faiface/pixel"
	"github.com/faiface/pixel/imdraw"
	"github.com/faiface/pixel/pixelgl"
)

type IController interface {
	Update(win *pixelgl.Window)
}

type Objecter interface {
	Draw(imd *imdraw.IMDraw)
	Rect() *pixel.Rect
}
