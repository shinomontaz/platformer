package submenu

import (
	"fmt"
	"image/color"
	"platformer/common"
	"platformer/events"

	"github.com/faiface/pixel"
	"github.com/faiface/pixel/text"
	"golang.org/x/image/colornames"
)

type Main struct {
	id           int
	m            Menu
	items        []Item
	curr         int
	defaultColor color.Color
	selectColor  color.Color
	atlas        *text.Atlas
	isfullscreen bool
}

func NewMain(m Menu) *Main {
	fnt := common.GetFont("menu") // TODO: use special "menu" font
	msm := &Main{
		id:           MAIN,
		m:            m,
		items:        make([]Item, 0),
		defaultColor: colornames.Whitesmoke,
		selectColor:  colornames.Red,
		atlas:        text.NewAtlas(fnt, text.ASCII),
	}

	msm.addItem("New game", func() {
		fmt.Println("New game")
	})
	msm.addItem("Video", func() {
		//		msm.isfullscreen = !msm.isfullscreen
		//		msm.m.OnVideo(msm.isfullscreen)
		msm.m.SetState(VIDEO)
	})
	msm.addItem("Sound", func() {
		msm.m.SetState(SOUND)
	})
	msm.addItem("Quit", func() {
		fmt.Println("quit")
	})

	return msm
}

func (m *Main) Start() {

}

func (m *Main) GetId() int {
	return m.id
}

func (m *Main) addItem(name string, action func()) {
	r := m.m.GetRect()
	txt := text.New(r.Min, m.atlas)
	txt.Color = m.defaultColor
	item := Item{
		action: action,
		txt:    txt,
		title:  name,
	}
	m.items = append(m.items, item)
}

func (m *Main) Update(dt float64) {

}

func (m *Main) Listen(e int, v pixel.Vec) {
	// update m.curr or call action of current element
	if v.Y > 0 {
		m.curr = (m.curr - 1 + len(m.items)) % len(m.items)
	}
	if v.Y < 0 {
		m.curr = (m.curr + 1) % len(m.items)
	}
	if e == events.ENTER {
		m.items[m.curr].action()
	}
}

func (m *Main) Draw(t pixel.Target) {
	r := m.m.GetRect()
	w := r.W()
	h := r.H()

	itemScale := 1.0 //w / h / 3
	offsetX := w / 2
	offsetY := 0.0

	for i := range m.items {
		m.items[i].txt.Clear()
		m.items[i].txt.Dot.X -= m.items[i].txt.BoundsOf(m.items[i].title).W() / 2
		m.items[i].txt.Color = m.defaultColor

		if i == m.curr {
			m.items[i].txt.Color = m.selectColor
		}

		fmt.Fprintln(m.items[i].txt, m.items[i].title)

		offsetY = h/1.6 - float64(i)*m.items[i].txt.LineHeight*itemScale
		m.items[i].txt.Draw(t, pixel.IM.Scaled(pixel.ZV, itemScale).Moved(pixel.V(offsetX, offsetY)))
	}
}
