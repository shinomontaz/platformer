package main

import (
	"fmt"
	"platformer/common"
	"platformer/menu"

	"github.com/faiface/pixel"
	"github.com/faiface/pixel/pixelgl"
	"github.com/faiface/pixel/text"
)

var mainmenu *menu.Menu
var displaymenu *menu.Menu
var activemenu *menu.Menu

func initMenu(win *pixelgl.Window) {
	videoModes := pixelgl.PrimaryMonitor().VideoModes()
	currentVideoMode := len(videoModes) - 1
	isFullscreen := false

	// main menu
	mainmenu = menu.New(currBounds)
	activemenu = mainmenu

	fnt := common.GetFont("menu")
	atlas := text.NewAtlas(fnt, text.ASCII)

	txt := text.New(pixel.V(0, 0), atlas)
	it := menu.NewItem("New game", txt, menu.WithAction(func() {
		ismenu = false
		mainmenu.SetActive(false)
	}))

	mainmenu.AddItem(it)
	it.Select(true)

	txt = text.New(pixel.V(0, 0), atlas)
	it = menu.NewItem("Display", txt, menu.WithAction(func() {
		activemenu = displaymenu
		mainmenu.SetActive(false)
		activemenu.SetActive(true)
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
	it = menu.NewItem(fmt.Sprintf("%20v: %-10v", "Resolution", fmt.Sprintf("%v x %v", mode.Width, mode.Height)), txt, menu.WithAction(func() {
		fmt.Println("action!!!")
	}))
	displaymenu.AddItem(it)
	it.Select(true)

	txt = text.New(pixel.V(0, 0), atlas)
	it = menu.NewItem(fmt.Sprintf("%20v: %-10v", "Fullscreen", isFullscreen), txt)
	displaymenu.AddItem(it)

	txt = text.New(pixel.V(0, 0), atlas)
	it = menu.NewItem("Quit", txt, menu.WithAction(func() {
		activemenu = mainmenu
		displaymenu.SetActive(false)
		activemenu.SetActive(true)
	}))
	displaymenu.AddItem(it)

	activemenu.SetActive(true)

	ctrl.Subscribe(mainmenu)
	ctrl.Subscribe(displaymenu)
}

func menuFunc(win *pixelgl.Window, dt float64) {
	win.Clear(rgba)

	activemenu.Update(dt)
	// draw menu background
	activemenu.Draw(win)
}
