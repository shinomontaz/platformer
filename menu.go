package main

import (
	"platformer/menu"

	"github.com/faiface/pixel/pixelgl"
)

var mainmenu *menu.Menu

func initMenu(win *pixelgl.Window) {
	mainmenu = menu.NewMain(currBounds)

	// displaymenu = menu.NewMain(currBounds)
	// soundmenu = menu.NewMain(currBounds)

	mainmenu.VideoCallback(func(fs bool) {
		isfullscreen = fs
		if isfullscreen {
			mons := pixelgl.Monitors()
			if len(mons) > 0 {
				win.SetMonitor(mons[0])
			}
		} else {
			win.SetMonitor(nil)
		}

		win.SetBounds(currBounds)
	})

	ctrl.Subscribe(mainmenu)
}

func menuFunc(win *pixelgl.Window, dt float64) {
	win.Clear(rgba)

	mainmenu.Update(dt)
	mainmenu.Draw(win)
}
