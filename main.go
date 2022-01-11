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
	currBounds    pixel.Rect
	initialCenter pixel.Vec
	lastPos       pixel.Vec
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
		currBounds = currBounds.Moved(lastPos.To(pos))

		ctrl.Update() // - here we capture control signals, so actor physics receive input from controller
		hero.Update(dt)
		w.Update(currBounds)

		w.Draw(win)
		hero.Draw(win)

		lastPos = pos
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
	w = world.New("my.tmx")
	w.SetGravity(config.WorldConfig.Gravity)
	currBounds = w.Data()

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

	playerRect := pixel.R(0, 0, config.PlayerConfig.Width, config.PlayerConfig.Height)
	initialCenter = currBounds.Center()

	playerAnims := animation.New(playerRect, config.PlayerConfig.Margin)
	for _, anim := range config.PlayerConfig.Anims {
		playerAnims.SetAnim(anim.Name, anim.File, anim.Frames)
	}

	hero = actor.New(w, playerAnims, playerRect.ResizedMin(pixel.Vec{config.PlayerConfig.Width * 1.25, config.PlayerConfig.Height * 1.25}), config.PlayerConfig.Run, config.PlayerConfig.Walk)
	hero.Move(initialCenter)
	ctrl.Subscribe(hero)

	lastPos = hero.GetPos()

	gameLoop(win)
}

func main() {
	pixelgl.Run(run)
}
