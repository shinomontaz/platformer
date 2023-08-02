package stages

import (
	"fmt"
	"platformer/bindings"
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
	assetloader  *common.Loader
	mainmenu     *menu.Menu
	displaymenu  *menu.Menu
	controlsmenu *menu.Menu
	soundmenu    *menu.Menu
	activemenu   *menu.Menu

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
	logopic, err := m.assetloader.LoadPicture("logo.png")
	logo := pixel.NewSprite(logopic, pixel.R(0, 0, logopic.Bounds().W(), logopic.Bounds().H()))
	if err != nil {
		panic(err)
	}
	soundmenupic, err := m.assetloader.LoadPicture("sound.png")
	soundlogo := pixel.NewSprite(soundmenupic, pixel.R(0, 0, soundmenupic.Bounds().W(), soundmenupic.Bounds().H()))
	if err != nil {
		panic(err)
	}
	displaymenupic, err := m.assetloader.LoadPicture("display.png")
	displaylogo := pixel.NewSprite(displaymenupic, pixel.R(0, 0, displaymenupic.Bounds().W(), displaymenupic.Bounds().H()))
	if err != nil {
		panic(err)
	}

	menurect := pixel.R(0, 0, 200, 240)
	m.mainmenu = menu.New(menurect.Moved(m.currBounds.Center().Sub(menurect.Center())), menu.WithLogo(logo))

	m.activemenu = m.mainmenu

	fnt := common.GetFont("menu28")
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
	it = menu.NewItem("Controls", txt, menu.WithAction(func() {
		m.activemenu = m.controlsmenu
		m.controlsmenu.SetActive(true)
		m.mainmenu.SetActive(false)
		m.controlsmenu.Select(0)
	}))
	m.mainmenu.AddItem(it)

	txt = text.New(pixel.V(0, 0), m.atlas)
	it = menu.NewItem("Quit", txt, menu.WithAction(func() {
		m.Notify(events.STAGEVENT_QUIT)
	}))
	m.mainmenu.AddItem(it)

	// display menu
	displaymenurect := pixel.R(0, 0, 400, 140)
	m.displaymenu = menu.New(displaymenurect.Moved(m.currBounds.Center().Sub(displaymenurect.Center())), menu.WithLogo(displaylogo))

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
		menu.WithHandle(func(b pixelgl.Button) {
			switch b {
			case pixelgl.KeyLeft:
				fallthrough
			case pixelgl.KeyRight:
				config.Opts.Fullscreen = !config.Opts.Fullscreen
			}
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

	m.SetupControlsMenu()

	// sound menu
	soundmenurect := pixel.R(0, 0, 200, 240)
	m.soundmenu = menu.NewSound(soundmenurect.Moved(m.currBounds.Center().Sub(soundmenurect.Center())), m.atlas, menu.WithQuit(soundMenuQuit(m)), menu.WithLogo(soundlogo))

	m.ctrl.AddKeyListener(m)
	m.isReady = true
}

func (m *Menu) SetupControlsMenu() {
	controlsmenurect := pixel.R(0, 0, 400, 140)
	m.controlsmenu = menu.New(controlsmenurect.Moved(m.currBounds.Center().Sub(controlsmenurect.Center())))
	activated := ""
	for _, ka := range bindings.KeyActions {
		ka := ka
		title := bindings.KeyActionNames[bindings.KeyAction[ka]]
		kaId := bindings.KeyAction[ka]
		txt := text.New(pixel.V(0, 0), m.atlas)
		it := menu.NewItem(fmt.Sprintf("%v: %-10v", title, bindings.Active.GetBinding(kaId)), txt,
			menu.WithAction(func() { // on Enter
				if activated == "" {
					activated = ka
					m.ctrl.SetListenAll(true)
					m.controlsmenu.UpdateSelectedItemText(fmt.Sprintf("%v: %-10v", title, "?????"))
					m.controlsmenu.AcceptInput(true)
				}
			}),
			menu.WithHandle(func(b pixelgl.Button) {
				m.ctrl.SetListenAll(false)
				m.controlsmenu.AcceptInput(false)
				if activated == "" {
					return
				}

				bindings.Active.SetBind(b, kaId)
				activated = ""
				m.controlsmenu.UpdateSelectedItemText(fmt.Sprintf("%v: %-10v", title, bindings.Active.GetBinding(kaId)))
				m.ctrl.SetListenAll(false)

			}),
		)
		m.controlsmenu.AddItem(it)
	}

	txt := text.New(pixel.V(0, 0), m.atlas)
	it := menu.NewItem("Quit", txt, menu.WithAction(func() {
		bindings.Active.Save()
		m.saveOptions()
		m.controlsmenu.SetActive(false)
		m.mainmenu.SetActive(true)
		m.activemenu = m.mainmenu
	}))
	m.controlsmenu.AddItem(it)
}

func (m *Menu) KeyAction(key pixelgl.Button) {
	m.activemenu.KeyAction(key)
}

func (m *Menu) Start() {
	if !m.isReady {
		return
	}

	m.mainmenu.Select(0)

	soundmenurect := pixel.R(0, 0, 200, 240)
	m.soundmenu.Invoke(soundmenurect.Moved(m.currBounds.Center().Sub(soundmenurect.Center())), menu.WithQuit(soundMenuQuit(m)))
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

func soundMenuQuit(m *Menu) func() {
	return func() {
		m.Notify(events.GAMEVENT_UPDATEVOLUME)
		m.saveOptions()
		m.soundmenu.SetActive(false)
		m.mainmenu.SetActive(true)
		m.activemenu = m.mainmenu
	}
}
