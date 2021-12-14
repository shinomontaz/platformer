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

	world := NewWorld(config.WorldConfig.Width, config.WorldConfig.Heigth)
	for _, p := range config.WorldConfig.Platforms {
		world.platforms = append(world.platforms, NewPlatform(pixel.R(p[0], p[1], p[2], p[3]).Moved(win.Bounds().Center())))
	}

	ctrl := controller.New(win)
	phys := NewPhys(world)
	phys.rect = pixel.R(0, 0, config.PlayerConfig.Width/2, config.PlayerConfig.Height*0.75)
	phys.runSpeed = config.PlayerConfig.Run
	phys.walkSpeed = config.PlayerConfig.Walk
	phys.jumpSpeed = config.WorldConfig.Gravity * 30
	phys.gravity = config.WorldConfig.Gravity

	initialCenter := win.Bounds().Center()
	phys.rect = phys.rect.Moved(initialCenter)

	playerAnims := animation.New(pixel.R(0, 0, config.PlayerConfig.Width, config.PlayerConfig.Height))
	for _, anim := range config.PlayerConfig.Anims {
		playerAnims.SetAnim(anim.Name, anim.File, anim.Frames)
	}

	hero := actor.New(&phys, playerAnims)
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

		camPos = pixel.Lerp(camPos, initialCenter.Sub(pos), 1-math.Pow(1.0/128, dt))
		cam := pixel.IM.Moved(camPos)

		win.SetMatrix(cam)
		currBounds = cfg.Bounds.Moved(initialCenter.Sub(pos).Scaled(-1))

		ctrl.Update() // - here we capture control signals, so actor physics receive input from controller
		hero.Update(dt)
		world.Update(currBounds)

		world.Draw(win)
		hero.Draw(win)
		phys.draw(win)

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
