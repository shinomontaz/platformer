package gamestate

import (
	"fmt"
	"platformer/controller"
	"platformer/events"

	"github.com/shinomontaz/pixel/pixelgl"
)

type Victory struct {
	Common
	win  *pixelgl.Window
	ctrl *controller.Controller
}

func NewVictory(game Gamer, win *pixelgl.Window) *Victory {
	v := &Victory{
		Common: Common{
			game: game,
			id:   VICTORY,
		},
		win:  win,
		ctrl: controller.New(win, true),
	}

	v.ctrl.AddKeyListener(v) // to listen ESCAPE  keyborad event

	return v
}

func (v *Victory) Update(dt float64) {
	if dt > 0 {
		v.ctrl.Update()
	}
}

func (v *Victory) Draw(win *pixelgl.Window) {
	//	camPos := v.lastPos.Add(pixel.V(0, 150))
}

func (v *Victory) GetId() int {
	return v.id
}

func (v *Victory) Start() {
	fmt.Println("victory screen")
}

func (v *Victory) KeyAction(key pixelgl.Button) {
	switch key {
	case pixelgl.KeyEscape: // from controller
		//		v.game.SetState(MENU)
		v.game.Notify(events.STAGEVENT_QUIT)
	}
}
