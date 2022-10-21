package main

import (
	"fmt"
	"math/rand"
	"time"

	"platformer/controller"
	"platformer/stages"

	"net/http"
	_ "net/http/pprof"

	"github.com/shinomontaz/pixel"
	"github.com/shinomontaz/pixel/pixelgl"
)

var (
	//	b *background.Pback
	//	b          *background.Back
	//	u          *ui.Ui
	//	w          *world.World
	win *pixelgl.Window
	//	hero       *actor.Actor
	ctrl       *controller.Controller
	title      string     = "platformer"
	currBounds pixel.Rect // current viewport

	//	initialCenter pixel.Vec
	//	lastPos       pixel.Vec
	isquit  bool
	isdebug bool

	currStage    stages.Stager
	stgs         map[int]stages.Stager
	loadingStage stages.Stager
)

func init() {
	rand.Seed(time.Now().UTC().UnixNano())

	// read start config

	// load video mode and sound volumes
	initRuntime()
}

func run() {
	cfg := pixelgl.WindowConfig{
		Title:  title,
		Bounds: currBounds,
		VSync:  true,
	}

	win, err := pixelgl.NewWindow(cfg)
	if err != nil {
		panic(err)
	}
	win.SetSmooth(true)
	ctrl = controller.New(win)

	initScreen(win)

	stgs = make(map[int]stages.Stager, 0)
	loadingStage = stages.NewLoading(inform, assetloader)
	stgs[stages.LOADING] = loadingStage
	stgs[stages.MENU] = stages.NewMenu(inform, assetloader, ctrl, currBounds)
	stgs[stages.GAME] = stages.NewGame(inform, assetloader, ctrl, currBounds)

	currStage = stgs[stages.LOADING]

	currStage.SetUp(stages.WithJob(stgs[stages.MENU].Init), stages.WithNext(stages.EVENT_DONE, stages.MENU))
	currStage.Init()
	currStage.Start()

	if startConfig.TestFlag {
		go func() {
			http.ListenAndServe("localhost:5000", nil)
		}()
	}

	last := time.Now()
	for !win.Closed() && !isquit {
		dt := time.Since(last).Seconds()
		last = time.Now()
		ctrl.Update() // - here we capture control signals, so actor physics receive input from controller

		currStage.Run(win, dt)

		win.Update()
	}
}

func main() {
	pixelgl.Run(run)
}

func inform(e int) {
	switch e {
	case stages.EVENT_DONE:
		fmt.Println("event done")
		next, ok := currStage.GetNext(stages.EVENT_DONE)
		if ok {
			setStage(next)
		}
	case stages.EVENT_ENTER:
		fmt.Println("event enter")
		next, ok := currStage.GetNext(stages.EVENT_ENTER)
		if ok {
			setStage(next)
		}
	case stages.EVENT_QUIT:
		fmt.Println("event quit")
		next, ok := currStage.GetNext(stages.EVENT_QUIT)
		if ok {
			setStage(next)
		} else {
			isquit = true
		}

	case stages.EVENT_INITSCREEN:
		initScreen(win)
	case stages.EVENT_NOTREADY:
		fmt.Println("event not ready")
		loadingStage.SetUp(stages.WithJob(currStage.Init), stages.WithNext(stages.EVENT_DONE, currStage.GetID()))
		setStage(loadingStage.GetID())
	}
}

func setStage(id int) {
	currStage = stgs[id]
	currStage.Start()
}
