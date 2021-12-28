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
	"platformer/phys"
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

func gameLoop() {
	cfg := pixelgl.WindowConfig{
		Title:  "platformer",
		Bounds: pixel.R(0, 0, config.WorldConfig.Width, config.WorldConfig.Heigth),
		VSync:  true,
	}

	win, err := pixelgl.NewWindow(cfg)
	if err != nil {
		panic(err)
	}

	win.SetSmooth(true)

	w := world.New(config.WorldConfig.Width, config.WorldConfig.Heigth)
	for _, p := range config.WorldConfig.Platforms {
		w.Add(world.NewPlatform(pixel.R(p[0], p[1], p[2], p[3]).Moved(win.Bounds().Center())))
	}

	ctrl := controller.New(win)

	mainRect := pixel.R(0, 0, config.PlayerConfig.Width, config.PlayerConfig.Height)
	initialCenter := win.Bounds().Center()

	fmt.Println("mainRect, mainRect.Moved(initialCenter)", mainRect, mainRect.Moved(initialCenter))

	p := phys.New(mainRect.Moved(initialCenter), config.PlayerConfig.Run, config.PlayerConfig.Walk, config.WorldConfig.Gravity*30, config.WorldConfig.Gravity)
	p.SetQt(w.GetQt())

	playerAnims := animation.New(mainRect)
	for _, anim := range config.PlayerConfig.Anims {
		playerAnims.SetAnim(anim.Name, anim.File, anim.Frames)
	}

	hero := actor.New(&p, playerAnims)
	ctrl.Subscribe(hero)
	currBounds := cfg.Bounds

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

		//		fmt.Println("hero.GetPos()", pos)

		camPos = pixel.Lerp(camPos, initialCenter.Sub(pos), 1-math.Pow(1.0/128, dt))
		cam := pixel.IM.Moved(camPos)

		win.SetMatrix(cam)
		currBounds = cfg.Bounds.Moved(initialCenter.Sub(pos).Scaled(-1))

		ctrl.Update() // - here we capture control signals, so actor physics receive input from controller
		hero.Update(dt)
		w.Update(currBounds)

		w.Draw(win)
		hero.Draw(win)
		p.Draw(win)

		win.Update()

		frames++
		select {
		case <-frametime:
			//			hero.Tick()
		case <-second:
			win.SetTitle(fmt.Sprintf("%s | FPS: %d", cfg.Title, frames))
			frames = 0
		default:
		}
	}
}

func main() {
	pixelgl.Run(gameLoop)
}
