package stages

import (
	"fmt"
	"platformer/common"
	"platformer/config"
	"platformer/controller"
	"platformer/events"
	"platformer/menu"

	"github.com/shinomontaz/pixel"
	"github.com/shinomontaz/pixel/pixelgl"
	"github.com/shinomontaz/pixel/text"
)

type Menu struct {
	Common
	assetloader *common.Loader
	mainmenu    *menu.Menu
	displaymenu *menu.Menu
	activemenu  *menu.Menu

	mainmenuback *menu.Back

	ctrl       *controller.Controller
	currBounds pixel.Rect
	isGame     bool
	atlas      *text.Atlas
}

func NewMenu(f Inform, l *common.Loader, win *pixelgl.Window, currBounds pixel.Rect) *Menu {
	return &Menu{
		Common: Common{
			id:       MENU,
			done:     make(chan struct{}),
			inform:   f,
			eventMap: map[int]int{events.STAGEVENT_NEXT: GAME},
		},
		assetloader: l,
		ctrl:        controller.New(win),
		currBounds:  currBounds,
	}
}

func (m *Menu) Init() {
	videoModes := pixelgl.PrimaryMonitor().VideoModes()
	currentVideoMode := len(videoModes) - 1

	m.mainmenuback = menu.NewBack(m.currBounds, m.assetloader)

	// main menu
	m.mainmenu = menu.New(m.currBounds)
	m.activemenu = m.mainmenu

	fnt := common.GetFont("menu")
	m.atlas = text.NewAtlas(fnt, text.ASCII)

	txt := text.New(pixel.V(0, 0), m.atlas)
	it := menu.NewItem("New game", txt, menu.WithAction(func() {
		m.Notify(events.STAGEVENT_NEXT)
		m.isGame = true
	}))

	m.mainmenu.AddItem(it)

	txt = text.New(pixel.V(0, 0), m.atlas)
	it = menu.NewItem("Display", txt, menu.WithAction(func() {
		m.activemenu = m.displaymenu
		m.mainmenu.SetActive(false)
		m.displaymenu.SetActive(true)
	}))
	m.mainmenu.AddItem(it)

	txt = text.New(pixel.V(0, 0), m.atlas)
	it = menu.NewItem("Sound", txt)
	m.mainmenu.AddItem(it)

	txt = text.New(pixel.V(0, 0), m.atlas)
	it = menu.NewItem("Quit", txt, menu.WithAction(func() {
		m.Notify(events.STAGEVENT_QUIT)
	}))
	m.mainmenu.AddItem(it)

	// display menu
	m.displaymenu = menu.New(m.currBounds)
	txt = text.New(pixel.V(0, 0), m.atlas)

	mode := videoModes[currentVideoMode]
	it = menu.NewItem(fmt.Sprintf("%v: %-10v", "Resolution", fmt.Sprintf("%v x %v", mode.Width, mode.Height)), txt,
		menu.WithAction(func() {
			fmt.Println("action!!!")
		}))
	m.displaymenu.AddItem(it)
	it.Select(true)

	txt = text.New(pixel.V(0, 0), m.atlas)
	it = menu.NewItem(fmt.Sprintf("%v: %-10v", "Fullscreen", config.Opts.Fullscreen), txt,
		menu.WithHandle(func(e int, v pixel.Vec) {
			config.Opts.Fullscreen = !config.Opts.Fullscreen
			m.displaymenu.UpdateSelectedItemText(fmt.Sprintf("%v: %-10v", "Fullscreen", config.Opts.Fullscreen))
		}),
		menu.WithAction(func() {
			m.Notify(events.GAMEVENT_INITSCREEN)
			m.saveOptions()
		}))
	m.displaymenu.AddItem(it)

	txt = text.New(pixel.V(0, 0), m.atlas)
	it = menu.NewItem("Quit", txt, menu.WithAction(func() {
		m.activemenu = m.mainmenu
		m.displaymenu.SetActive(false)
		m.mainmenu.SetActive(true)
	}))
	m.displaymenu.AddItem(it)

	m.ctrl.AddListener(m.mainmenu)
	m.ctrl.AddListener(m.displaymenu)

	m.isReady = true

}

func (m *Menu) Start() {
	if !m.isReady {
		return
	}

	m.mainmenu.Select(0)

	m.activemenu.SetActive(true)
	m.isActive = true
}

func (m *Menu) Run(win *pixelgl.Window, dt float64) {
	if !m.isReady {
		m.Notify(events.STAGEVENT_NOTREADY)
		return
	}

	m.ctrl.Update()

	win.Clear(rgba)

	m.mainmenuback.Update(dt)
	m.mainmenuback.Draw(win)

	m.activemenu.Update(dt)
	m.activemenu.Draw(win)
}

func (m *Menu) saveOptions() {
	go config.SaveRuntime()
}
