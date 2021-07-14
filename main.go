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
	fmt.Println(config.WorldConfig)
	// Platforms [][]int `json:"platforms"`
	// Width     int     `json:"width"`
	// Heigth    int     `json:"heigth"`
	// Gravity   int     `json:"gravity"`
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

	//	world := NewWorld(WIDTH, HEIGTH, hero.getPos())

	phys := phys{
		rect:      pixel.R(0, 0, config.PlayerConfig.Width, config.PlayerConfig.Height),
		runSpeed:  config.PlayerConfig.Walk,
		walkSpeed: config.PlayerConfig.Walk,
		ground:    true,
	}
	hero := Hero{
		phys:  &phys,
		rect:  phys.rect,
		anims: make(map[string]*Anim, 0),
	}
	for _, anim := range config.PlayerConfig.Anims {
		hero.SetAnim(anim.Name, anim.File, anim.Frames)
	}
	pltfms := make([]*platform, 0)
	for _, p := range config.WorldConfig.Platforms {
		pltfms = append(pltfms, NewPlatform(pixel.R(p[0], p[1], p[2], p[3])))
	}

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

		camPos = pixel.Lerp(camPos, hero.getPos(), 1-math.Pow(1.0/128, dt))
		cam := pixel.IM.Moved(camPos)

		win.SetMatrix(cam)
		//		ctrl.update(win) - here we capture control signals, so physics receive input from controller

		//		world.Draw(win, hero.getPos(), camPos)
		//		hero.Update(world.Platforms)
		//		hero.Draw(win, win.Bounds().Center().Sub(hero.getPos()).Sub(pixel.V(0, 140.0)))

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
