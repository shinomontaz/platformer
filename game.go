package main

import (
	"fmt"
	"image/color"
	"platformer/actor"
	"platformer/animation"
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
		//		deltaVec = pixel.Lerp(deltaVec, lastPos.To(pos), 1-math.Pow(1.0/128, dt))

		deltaVec = lastPos.To(pos)
		//		camPos = pixel.Lerp(camPos, pos.Add(pixel.V(0, 150)), 1-math.Pow(1.0/128, dt))
		camPos = pos.Add(pixel.V(0, 150))

		currBounds = currBounds.Moved(deltaVec)

		w.Update(currBounds, dt)
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

func loadAnimations() { // load animations
	for _, anim := range config.AnimConfig {
		animation.Load(&anim)
	}
	for name, cfg := range config.Spells {
		magic.Load(name, &cfg)
	}
}

func initGame(win *pixelgl.Window) *world.World {
	animation.Init(assetloader)
	actor.Init(assetloader)
	ui.Init(assetloader)

	loadAnimations()

	w, err := world.New("ep1.tmx", currBounds, world.WithLoader(assetloader))
	if err != nil {
		panic(err)
	}

	w.IsDebug = isdebug
	w.InitEnemies()

	magic.SetWorld(w)

	initialCenter = w.GetCenter()
	hero = factories.NewActor(config.Profiles["player"], w)
	hero.Move(initialCenter)
	ctrl.Subscribe(hero)
	w.AddHero(hero)
	u = ui.New(hero, currBounds)

	currBounds = currBounds.Moved(initialCenter.Sub(pixel.V(currBounds.W()/2, currBounds.H()/2)))

	lastPos = hero.GetPos()
	camPos = lastPos.Add(pixel.V(0, 150))

	b = background.New(lastPos, currBounds.Moved(pixel.V(0, 150)), assetloader, "gamebackground.png")

	w.SetBackground(b)

	return w
}
