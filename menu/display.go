package menu

import (
	"fmt"
	"image/color"
	"platformer/common"
	"platformer/events"

	"github.com/faiface/pixel"
	"github.com/faiface/pixel/pixelgl"
	"github.com/faiface/pixel/text"
	"golang.org/x/image/colornames"
)

type Display struct {
	parent       Menu
	id           int
	items        []Item
	curr         int
	defaultColor color.Color
	selectColor  color.Color
	atlas        *text.Atlas
	isfullscreen bool
	currMode     int
	videoModes   []pixelgl.VideoMode
	txt          *text.Text
}

func NewVideo(m Menu) *Video {
	fnt := common.GetFont("menu") // TODO: use special "menu" font
	r := m.GetRect()
	msm := &Video{
		id:           VIDEO,
		m:            m,
		items:        make([]Item, 0),
		defaultColor: colornames.Whitesmoke,
		selectColor:  colornames.Red,
		atlas:        text.NewAtlas(fnt, text.ASCII),
		txt:          text.New(r.Min, text.NewAtlas(fnt, text.ASCII)),
	}
	msm.videoModes = pixelgl.PrimaryMonitor().VideoModes()

	return msm
}

func (m *Video) GetId() int {
	return m.id
}

func (m *Video) Start() {

}

func (m *Video) Update(dt float64) {

}

func (m *Video) Listen(e int, v pixel.Vec) {
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

func (m *Video) Draw(t pixel.Target) {
	r := m.m.GetRect()
	w := r.W()
	h := r.H()

	fntsmall := common.GetFont("menusmall")
	smallatlas := text.NewAtlas(fntsmall, text.ASCII)

	offsetX := w / 2

	// draw title

	title := "Video options"
	m.txt.Dot.X = w - m.txt.BoundsOf(title).W()/2
	m.txt.Dot.Y = h - 3*m.txt.BoundsOf(title).H()
	fmt.Fprintln(m.txt, title)

	offsetY := m.txt.Dot.Y - m.txt.LineHeight

	m.txt.Draw(t, pixel.IM.Moved(pixel.V(offsetX, offsetY)))

	mode := m.videoModes[m.currentVideoMode]
	txt := text.New(r.Min, smallatlas)
	title := fmt.Sprintf("%v x %v", mode.Width, mode.Height)
	txt.Color = m.defaultColor

	// for i, mode := range m.videoModes {
	// 	txt := text.New(r.Min, smallatlas)
	// 	title := fmt.Sprintf("%v x %v", mode.Width, mode.Height)
	// 	txt.Color = m.defaultColor

	// 	// if i == m.curr {
	// 	// 	m.items[i].txt.Color = m.selectColor
	// 	// }

	// 	fmt.Fprintln(txt, title)

	// 	if i%2 == 0 {
	// 		txt.Dot.X -= txt.BoundsOf(title).W() + 50.0
	// 		offsetY = h/1.6 - float64(i)*txt.LineHeight
	// 	} else {
	// 		txt.Dot.X += 50.0
	// 	}

	// 	txt.Draw(t, pixel.IM.Moved(pixel.V(offsetX, offsetY)))
	// }
}
