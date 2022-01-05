package main

import (
	"fmt"
	"image/color"
	"math"
	"math/rand"
	"time"

	"platformer/actor"
	"platformer/animation"
	"platformer/config"
	"platformer/world"

	"platformer/controller"

	"github.com/faiface/pixel"
	"github.com/faiface/pixel/pixelgl"
)

func init() {
	// create world
	// create physics
	// create hero

	rand.Seed(time.Now().UTC().UnixNano())
}

var (
	w             *world.World
	hero          *actor.Actor
	ctrl          *controller.Controller
	title         string = "platformer"
	initBounds    pixel.Rect
	initialCenter pixel.Vec
)

func gameLoop(win *pixelgl.Window) {

	var (
		camPos    = pixel.ZV
		frames    = 0
		second    = time.Tick(time.Second)
		frametime = time.Tick(120 * time.Millisecond)
	)

	last := time.Now()
	rgba := color.RGBA{205, 231, 244, 1}

	for !win.Closed() {
		dt := time.Since(last).Seconds()
		last = time.Now()
		win.Clear(rgba)

		pos := hero.GetPos()

		camPos = pixel.Lerp(camPos, initialCenter.Sub(pos), 1-math.Pow(1.0/128, dt))
		cam := pixel.IM.Moved(camPos)

		win.SetMatrix(cam)
		currBounds := initBounds.Moved(initialCenter.Sub(pos).Scaled(-1))

		ctrl.Update() // - here we capture control signals, so actor physics receive input from controller
		hero.Update(dt)
		w.Update(currBounds)

		w.Draw(win)
		hero.Draw(win)

		win.Update()

		frames++
		select {
		case <-frametime:
			//			hero.Tick()
		case <-second:
			win.SetTitle(fmt.Sprintf("%s | FPS: %d", title, frames))
			frames = 0
		default:
		}
	}
}

func run() {
	initBounds = pixel.R(0, 0, config.WorldConfig.Width, config.WorldConfig.Heigth)
	cfg := pixelgl.WindowConfig{
		Title:  title,
		Bounds: initBounds,
		VSync:  true,
	}

	win, err := pixelgl.NewWindow(cfg)
	if err != nil {
		panic(err)
	}

	win.SetSmooth(true)

	w = world.New(pixel.R(-2000000.0, -1000.0, 20000000.0, 1000.0))
	for _, p := range config.WorldConfig.Platforms {
		w.Add(world.NewPlatform(pixel.R(p[0], p[1], p[2], p[3]).Moved(win.Bounds().Center())))
	}

	w.SetGravity(config.WorldConfig.Gravity)

	ctrl = controller.New(win)

	mainRect := pixel.R(0, 0, config.PlayerConfig.Width, config.PlayerConfig.Height)
	initialCenter = win.Bounds().Center()

	playerAnims := animation.New(mainRect)
	for _, anim := range config.PlayerConfig.Anims {
		playerAnims.SetAnim(anim.Name, anim.File, anim.Frames)
	}

	hero = actor.New(w, playerAnims, mainRect, config.PlayerConfig.Run, config.PlayerConfig.Walk)
	hero.Move(initialCenter)
	ctrl.Subscribe(hero)

	gameLoop(win)
}

func main() {
	pixelgl.Run(run)
}
