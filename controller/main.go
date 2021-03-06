package controller

import (
	"platformer/common"
	"platformer/events"

	"github.com/shinomontaz/pixel"
	"github.com/shinomontaz/pixel/pixelgl"
)

type Controller struct {
	win  *pixelgl.Window
	vec  pixel.Vec
	sbrs []common.Subscriber
}

func New(win *pixelgl.Window) *Controller {
	ctrl := &Controller{
		win:  win,
		sbrs: make([]common.Subscriber, 0),
	}

	return ctrl
}

func (pc *Controller) Subscribe(s common.Subscriber) {
	pc.sbrs = append(pc.sbrs, s)
}

func (pc *Controller) Notify(e int) {
	for _, s := range pc.sbrs {
		s.Listen(e, pc.vec)
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

	if pc.win.JustPressed(pixelgl.KeyUp) {
		pc.vec.Y++
		isMoved = true
	}
	if pc.win.JustPressed(pixelgl.KeyDown) {
		pc.vec.Y--
		isMoved = true
	}

	if pc.win.JustPressed(pixelgl.KeyEnter) {
		pc.Notify(events.ENTER)
	}

	if pc.win.JustReleased(pixelgl.KeyLeft) || pc.win.JustReleased(pixelgl.KeyRight) || pc.win.JustReleased(pixelgl.KeyUp) {
		isMoved = true
	}

	if isMoved {
		pc.Notify(events.MOVE)
	}

}
