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
	"platformer/creatures"
	"platformer/events"
	"platformer/factories"
	"platformer/loot"
	"platformer/magic"
	"platformer/sound"
	"platformer/talks"
	"platformer/ui"
	"platformer/world"

	"platformer/stages/gamestate"
	"time"

	"github.com/shinomontaz/pixel"
	"github.com/shinomontaz/pixel/pixelgl"
)

type Game struct {
	Common
	assetloader   *common.Loader
	initialBounds pixel.Rect
	u             *ui.Ui
	w             *world.World
	hero          *actor.Actor
	ctrl          *controller.Controller
	initialCenter pixel.Vec
	lastPos       pixel.Vec

	state  Gamestater
	states map[int]Gamestater
	win    *pixelgl.Window
}

var (
	frames = 0
	second = time.Tick(time.Second)
	rgba   = color.RGBA{123, 175, 213, 1}
	title  = "platformer"
)

func NewGame(f Inform, l *common.Loader, win *pixelgl.Window, currBounds pixel.Rect) *Game {
	return &Game{
		Common: Common{
			id:       GAME,
			done:     make(chan struct{}),
			inform:   f,
			eventMap: map[int]int{events.STAGEVENT_QUIT: MENU},
		},
		assetloader:   l,
		ctrl:          controller.New(win),
		initialBounds: currBounds,
		win:           win,
	}
}

func (g *Game) Start() {
	if !g.isReady {
		return
	}
	sound.PlayMusic("main")
	g.isActive = true
}

func (g *Game) Init() {
	currBounds := g.initialBounds
	animation.Init(g.assetloader)
	actor.Init(g.assetloader) // to load portrait only
	ui.Init(g.assetloader)
	sound.Init(g.assetloader)

	for _, anim := range config.AnimConfig {
		animation.Load(&anim)
	}
	for name, cfg := range config.Spells {
		magic.Load(name, &cfg)
	}

	w, err := world.New("ep2.tmx", currBounds, world.WithLoader(g.assetloader))
	if err != nil {
		panic(err)
	}
	g.w = w
	g.hero = factories.NewActor(config.Profiles["player"], g.w)

	loot.Init(g.w, config.Loots)
	talks.Init(g.assetloader)
	creatures.Init()
	list := g.w.GetMetas()
	for _, o := range list {
		if o.Class == "enemy" {
			enemy := factories.NewActor(config.Profiles[o.Name], g.w)
			enemy.Move(pixel.V(o.X, o.Y))
			if o.Properties.GetString("reward") != "" {
				// make OnKill handler
				enemy.SetOnKill(loot.AddCoin)
			}
			factories.NewAi(config.Profiles[o.Name].Type, enemy, w)
			creatures.AddEnemy(enemy)
		}
		if o.Class == "npc" {
			npc := factories.NewActor(config.Profiles[o.Name], g.w)
			npc.Move(pixel.V(o.X, o.Y))
			factories.NewAi(config.Profiles[o.Name].Type, npc, w)
			creatures.AddNpc(npc)
		}
		if o.Class == "coin" {
			loot.AddCoin(pixel.V(o.X, o.Y), pixel.ZV)
		}
	}
	creatures.SetHero(g.hero)

	//	activities.Init(creatures.List()) // new

	//	w.IsDebug = isdebug

	magic.SetWorld(g.w)

	g.initialCenter = g.w.GetCenter()
	g.hero.Move(g.initialCenter)

	g.u = ui.New(g.hero, currBounds)

	currBounds = currBounds.Moved(g.initialCenter.Sub(pixel.V(currBounds.W()/2, currBounds.H()/2)))

	g.lastPos = g.hero.GetPos()

	b := background.New(g.lastPos, currBounds.Moved(pixel.V(0, 150)), g.assetloader, "gamebackground.png")

	g.w.SetBackground(b)

	g.initStates(currBounds)

	g.isReady = true
}

func (g *Game) initStates(currBounds pixel.Rect) {
	sNormal := gamestate.NewNormal(g, currBounds, g.u, g.w, g.hero, g.win)
	sDead := gamestate.NewDead(g, g.w, g.hero, g.win)
	sMenu := gamestate.NewMenu(g, g.w, g.hero, g.win)
	//	sDialog := gamestate.NewDialog(g, g.w, g.hero, g.win)

	g.states = map[int]Gamestater{
		gamestate.NORMAL: sNormal,
		gamestate.DEAD:   sDead,
		gamestate.MENU:   sMenu,
		//		gamestate.DIALOG: sDialog,
	}

	g.SetState(gamestate.NORMAL)
}

func (g *Game) SetState(id int) {
	if s, ok := g.states[id]; ok {
		g.state = s
		g.state.Start()
	}
}

func (g *Game) Run(win *pixelgl.Window, dt float64) {
	if !g.isReady {
		g.Notify(events.STAGEVENT_NOTREADY)
		return
	}

	win.Clear(rgba)

	g.state.Update(dt)
	g.state.Draw(win)

	frames++
	select {
	case <-second:
		win.SetTitle(fmt.Sprintf("%s | FPS: %d", title, frames))
		frames = 0
	default:
	}
}

func (g *Game) Notify(e int) {
	if e == events.STAGEVENT_QUIT {
		g.isReady = false
	}
	g.inform(e)
}

func (g *Game) Listen(e int, v pixel.Vec) {
	if !g.isActive {
		return
	}
	if e == events.ESCAPE {
		g.Notify(events.STAGEVENT_QUIT)
	}
}
