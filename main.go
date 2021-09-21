package main

import (
	"fmt"
	"image/color"
	"math"
	"time"

	"platformer/config"

	"github.com/faiface/pixel"
	"github.com/faiface/pixel/pixelgl"
)

func init() {
	// create world
	// create physics
	// create hero
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

	ctrl := PlayerController{}

	phys := phys{
		rect:      pixel.R(0, 0, config.PlayerConfig.Width, config.PlayerConfig.Height),
		runSpeed:  config.PlayerConfig.Run,
		walkSpeed: config.PlayerConfig.Walk,
		jumpSpeed: config.WorldConfig.Gravity * 50,
		ground:    true,
		gravity:   config.WorldConfig.Gravity,
	}

	hero := Hero{
		phys:  &phys,
		rect:  phys.rect,
		anims: make(map[string]*Anim, 0),
		pos:   pixel.V(0.0, 0.0),
		dir:   1.0,
	}
	for _, anim := range config.PlayerConfig.Anims {
		hero.SetAnim(anim.Name, anim.File, anim.Frames)
	}

	initialCenter := win.Bounds().Center()
	phys.rect = phys.rect.Moved(initialCenter)

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

		pos := hero.getPos()
		currCenter := initialCenter.Add(pixel.Vec{X: -pos.X, Y: -pos.Y})
		camPos = pixel.Lerp(camPos, currCenter, 1-math.Pow(1.0/128, dt))

		cam := pixel.IM.Moved(camPos)

		win.SetMatrix(cam)

		world.Update(cfg.Bounds.Moved(currCenter))
		ctrl.Update(win) // - here we capture control signals, so physics receive input from controller
		phys.update(dt, ctrl.vec, world.Objects())
		hero.Update(dt)

		ctrl.SetGround(phys.ground)

		world.Draw(win)
		hero.draw(win)

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
