package menu

import (
	"platformer/events"

	"golang.org/x/image/colornames"

	"github.com/shinomontaz/pixel"
	"github.com/shinomontaz/pixel/imdraw"
)

type Menu struct {
	items    []*Item
	curr     int
	rect     pixel.Rect
	isactive bool
	onquit   func()
	title    string
	logo     *pixel.Sprite
	marginY  float64
	marginX  float64
	imd      *imdraw.IMDraw
}

func New(r pixel.Rect, opts ...MenuOption) *Menu {
	m := &Menu{
		items:   make([]*Item, 0),
		rect:    r,
		marginY: 30,
		marginX: 10,
		imd:     imdraw.New(nil),
	}

	for _, opt := range opts {
		opt(m)
	}

	return m
}

func (m *Menu) GetRect() pixel.Rect {
	return m.rect
}

func (m *Menu) SetActive(a bool) {
	m.isactive = a
	m.updateImdraw()
	// if m.isactive {
	// 	m.Select(0)
	// }
}

func (m *Menu) UpdateSelectedItemText(title string) {
	m.items[m.curr].title = title
}

func (m *Menu) Select(idx int) {
	m.items[m.curr].Select(false)
	m.curr = idx
	m.items[m.curr].Select(true)
}

func (m *Menu) Listen(e int, v pixel.Vec) {
	if !m.isactive {
		return
	}

	// if up or down - handle just here, otherwise make item handle it
	curr := m.curr
	ismoved := false
	if v.Y > 0 {
		curr = (curr - 1 + len(m.items)) % len(m.items)
		ismoved = true
	}
	if v.Y < 0 {
		curr = (curr + 1) % len(m.items)
		ismoved = true
	}
	if ismoved {
		m.Select(curr)
		return
	}
	if e == events.ENTER {
		m.items[m.curr].Action()
		return
	}

	m.items[m.curr].Listen(e, v)
}

func (m *Menu) AddItem(it *Item) {
	offsetY := m.rect.Max.Y - 2*m.marginY
	if m.logo != nil {
		offsetY -= m.logo.Frame().H() / 10
	}

	offsetX := 40.0
	b := it.Bounds()
	i := len(m.items)

	pos := pixel.V(m.rect.Min.X+m.marginX+offsetX, offsetY-float64(i)*b.H())

	it.Place(pos)

	m.updateByHeight(pos.Y - m.marginY)

	m.items = append(m.items, it)
}

func (m *Menu) ReplaceItem(idx int, it *Item) {
	pos := m.items[idx].GetPlace()
	it.Place(pos)
	m.items[idx] = it
}

func (m *Menu) Update(dt float64) {
	if !m.isactive {
		return
	}
	for _, it := range m.items {
		it.Update(dt)
	}
}

func (m *Menu) updateByHeight(h float64) {
	m.rect.Min.Y = h
}

func (m *Menu) updateImdraw() {
	if len(m.items) == 0 {
		return
	}

	offsetY := m.rect.Max.Y - m.marginY
	if m.logo != nil {
		offsetY -= m.logo.Frame().H() / 10
	}

	b := m.items[0].Bounds()
	offsetY = offsetY - float64(len(m.items))*b.H() - m.marginY

	m.updateByHeight(offsetY)

	m.imd.Clear()
	m.imd.Color = colornames.Darkslategray
	m.imd.Push(m.rect.Min)
	m.imd.Push(m.rect.Max)
	m.imd.Rectangle(0)

	m.imd.Color = colornames.Darkgray
	m.imd.Push(m.rect.Min.Add(pixel.Vec{3, 3}))
	m.imd.Push(m.rect.Max.Sub(pixel.Vec{3, 3}))
	m.imd.Rectangle(2)
}

func (m *Menu) Draw(t pixel.Target) {
	if !m.isactive {
		return
	}

	m.imd.Draw(t)

	if m.logo != nil {
		vec := pixel.V(m.rect.Center().X, m.rect.Max.Y+m.logo.Frame().H()/10)
		m.logo.Draw(t, pixel.IM.Moved(vec))
	}

	for _, it := range m.items {
		it.Draw(t)
	}
}

func (m *Menu) Invoke(r pixel.Rect, opts ...MenuOption) {
	m.rect = r
	for _, o := range opts {
		o(m)
	}

	m.updateImdraw()
}
