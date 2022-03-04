package main

import (
	"fmt"
	"image/color"
	"math"
	"platformer/animation"
	"platformer/common"
	"platformer/menu"

	"github.com/faiface/pixel"
	"github.com/faiface/pixel/pixelgl"
	"github.com/faiface/pixel/text"
	"golang.org/x/image/colornames"
)

var (
	mainmenu       *menu.Menu
	displaymenu    *menu.Menu
	activemenu     *menu.Menu
	mainmenurba    color.Color
	anims          *animation.Anims
	animSpriteNum  int
	currtime       float64
	campfiresprite *pixel.Sprite
)

func initMenu(win *pixelgl.Window) {
	videoModes := pixelgl.PrimaryMonitor().VideoModes()
	currentVideoMode := len(videoModes) - 1
	isFullscreen := false
	mainmenurba = colornames.Black

	anims = animation.Get("scenery")
	campfiresprite = pixel.NewSprite(nil, pixel.Rect{})

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
	win.Clear(mainmenurba)

	currtime += dt
	animSpriteNum = int(math.Floor(currtime / 0.2))
	pic, rect := anims.GetSprite("campfire", animSpriteNum)
	campfiresprite.Set(pic, rect)
	c := currBounds.Min
	campfiresprite.Draw(win, pixel.IM.Moved(pixel.V(c.X+100, c.Y+100)))

	activemenu.Update(dt)
	activemenu.Draw(win)
}
