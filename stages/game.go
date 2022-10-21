package stages

import (
	"fmt"
	"image/color"
	"platformer/actor"
	"platformer/animation"
	"platformer/background"
	"platformer/common"
	"platformer/config"
	"platformer/controller"
	"platformer/events"
	"platformer/factories"
	"platformer/magic"
	"platformer/sound"
	"platformer/ui"
	"platformer/world"
	"time"

	"github.com/shinomontaz/pixel"
	"github.com/shinomontaz/pixel/pixelgl"
)

type Game struct {
	Common
	assetloader   *common.Loader
	currBounds    pixel.Rect
	u             *ui.Ui
	w             *world.World
	hero          *actor.Actor
	ctrl          *controller.Controller
	initialCenter pixel.Vec
	lastPos       pixel.Vec
}

var (
	camPos   = pixel.ZV
	frames   = 0
	second   = time.Tick(time.Second)
	rgba     = color.RGBA{123, 175, 213, 1}
	deltaVec = pixel.ZV
	title    = "platformer"
)

func NewGame(f Inform, l *common.Loader, ctrl *controller.Controller, currBounds pixel.Rect) *Game {
	return &Game{
		Common: Common{
			id:       GAME,
			done:     make(chan struct{}),
			inform:   f,
			eventMap: map[int]int{EVENT_QUIT: MENU},
		},
		assetloader: l,
		ctrl:        ctrl,
		currBounds:  currBounds,
	}
}

func (g *Game) Start() {
	if !g.isReady {
		return
	}
	sound.Init(g.assetloader)
	sound.PlayMusic("main")
	g.isActive = true
}

func (g *Game) Init() {
	animation.Init(g.assetloader)
	actor.Init(g.assetloader)
	ui.Init(g.assetloader)

	for _, anim := range config.AnimConfig {
		animation.Load(&anim)
	}
	for name, cfg := range config.Spells {
		magic.Load(name, &cfg)
	}

	w, err := world.New("ep1.tmx", g.currBounds, world.WithLoader(g.assetloader))
	if err != nil {
		panic(err)
	}
	g.w = w

	//	w.IsDebug = isdebug
	g.w.InitEnemies()

	magic.SetWorld(g.w)

	g.initialCenter = g.w.GetCenter()
	g.hero = factories.NewActor(config.Profiles["player"], g.w)
	g.hero.Move(g.initialCenter)
	g.ctrl.Subscribe(g.hero)
	g.w.AddHero(g.hero)
	g.u = ui.New(g.hero, g.currBounds)

	g.currBounds = g.currBounds.Moved(g.initialCenter.Sub(pixel.V(g.currBounds.W()/2, g.currBounds.H()/2)))

	g.lastPos = g.hero.GetPos()
	camPos = g.lastPos.Add(pixel.V(0, 150))

	b := background.New(g.lastPos, g.currBounds.Moved(pixel.V(0, 150)), g.assetloader, "gamebackground.png")

	g.w.SetBackground(b)

	g.ctrl.Subscribe(g)

	g.isReady = true
}

func (g *Game) Run(win *pixelgl.Window, dt float64) {
	if !g.isReady {
		g.Notify(EVENT_NOTREADY)
		return
	}

	win.Clear(rgba)

	pos := g.hero.GetPos()
	sound.Update(pos)
	if dt > 0 {
		deltaVec = g.lastPos.To(pos)
		camPos = pos.Add(pixel.V(0, 150))

		g.currBounds = g.currBounds.Moved(deltaVec)

		g.w.Update(g.currBounds, dt)
	}

	g.w.Draw(win, pos, camPos)
	g.u.Draw(win)

	g.lastPos = pos

	frames++
	select {
	case <-second:
		win.SetTitle(fmt.Sprintf("%s | FPS: %d", title, frames))
		frames = 0
	default:
	}
}

func (g *Game) Listen(e int, v pixel.Vec) {
	if !g.isActive {
		return
	}
	if e == events.ESCAPE {
		g.Notify(EVENT_QUIT)
	}
}
