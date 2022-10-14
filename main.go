package main

import (
	"math/rand"
	"time"

	"platformer/actor"
	"platformer/background"
	"platformer/ui"
	"platformer/world"

	"platformer/controller"

	"net/http"
	_ "net/http/pprof"

	"github.com/shinomontaz/pixel"
	"github.com/shinomontaz/pixel/pixelgl"
)

func init() {
	rand.Seed(time.Now().UTC().UnixNano())

	// read start config

	// load video mode and sound volumes
	initRuntime()
}

var (
	//	b *background.Pback
	b          *background.Back
	u          *ui.Ui
	w          *world.World
	win        *pixelgl.Window
	hero       *actor.Actor
	ctrl       *controller.Controller
	title      string     = "platformer"
	currBounds pixel.Rect // current viewport

	initialCenter pixel.Vec
	lastPos       pixel.Vec
	ismenu        bool
	isquit        bool
	isdebug       bool
)

func gameLoop(win *pixelgl.Window) {
	last := time.Now()

	for !win.Closed() && !isquit {
		dt := time.Since(last).Seconds()
		last = time.Now()
		ctrl.Update() // - here we capture control signals, so actor physics receive input from controller

		if ismenu {
			menuFunc(win, dt)
		} else {
			gameFunc(win, dt)
		}
		win.Update()
	}
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
	initMenu(win)

	ismenu = true

	go func() {
		http.ListenAndServe("localhost:5000", nil)
	}()

	gameLoop(win)
}

func main() {
	pixelgl.Run(run)
}
