package inventory

import (
	"fmt"
	"platformer/common"
	"platformer/loot"

	"github.com/shinomontaz/pixel"
	"github.com/shinomontaz/pixel/text"
	"golang.org/x/image/colornames"
)

var (
	limit    int
	list     map[int]int
	piclist  map[int]*pixel.Sprite
	loader   *common.Loader
	atlas    *text.Atlas
	atlasbig *text.Atlas
)

func Init(l *common.Loader) {
	list = make(map[int]int)
	piclist = make(map[int]*pixel.Sprite)
	limit = 10
	loader = l
	atlas = text.NewAtlas(common.GetFont("regular16"), text.ASCII)
	atlasbig = text.NewAtlas(common.GetFont("regular32"), text.ASCII)
}

func HaveCoins() bool {
	_, ok := list[loot.COIN]
	return ok
}

func PayCoins(n int) {
	if _, ok := list[loot.COIN]; ok {
		for i := 0; i < n; i++ {
			if list[loot.COIN] == 0 {
				break
			}
			list[loot.COIN]--
		}
	}
}

func Add(loots []*loot.Loot) {
	if len(loots) == 0 {
		return
	}
	for _, l := range loots {
		id := l.GetId()
		if _, exists := list[id]; !exists {
			list[id] = 0
			piclist[id] = loadPic(l.GetImagePath())
		}
		list[id]++
	}
}

func Remove(l *loot.Loot) {
	id := l.GetId()
	if _, exists := list[id]; !exists {
		return
	}
	list[id]--
	if list[id] <= 0 {
		delete(list, id)
	}
}

func loadPic(path string) *pixel.Sprite {
	if path == "" {
		return nil
	}
	prt, err := loader.LoadPicture(path)
	if err != nil {
		panic(err)
	}

	return pixel.NewSprite(prt, pixel.R(0, 0, prt.Bounds().W(), prt.Bounds().H()))
}

func Draw(t pixel.Target, m pixel.Matrix) {
	// draw coins and keys
	if _, ok := list[loot.COIN]; ok {
		piclist[loot.COIN].Draw(t, m)

		vec := pixel.Vec{20, -16}
		txt := text.New(pixel.ZV, atlas)
		txt.LineHeight = atlas.LineHeight() * 0.5
		txt.Color = colornames.Whitesmoke
		fmt.Fprintf(txt, "x")
		txt2 := text.New(pixel.Vec{10, 0}, atlasbig)
		txt2.LineHeight = atlasbig.LineHeight() * 1.3
		txt2.Color = colornames.Whitesmoke
		fmt.Fprintf(txt2, "%d", list[loot.COIN])

		txt.Draw(t, m.Moved(vec))
		txt2.Draw(t, m.Moved(vec))
	}

	if _, ok := list[loot.KEY]; ok {
		v := pixel.Vec{0, -40}
		for i := 0; i < list[loot.KEY]; i++ {
			piclist[loot.KEY].Draw(t, m.Moved(v))
			v.X += 15
		}
	}

}
