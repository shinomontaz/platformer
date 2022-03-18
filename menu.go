package main

import (
	"fmt"
	"platformer/common"
	"platformer/config"
	"platformer/menu"

	"github.com/faiface/pixel"
	"github.com/faiface/pixel/pixelgl"
	"github.com/faiface/pixel/text"
)

var (
	mainmenu    *menu.Menu
	displaymenu *menu.Menu
	activemenu  *menu.Menu

	mainmenuback *menu.Back
)

func initMenu(win *pixelgl.Window) {
	videoModes := pixelgl.PrimaryMonitor().VideoModes()
	currentVideoMode := len(videoModes) - 1

	mainmenuback = menu.NewBack(currBounds)

	// main menu
	mainmenu = menu.New(currBounds)
	activemenu = mainmenu

	fnt := common.GetFont("menu")
	atlas := text.NewAtlas(fnt, text.ASCII)

	txt := text.New(pixel.V(0, 0), atlas)
	it := menu.NewItem("New game", txt, menu.WithAction(func() {
		ismenu = false
		startGame()
		mainmenu.SetActive(false)
	}))

	mainmenu.AddItem(it)
	it.Select(true)

	txt = text.New(pixel.V(0, 0), atlas)
	it = menu.NewItem("Display", txt, menu.WithAction(func() {
		activemenu = displaymenu
		mainmenu.SetActive(false)
		displaymenu.SetActive(true)
	}))
	mainmenu.AddItem(it)

	txt = text.New(pixel.V(0, 0), atlas)
	it = menu.NewItem("Sound", txt)
	mainmenu.AddItem(it)

	txt = text.New(pixel.V(0, 0), atlas)
	it = menu.NewItem("Quit", txt, menu.WithAction(func() {
		isquit = true
	}))
	mainmenu.AddItem(it)

	// display menu
	displaymenu = menu.New(currBounds)
	txt = text.New(pixel.V(0, 0), atlas)

	mode := videoModes[currentVideoMode]
	it = menu.NewItem(fmt.Sprintf("%v: %-10v", "Resolution", fmt.Sprintf("%v x %v", mode.Width, mode.Height)), txt,
		menu.WithAction(func() {
			fmt.Println("action!!!")
		}))
	displaymenu.AddItem(it)
	it.Select(true)

	txt = text.New(pixel.V(0, 0), atlas)
	it = menu.NewItem(fmt.Sprintf("%v: %-10v", "Fullscreen", config.Opts.Fullscreen), txt,
		menu.WithHandle(func(e int, v pixel.Vec) {
			config.Opts.Fullscreen = !config.Opts.Fullscreen
			displaymenu.UpdateSelectedItemText(fmt.Sprintf("%v: %-10v", "Fullscreen", config.Opts.Fullscreen))
		}),
		menu.WithAction(func() {
			// fullscreen toggle
			initScreen(win)
			SaveOptions()
		}))
	displaymenu.AddItem(it)

	txt = text.New(pixel.V(0, 0), atlas)
	it = menu.NewItem("Quit", txt, menu.WithAction(func() {
		activemenu = mainmenu
		displaymenu.SetActive(false)
		mainmenu.SetActive(true)
	}))
	displaymenu.AddItem(it)

	activemenu.SetActive(true)

	ctrl.Subscribe(mainmenu)
	ctrl.Subscribe(displaymenu)
}

func menuFunc(win *pixelgl.Window, dt float64) {
	win.Clear(rgba)

	mainmenuback.Update(dt)
	mainmenuback.Draw(win)

	activemenu.Update(dt)
	activemenu.Draw(win)
}
