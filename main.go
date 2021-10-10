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
	for _, e := range config.WorldConfig.Enemies {
		world.enemies = append(world.enemies, NewEnemy(*e, config.WorldConfig))
	}

	ctrl := controller.New(win)
	phys := NewPhys()
	phys.rect = pixel.R(0, 0, config.PlayerConfig.Width/2, config.PlayerConfig.Height*0.75)
	phys.runSpeed = config.PlayerConfig.Run
	phys.walkSpeed = config.PlayerConfig.Walk
	phys.jumpSpeed = config.WorldConfig.Gravity * 50
	phys.gravity = config.WorldConfig.Gravity

	// 	id:    1,
	// 	phys:  &phys,
	// 	rect:  pixel.R(0, 0, config.PlayerConfig.Width, config.PlayerConfig.Height),
	// 	anims: make(map[string]*Anim, 0),
	// 	pos:   pixel.V(0.0, 0.0),
	// 	dir:   1.0,
	// }

	anims := animation.New(pixel.R(0, 0, config.PlayerConfig.Width, config.PlayerConfig.Height))
	for _, anim := range config.PlayerConfig.Anims {
		anims.SetAnim(anim.Name, anim.File, anim.Frames)
	}

	hero := actor.New(&phys, anims)
	ctrl.Subscribe(hero)

	initialCenter := win.Bounds().Center()
	phys.rect = phys.rect.Moved(initialCenter)
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

		pos := hero.getPos()
		camPos = pixel.Lerp(camPos, initialCenter.Sub(pos), 1-math.Pow(1.0/128, dt))
		cam := pixel.IM.Moved(camPos)

		win.SetMatrix(cam)
		currBounds = cfg.Bounds.Moved(initialCenter.Sub(pos).Scaled(-1))

		world.Update(currBounds)
		for _, e := range world.currenm {
			e.ai.Update(pos, world.Objects())
			e.p.update(dt, e.ai.vec, world.Objects())
			e.a.Update(dt, e.ai.cmd)
			//			atks.Add(*e.ai.attack)
			e.p.draw(win)
		}

		ctrl.Update(dt, win) // - here we capture control signals, so physics receive input from controller
		//		atks.Add(*ctrl.attack)
		atks.Update(dt, world.Objects())
		phys.update(dt, ctrl.vec, world.Objects())

		hero.Update(dt, ctrl.cmd)

		ctrl.SetGround(phys.ground)

		phys.draw(win)
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
