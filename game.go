package main

import (
	"fmt"
	"image/color"
	"math"
	"platformer/background"
	"platformer/config"
	"platformer/factories"
	"platformer/magic"
	"platformer/sound"
	"platformer/ui"
	"platformer/world"
	"time"

	"github.com/faiface/pixel"
	"github.com/faiface/pixel/pixelgl"
)

// main game loop and logic implementation

var (
	camPos = pixel.ZV
	frames = 0
	second = time.Tick(time.Second)
	rgba   = color.RGBA{123, 175, 213, 1}
)

func gameFunc(win *pixelgl.Window, dt float64) {
	win.Clear(rgba)

	pos := hero.GetPos()
	sound.Update(pos)
	if dt > 0 {
		deltaVec := lastPos.To(pos)
		camPos = pixel.Lerp(camPos, initialCenter.Sub(pos).Sub(pixel.V(0, 150)), 1-math.Pow(1.0/128, dt)) // standart with moving cam slightly down
		cam := pixel.IM.Moved(camPos)

		win.SetMatrix(cam)
		currBounds = currBounds.Moved(deltaVec)

		w.Update(currBounds, dt)
	}

	b.Draw(win, pos, camPos)
	w.Draw(win)
	u.Draw(win, pos, camPos)

	lastPos = pos

	frames++
	select {
	case <-second:
		win.SetTitle(fmt.Sprintf("%s | FPS: %d", title, frames))
		frames = 0
	default:
	}
}

func initGame(win *pixelgl.Window) {
	w = world.New("my.tmx")
	w.InitEnemies()

	magic.SetWorld(w)
	//	sound.PlayMusic("main")
	// mons := pixelgl.Monitors()
	// if len(mons) > 0 {
	// 	win.SetMonitor(mons[0])
	// }

	initialCenter = w.GetCenter()
	currBounds = currBounds.Moved(initialCenter.Sub(pixel.V(currBounds.W()/2, currBounds.H()/2)))
	//	currBounds = currBounds.Moved(pixel.V(200, 0))

	win.SetBounds(currBounds)

	hero = factories.NewActor(config.Profiles["player"], w)
	hero.Move(initialCenter)
	ctrl.Subscribe(hero)
	w.AddHero(hero)

	u = ui.New(hero, currBounds)
	lastPos = hero.GetPos()
	b = background.New(lastPos, currBounds.Moved(pixel.Vec{0, 100}), "assets/gamebackground.png")
}
