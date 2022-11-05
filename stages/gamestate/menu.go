package gamestate

import (
	"image/color"
	"platformer/actor"
	"platformer/common"
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
	win        *pixelgl.Window
	ctrl       *controller.Controller
	cnv        *pixelgl.Canvas
	uTime      float32
	ingamemenu *menu.Menu
	atlas      *text.Atlas
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
		win:  win,
		ctrl: controller.New(win),
	}

	m.cnv = pixelgl.NewCanvas(win.Bounds())
	m.cnv.SetSmooth(true)

	m.cnv.SetUniform("uTime", &m.uTime)
	m.cnv.SetFragmentShader(fragSource3)

	m.ingamemenu = menu.New(win.Bounds())

	fnt := common.GetFont("menu")
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
	it = menu.NewItem("Quit to main", txt, menu.WithAction(func() {
		m.game.Notify(events.STAGEVENT_QUIT)
	}))
	m.ingamemenu.AddItem(it)

	m.ctrl.AddListener(m.ingamemenu)

	return m
}

func (m *Menu) Update(dt float64) {
	if dt > 0 {
		m.ctrl.Update()
		m.w.Update(m.currBounds, dt)
		m.uTime = float32(dt)
	}
}

func (m *Menu) Draw(win *pixelgl.Window) {
	camPos := m.lastPos.Add(pixel.V(0, 150))

	m.cnv.Clear(color.RGBA{0, 0, 0, 1})

	m.w.Draw(m.cnv, m.lastPos, camPos, m.cnv.Bounds().Center())

	m.cnv.Draw(win, pixel.IM.Moved(win.Bounds().Center()))
	m.ingamemenu.Draw(win)
}

func (d *Menu) GetId() int {
	return d.id
}

func (m *Menu) Start() {
	m.currBounds = m.w.GetViewport()
	m.lastPos = m.hero.GetPos()

	m.ingamemenu.Select(0)
	m.ingamemenu.SetActive(true)
}

func (m *Menu) Listen(e int, v pixel.Vec) {
}
