package main

import (
	"encoding/json"
	"io/ioutil"
	"platformer/config"
	"platformer/sound"

	"github.com/faiface/pixel/pixelgl"

	"github.com/faiface/pixel"
)

var (
	isfullscreen bool
	winHeight    float64
	winWidth     float64
)

// get window mode, sound volumes and menu/game mode
// currBounds
func initRuntime() {
	winWidth = config.Opts.WindowWidth
	winHeight = config.Opts.WindowHeight

	currBounds = pixel.R(0, 0, winWidth, winHeight)
}

func SaveOptions() {
	json, _ := json.Marshal(config.Opts)
	if err := ioutil.WriteFile("config/options.json", json, 0644); err != nil {
		panic("Failed to save configuration")
	}
}

func startGame() {
	sound.PlayMusic("main")
}

func initScreen(win *pixelgl.Window) {
	if config.Opts.Fullscreen {
		win.SetMonitor(pixelgl.PrimaryMonitor())
		win.SetBounds(currBounds)
	} else {
		win.SetMonitor(nil)
		win.SetBounds(currBounds)
	}
}
