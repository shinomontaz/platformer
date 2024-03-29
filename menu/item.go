package menu

import (
	"fmt"
	"image/color"

	"github.com/shinomontaz/pixel"
	"github.com/shinomontaz/pixel/pixelgl"
	"github.com/shinomontaz/pixel/text"
	"golang.org/x/image/colornames"
)

type Item struct {
	title        string
	txt          *text.Text
	defaultColor color.Color
	selectColor  color.Color

	handle func(pixelgl.Button)
	action func()
	rect   pixel.Rect

	selected bool
}

func NewItem(title string, txt *text.Text, opts ...ItemOption) *Item {
	i := &Item{
		title:        title,
		txt:          txt,
		defaultColor: colornames.Whitesmoke,
		selectColor:  colornames.Red,
		handle:       func(b pixelgl.Button) {},
		action:       func() {},
	}

	for _, opt := range opts {
		opt(i)
	}
	return i
}

func (i *Item) Action() {
	i.action()
}

func (i *Item) Select(s bool) {
	i.selected = s
}

func (i *Item) Update(dt float64) {
	// i.txt.Color = i.defaultColor
	// if i.selected {
	// 	i.txt.Color = i.selectColor
	// }
}

func (i *Item) Bounds() pixel.Rect {
	return i.txt.BoundsOf(i.title)
}

func (i *Item) Place(pos pixel.Vec) {
	//	i.txt.Dot = pos
	i.rect = i.txt.Bounds().Moved(pos)
}

func (i *Item) GetPlace() pixel.Vec {
	return i.rect.Center()
}

func (i *Item) KeyAction(b pixelgl.Button) { // for left and right volume change
	i.handle(b)
}

func (i *Item) Draw(t pixel.Target) {
	i.txt.Clear()
	i.txt.Color = i.defaultColor
	if i.selected {
		i.txt.Color = i.selectColor
	}
	fmt.Fprintln(i.txt, i.title)
	i.txt.Draw(t, pixel.IM.Moved(i.rect.Center()))

}
