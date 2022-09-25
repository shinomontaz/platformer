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
	camPos   = pixel.ZV
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

		camPos = pixel.Lerp(camPos, pos.Add(pixel.V(0, 150)), 1-math.Pow(1.0/128, dt))
		currBounds2 := currBounds.Moved(initialCenter.Sub(pixel.V(currBounds.W()/2, currBounds.H()/2))).Moved(deltaVec).Moved(pixel.ZV.Add(pixel.V(0, 150)))

		w.Update(currBounds2, dt)

	}

	w.Draw(win, pos, camPos)
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
	//	w = world.New("assets/ep2.tmx", currBounds)

	w.InitEnemies()

	magic.SetWorld(w)

	initialCenter = w.GetCenter()

	hero = factories.NewActor(config.Profiles["player"], w)
	hero.Move(initialCenter)
	ctrl.Subscribe(hero)
	w.AddHero(hero)

	u = ui.New(hero, currBounds)
	lastPos = hero.GetPos()
	camPos = lastPos.Add(pixel.V(0, 150))

	b = background.New(lastPos, currBounds.Moved(pixel.V(-currBounds.W()/2, 0)).Moved(pixel.V(0, 150)), "assets/gamebackground.png")

	w.SetBackground(b)
}
