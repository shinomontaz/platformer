package gamestate

import (
	"image/color"
	"platformer/actor"
	"platformer/common"
	"platformer/config"
	"platformer/controller"
	"platformer/events"
	"platformer/menu"
	"platformer/world"

	"github.com/shinomontaz/pixel/pixelgl"
	"github.com/shinomontaz/pixel/text"

	"github.com/shinomontaz/pixel"
)

var fragSource3 = `
#version 330 core

in vec2  vTexCoords;

out vec4 fragColor;

uniform vec4 uTexBounds;
uniform sampler2D uTexture;
uniform float uTime;

void main() {
	// Get our current screen coordinate
	vec2 t = (vTexCoords - uTexBounds.xy) / uTexBounds.zw;

	// Sum our 3 color channels
	float sum  = texture(uTexture, t).r;
	      sum += texture(uTexture, t).g;
	      sum += texture(uTexture, t).b;

	// Divide by 3, and set the output to the result
	vec4 color = vec4( sum/3, sum/3, sum/3, 1.0);
	fragColor = color;
}
`

type Menu struct {
	Common
	cnv        *pixelgl.Canvas
	uTime      float32
	ingamemenu *menu.Menu

	displaymenu *menu.Menu
	soundmenu   *menu.Menu

	activemenu *menu.Menu
	atlas      *text.Atlas

	ctrl *controller.Controller
}

func NewMenu(game Gamer, w *world.World, hero *actor.Actor, win *pixelgl.Window) *Menu {
	m := &Menu{
		Common: Common{
			game:    game,
			id:      MENU,
			w:       w,
			hero:    hero,
			lastPos: pixel.ZV,
		},
		ctrl: controller.New(win, true),
	}

	m.cnv = pixelgl.NewCanvas(win.Bounds())
	m.cnv.SetSmooth(true)

	m.cnv.SetUniform("uTime", &m.uTime)
	m.cnv.SetFragmentShader(fragSource3)

	menurect := pixel.R(0, 0, 240, 240)
	m.ingamemenu = menu.New(menurect.Moved(win.Bounds().Center().Sub(menurect.Center())))
	m.activemenu = m.ingamemenu

	fnt := common.GetFont("menu28")
	m.atlas = text.NewAtlas(fnt, text.ASCII)

	txt := text.New(pixel.V(0, 0), m.atlas)
	it := menu.NewItem("Resume game", txt, menu.WithAction(func() {
		m.game.SetState(NORMAL)
	}))

	m.ingamemenu.AddItem(it)

	txt = text.New(pixel.V(0, 0), m.atlas)
	it = menu.NewItem("Save", txt)
	m.ingamemenu.AddItem(it)

	txt = text.New(pixel.V(0, 0), m.atlas)
	it = menu.NewItem("Sound", txt, menu.WithAction(func() {
		m.ingamemenu.SetActive(false)
		m.soundmenu.SetActive(true)
		m.soundmenu.Select(0)
		m.activemenu = m.soundmenu
	}))
	m.ingamemenu.AddItem(it)

	txt = text.New(pixel.V(0, 0), m.atlas)
	it = menu.NewItem("Quit to main", txt, menu.WithAction(func() {
		m.game.Notify(events.STAGEVENT_QUIT)
	}))
	m.ingamemenu.AddItem(it)

	// sound menu
	soundmenurect := pixel.R(0, 0, 200, 240)
	m.soundmenu = menu.NewSound(soundmenurect.Moved(win.Bounds().Center().Sub(soundmenurect.Center())), m.atlas, menu.WithQuit(soundMenuQuit(m)))

	m.ctrl.AddListener(m)

	return m
}

func (m *Menu) Update(dt float64) {
	if dt > 0 {
		m.ctrl.Update()
		m.w.Update(m.currBounds, m.lastPos, dt)
		m.uTime = float32(dt)
	}
}

func (m *Menu) Draw(win *pixelgl.Window) {
	camPos := m.lastPos.Add(pixel.V(0, 150))

	m.cnv.Clear(color.RGBA{0, 0, 0, 1})

	m.w.Draw(m.cnv, m.lastPos, camPos, m.cnv.Bounds().Center())

	m.cnv.Draw(win, pixel.IM.Moved(win.Bounds().Center()))
	m.activemenu.Draw(win)
}

func (d *Menu) GetId() int {
	return d.id
}

func (m *Menu) Start() {
	m.currBounds = m.w.GetViewport()
	m.lastPos = m.hero.GetPos()

	m.ingamemenu.Select(0)
	m.soundmenu.Invoke(m.currBounds, menu.WithQuit(soundMenuQuit(m)))

	m.activemenu.SetActive(true)
}

func (m *Menu) KeyEvent(key pixelgl.Button) {
	m.activemenu.KeyEvent(key)
}

func (m *Menu) saveOptions() {
	go config.SaveRuntime()
}

func soundMenuQuit(m *Menu) func() {
	return func() {
		m.game.Notify(events.GAMEVENT_UPDATEVOLUME)
		m.saveOptions()
		m.soundmenu.SetActive(false)
		m.ingamemenu.SetActive(true)
		m.activemenu = m.ingamemenu
	}
}
