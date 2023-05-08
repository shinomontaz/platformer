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
	soundmenu   *menu.Menu
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
		ctrl:        controller.New(win, true),
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
		m.displaymenu.SetActive(true)
		m.mainmenu.SetActive(false)
		m.displaymenu.Select(0)
	}))
	m.mainmenu.AddItem(it)

	txt = text.New(pixel.V(0, 0), m.atlas)
	it = menu.NewItem("Sound", txt, menu.WithAction(func() {
		m.mainmenu.SetActive(false)
		m.soundmenu.SetActive(true)
		m.soundmenu.Select(0)
		m.activemenu = m.soundmenu
	}))
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
			fmt.Println("resolution!!!")
		}))
	m.displaymenu.AddItem(it)
	it.Select(true)

	txt = text.New(pixel.V(0, 0), m.atlas)
	it = menu.NewItem(fmt.Sprintf("%v: %-10v", "Fullscreen", config.Opts.Fullscreen), txt,
		menu.WithHandle(func(e int, v pixel.Vec) {
			if v.X == 0 {
				return
			}
			config.Opts.Fullscreen = !config.Opts.Fullscreen
			m.displaymenu.UpdateSelectedItemText(fmt.Sprintf("%v: %-10v", "Fullscreen", config.Opts.Fullscreen))
		}),
	)
	m.displaymenu.AddItem(it)

	txt = text.New(pixel.V(0, 0), m.atlas)
	it = menu.NewItem("Quit", txt, menu.WithAction(func() {
		m.Notify(events.GAMEVENT_INITSCREEN)
		m.saveOptions()
		m.displaymenu.SetActive(false)
		m.mainmenu.SetActive(true)
		m.activemenu = m.mainmenu
	}))
	m.displaymenu.AddItem(it)

	// sound menu
	m.soundmenu = menu.NewSound(m.currBounds, m.atlas, menu.WithQuit(soundMenuQuit(m)))

	m.ctrl.AddListener(m)
	m.isReady = true
}

func (m *Menu) Listen(e int, v pixel.Vec) {
	m.activemenu.Listen(e, v)
}

func (m *Menu) Start() {
	if !m.isReady {
		return
	}

	m.mainmenu.Select(0)

	m.soundmenu.Invoke(m.currBounds, menu.WithQuit(soundMenuQuit(m)))
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

	//	fmt.Println(m.activemenu)

	m.activemenu.Draw(win)
}

func (m *Menu) saveOptions() {
	go config.SaveRuntime()
}

func soundMenuQuit(m *Menu) func() {
	return func() {
		m.Notify(events.GAMEVENT_UPDATEVOLUME)
		m.saveOptions()
		m.soundmenu.SetActive(false)
		m.mainmenu.SetActive(true)
		m.activemenu = m.mainmenu
	}
}
