package controller

import (
	"platformer/events"

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
		win:  win,
		sbrs: make(map[int]Subscriber),
	}

	return ctrl
}

func (pc *Controller) Subscribe(s Subscriber) {
	pc.sbrs[s.GetId()] = s
}

func (pc *Controller) Notify(e int) {
	for _, s := range pc.sbrs {
		s.Notify(e, pc.vec)
	}
}

func (pc *Controller) Update() {
	pc.vec = pixel.ZV
	isMoved := false

	if pc.win.JustPressed(pixelgl.KeyEscape) {
		pc.Notify(events.ESCAPE)
	}

	if pc.win.JustPressed(pixelgl.KeyLeftControl) {
		pc.Notify(events.CTRL)
	}

	if pc.win.JustPressed(pixelgl.KeyLeftShift) || pc.win.JustReleased(pixelgl.KeyLeftShift) {
		pc.Notify(events.SHIFT)
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

	if pc.win.JustReleased(pixelgl.KeyLeft) || pc.win.JustReleased(pixelgl.KeyRight) || pc.win.JustReleased(pixelgl.KeyUp) {
		isMoved = true
	}

	if isMoved {
		pc.Notify(events.MOVE)
	}

}
