package main

import (
	"fmt"
	"platformer/events"
	"time"

	"platformer/stages"

	"net/http"
	_ "net/http/pprof"

	"github.com/shinomontaz/pixel"
	"github.com/shinomontaz/pixel/pixelgl"
)

var (
	win        *pixelgl.Window
	title      string     = "platformer"
	currBounds pixel.Rect // current viewport

	isquit  bool
	isdebug bool

	currStage    stages.Stager
	stgs         map[int]stages.Stager
	loadingStage stages.Stager
)

func init() {
	// load video mode and sound volumes
	initRuntime()
}

func run() {
	cfg := pixelgl.WindowConfig{
		Title:  title,
		Bounds: currBounds,
		VSync:  true,
	}

	var err error

	win, err = pixelgl.NewWindow(cfg)
	if err != nil {
		panic(err)
	}
	win.SetSmooth(true)

	initScreen(win)
	initSound()

	stgs = make(map[int]stages.Stager, 0)
	loadingStage = stages.NewLoading(inform, assetloader)
	stgs[stages.LOADING] = loadingStage
	stgs[stages.MENU] = stages.NewMenu(inform, assetloader, win, currBounds) // main menu
	stgs[stages.GAME] = stages.NewGame(inform, assetloader, win, currBounds)

	currStage = stgs[stages.LOADING]

	currStage.SetUp(stages.WithJob(stgs[stages.MENU].Init), stages.WithNext(events.STAGEVENT_DONE, stages.MENU))
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

		currStage.Run(win, dt)

		win.Update()
	}
}

func main() {
	pixelgl.Run(run)
}

func inform(e int) {
	switch e {
	case events.STAGEVENT_DONE:
		fmt.Println("event done")
		next, ok := currStage.GetNext(events.STAGEVENT_DONE)
		if ok {
			setStage(next)
		}
	case events.STAGEVENT_NEXT:
		fmt.Println("event next")
		next, ok := currStage.GetNext(events.STAGEVENT_NEXT)
		if ok {
			setStage(next)
		}
	case events.STAGEVENT_QUIT:
		fmt.Println("event quit")
		next, ok := currStage.GetNext(events.STAGEVENT_QUIT)
		if ok {
			setStage(next)
		} else {
			isquit = true
		}
	case events.GAMEVENT_INITSCREEN:
		initScreen(win)
	case events.GAMEVENT_UPDATEVOLUME:
		initSound()
	case events.STAGEVENT_NOTREADY:
		loadingStage.SetUp(stages.WithJob(currStage.Init), stages.WithNext(events.STAGEVENT_DONE, currStage.GetID()))
		setStage(loadingStage.GetID())
	}
}

func setStage(id int) {
	currStage = stgs[id]
	currStage.Start()
}
