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

	"github.com/shinomontaz/pixel"
	"github.com/shinomontaz/pixel/pixelgl"
)

// main game loop and logic implementation

var (
	frames   = 0
	second   = time.Tick(time.Second)
	rgba     = color.RGBA{123, 175, 213, 1}
	deltaVec = pixel.ZV
)

func gameFunc(win *pixelgl.Window, dt float64) {
	win.Clear(rgba)

	pos := hero.GetPos()
	sound.Update(pos)
	if dt > 0 {
		deltaVec = pixel.Lerp(deltaVec, lastPos.To(pos), 1-math.Pow(1.0/128, dt))
		currBounds = currBounds.Moved(deltaVec)

		w.Update(currBounds.Moved(pixel.ZV.Add(pixel.V(0, 150))), dt)
	}

	//	b.Draw(win, pos)
	w.Draw(win, pos)
	u.Draw(win)

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
	w = world.New("my.tmx", currBounds)
	w.InitEnemies()

	magic.SetWorld(w)

	initialCenter = w.GetCenter()
	currBounds = currBounds.Moved(initialCenter.Sub(pixel.V(currBounds.W()/2, currBounds.H()/2)))

	win.SetBounds(currBounds)

	hero = factories.NewActor(config.Profiles["player"], w)
	hero.Move(initialCenter)
	ctrl.Subscribe(hero)
	w.AddHero(hero)

	u = ui.New(hero, currBounds)
	lastPos = hero.GetPos()
	b = background.New(lastPos, currBounds.Moved(pixel.Vec{0, 100}), "assets/gamebackground.png")
	w.SetBackground(b)
}
