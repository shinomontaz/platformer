package controller

import (
	"fmt"
	"platformer/bindings"
	"platformer/common"

	"github.com/shinomontaz/pixel/pixelgl"
)

type Controller struct {
	win       *pixelgl.Window
	sbrs      []common.KeySubscriber
	jr        bool
	listenall bool
	currBind  map[int]pixelgl.Button
}

var typesBinds = map[int]string{
	bindings.CTRL:  "pressed",
	bindings.ENTER: "pressed",
	bindings.SHIFT: "default",
	bindings.LEFT:  "default",
	bindings.RIGHT: "default",
	bindings.UP:    "default",
	bindings.DOWN:  "default",
}

func New(win *pixelgl.Window, justReleased bool) *Controller {
	ctrl := &Controller{
		win:      win,
		sbrs:     make([]common.KeySubscriber, 0),
		jr:       justReleased,
		currBind: make(map[int]pixelgl.Button),
	}

	// type of event ( pressed, justpressed, just released) => pixelgl.Button
	// event => pixelgl
	df := bindings.Default
	for key, val := range df.List() {
		ctrl.currBind[key] = val
	}

	return ctrl
}

func (pc *Controller) AddListener(s common.KeySubscriber) {
	pc.sbrs = append(pc.sbrs, s)
}

func (pc *Controller) SetListenAll(listenall bool) {
	pc.listenall = listenall
}

func (pc *Controller) Notify(b pixelgl.Button) {
	for _, s := range pc.sbrs {
		s.KeyEvent(b)
	}
}

func (pc *Controller) Update() {
	if pc.listenall { // notify for first key pressed
		pc.Notify(pc.KeyAll())
		return
	}

	binded := pc.KeysBinded()
	for _, b := range binded {
		pc.Notify(b)
	}
}

func (pc *Controller) KeyAll() pixelgl.Button {
	for _, b := range AllButtons {
		if pc.win.Pressed(b) {
			return b
		}
	}
	return pixelgl.KeyUnknown
}

func (pc *Controller) KeysBinded() []pixelgl.Button {
	res := make([]pixelgl.Button, 0)
	if pc.win.JustPressed(pixelgl.KeyEscape) {
		res = append(res, pixelgl.KeyEscape)
	}

	for ev, button := range pc.currBind {
		if pc.jr {
			if pc.win.JustReleased(button) {
				res = append(res, button)
			}
		} else {
			switch typesBinds[ev] {
			case "pressed":
				if pc.win.JustPressed(button) {
					res = append(res, button)
				}
			case "released":
				if pc.win.JustReleased(button) {
					res = append(res, button)
				}
			default:
				if pc.win.Pressed(button) {
					res = append(res, button)
				}
			}
		}

	}
	return res
}

func (pc *Controller) DetectPressedButton() pixelgl.Button {
	fmt.Println("DetectPressedButton")
	for _, b := range AllButtons {
		if pc.win.Pressed(b) {
			return b
		}
	}
	return pixelgl.KeyUnknown
}
