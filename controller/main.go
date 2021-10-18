package controller

import (
	"github.com/faiface/pixel"
	"github.com/faiface/pixel/pixelgl"
)

type Controller struct {
	win  *pixelgl.Window
	vec  pixel.Vec
	sbrs map[int]Subscriber
}

func New(win *pixelgl.Window) *Controller {
	ctrl := &Controller{
		win: win,
	}

	return ctrl
}

func (pc *Controller) Subscribe(s Subscriber) {
	pc.sbrs[s.GetId()] = s
}

func (pc *Controller) Notify(e int) {
	for _, s := range pc.sbrs {
		s.Notify(e, &pc.vec)
	}
}

func (pc *Controller) Update() {
	pc.vec = pixel.ZV
	isMoved := false

	if pc.win.Pressed(pixelgl.KeyEscape) {
		pc.Notify(E_ESCAPE)
	}

	if pc.win.Pressed(pixelgl.KeyLeftControl) {
		pc.Notify(E_CTRL)
	}

	if pc.win.Pressed(pixelgl.KeyLeftShift) {
		pc.Notify(E_SHIFT)
	}

	if pc.win.Pressed(pixelgl.KeyLeft) {
		pc.vec.X--
		isMoved = true
	} else if pc.win.Pressed(pixelgl.KeyRight) {
		pc.vec.X++
		isMoved = true
	}

	if pc.win.Pressed(pixelgl.KeyUp) {
		pc.vec.Y++
		isMoved = true
	}

	if isMoved {
		pc.Notify(E_MOVE)
	}

}
