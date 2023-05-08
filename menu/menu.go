package menu

import (
	"platformer/events"

	"github.com/shinomontaz/pixel"
)

type Menu struct {
	items    []*Item
	curr     int
	rect     pixel.Rect
	isactive bool
	onquit   func()
}

func New(r pixel.Rect) *Menu {
	return &Menu{
		items: make([]*Item, 0),
		rect:  r,
	}
}

func (m *Menu) GetRect() pixel.Rect {
	return m.rect
}

func (m *Menu) SetActive(a bool) {
	m.isactive = a
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
	offsetY := m.rect.H() / 1.6
	b := it.Bounds()
	i := len(m.items)
	c := m.rect.Center()

	pos := pixel.V(c.X-m.rect.W()/16, offsetY-float64(i)*b.H())

	it.Place(pos)

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

func (m *Menu) Draw(t pixel.Target) {
	if !m.isactive {
		return
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
}
