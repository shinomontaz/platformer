package gamestate

import (
	"fmt"
	"platformer/actor"
	"platformer/controller"
	"platformer/events"
	"platformer/sound"
	"platformer/ui"
	"platformer/world"

	"github.com/shinomontaz/pixel/pixelgl"

	"github.com/shinomontaz/pixel"
)

type Normal struct {
	Common
	shader string
	ctrl   *controller.Controller
	win    *pixelgl.Window
}

func NewNormal(game Gamer, currBounds pixel.Rect, u *ui.Ui, w *world.World, hero *actor.Actor, win *pixelgl.Window) *Normal {
	n := &Normal{
		Common: Common{
			game:       game,
			id:         NORMAL,
			currBounds: currBounds,
			u:          u,
			w:          w,
			hero:       hero,
			lastPos:    pixel.ZV,
		},
		win:  win,
		ctrl: controller.New(win),
	}

	n.ctrl.AddListener(n.hero) // make hero listen keyboard input
	n.ctrl.AddListener(n)      // to listen ESCAPE  keyborad event
	n.hero.AddListener(n)      // to listen DIE event

	return n
}

func (n *Normal) Update(dt float64) {
	pos := n.hero.GetPos()
	sound.Update(pos)
	if dt > 0 {
		n.ctrl.Update()

		deltaVec := n.lastPos.To(pos)
		n.currBounds = n.currBounds.Moved(deltaVec)

		n.w.Update(n.currBounds, dt)
	}

	n.lastPos = pos
}

func (n *Normal) Draw(win *pixelgl.Window) {
	camPos := n.lastPos.Add(pixel.V(0, 150))

	n.w.Draw(win, n.lastPos, camPos, win.Bounds().Center())
	n.u.Draw(win)
}

func (n *Normal) GetId() int {
	return n.id
}

func (n *Normal) Start() {
	n.lastPos = n.hero.GetPos()
	//	n.win.Canvas().ResetFragmentShader()
	fmt.Println("applying base shader")
}

func (n *Normal) Listen(e int, v pixel.Vec) {
	switch e {
	case events.ESCAPE: // from controller
		n.game.SetState(MENU)
	case events.DIE: // from hero
		fmt.Println("handle event DIE")
		n.game.SetState(DEAD)
	}
}
