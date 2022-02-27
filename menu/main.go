package menu

import (
	"platformer/common"
	"platformer/menu/submenu"

	"github.com/faiface/pixel"
)

type Menu struct {
	currmenu Menuer
	submenus map[int]Menuer
	rect     pixel.Rect
	sbrs     []common.Subscriber
	vc       func(bool) // video mode callback
	smc      func(int)  // sound music callback
	sec      func(int)  // sound effects callback
}

func NewMain(r pixel.Rect) *Menu {
	m := &Menu{
		rect: r,
		sbrs: make([]common.Subscriber, 0),
	}

	smMain := submenu.NewMain(m)
	smVideo := submenu.NewVideo(m)

	m.submenus = map[int]Menuer{
		submenu.MAIN:  smMain,
		submenu.VIDEO: smVideo,
	}

	m.SetState(submenu.MAIN)

	return m
}

func (m *Menu) inform(e int, v pixel.Vec) {
	for _, s := range m.sbrs {
		s.Listen(e, v)
	}
}

func (m *Menu) AddListener(s common.Subscriber) {
	m.sbrs = append(m.sbrs, s)
}

func (m *Menu) SetState(id int) {
	m.currmenu = m.submenus[id]
	m.currmenu.Start()
}

func (m *Menu) GetRect() pixel.Rect {
	return m.rect
}

func (m *Menu) Update(dt float64) {
	m.currmenu.Update(dt)
}

func (m *Menu) Listen(e int, v pixel.Vec) {
	m.currmenu.Listen(e, v)
}

func (m *Menu) Draw(t pixel.Target) {
	m.currmenu.Draw(t)
}

func (m *Menu) VideoCallback(fs func(bool)) {
	m.vc = fs
}

func (m *Menu) AudioCallback(ch string, fs func(int)) {
	switch ch {
	case "music":
		m.smc = fs
	case "effects":
		m.sec = fs
	}
}

func (m *Menu) OnVideo(isfullscreen bool) {
	m.vc(isfullscreen)
}

func (m *Menu) OnAudio(ch string, volume int) {
	switch ch {
	case "music":
		m.smc(volume)
	case "effects":
		m.sec(volume)
	}
}
