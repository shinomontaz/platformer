package gamestate

import (
	"fmt"
	"platformer/actor"
	"platformer/ai"
	"platformer/controller"
	"platformer/events"
	"platformer/sound"
	"platformer/talks"
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
		ctrl: controller.New(win, false),
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
		//		creatures.Update(dt, w.GetVisiblePhys())
		talks.Update(dt)
		ai.Update(dt)

		//		n.w.SetCreatures(creatures.List()) // we need shaders on creatures ( npc and enemies ), creatures will be drawn on world draw
	}

	n.lastPos = pos
}

func (n *Normal) Draw(win *pixelgl.Window) {
	camPos := n.lastPos.Add(pixel.V(0, 150))
	cntr := win.Bounds().Center()

	n.w.Draw(win, n.lastPos, camPos, cntr)
	talks.Draw(win, camPos, cntr) // we do not need world shaders upon alerts and replys
	//	objects.Draw(win, camPos, cntr)
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
	case events.GAMEVENT_DIE: // from hero
		fmt.Println("handle event DIE")
		n.game.SetState(DEAD)
	case events.INTERACT: // from hero
		fmt.Println("handle event INTERACT")
		// world -> check interactions, get interaction action
		// if dialog =>
		// hero => set active dialog
		// n.game.SetState(DIALOG)
		// if reply => world add reply
		// if action => do action on world
	}
}
