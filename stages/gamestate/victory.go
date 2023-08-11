package gamestate

import (
	"fmt"
	"platformer/common"
	"platformer/controller"
	"platformer/events"

	"github.com/shinomontaz/pixel"
	"github.com/shinomontaz/pixel/pixelgl"
	"github.com/shinomontaz/pixel/text"
	"golang.org/x/image/colornames"
)

type Victory struct {
	Common
	win    *pixelgl.Window
	ctrl   *controller.Controller
	loader *common.Loader
	pic    *pixel.Sprite
	txt    *text.Text
}

func NewVictory(game Gamer, win *pixelgl.Window, l *common.Loader) *Victory {
	v := &Victory{
		Common: Common{
			game: game,
			id:   VICTORY,
		},
		win:    win,
		ctrl:   controller.New(win, true),
		loader: l,
	}

	v.ctrl.AddKeyListener(v) // to listen ESCAPE  keyborad event

	return v
}

func (v *Victory) Update(dt float64) {
	if dt > 0 {
		v.ctrl.Update()
	}
}

func (v *Victory) Draw(win *pixelgl.Window) {
	//	camPos := v.lastPos.Add(pixel.V(0, 150))
	bounds := win.Bounds()
	picbounds := v.pic.Frame()
	v.pic.Draw(win, pixel.IM.Moved(bounds.Center()).ScaledXY(bounds.Center(), pixel.V(bounds.W()/picbounds.W(), bounds.H()/picbounds.H())))
	v.txt.Draw(win, pixel.IM.Moved(bounds.Center()).Moved(pixel.Vec{-50, 200}))
}

func (v *Victory) GetId() int {
	return v.id
}

func (v *Victory) Start() {
	fmt.Println("victory screen")
	// geet image
	pic, err := v.loader.LoadPicture("victory.png")
	if err != nil {
		panic(err)
	}

	v.pic = pixel.NewSprite(pic, pixel.R(0, 0, pic.Bounds().W(), pic.Bounds().H()))

	fnt := common.GetFont("regular32")
	atlas := text.NewAtlas(fnt, text.ASCII)
	v.txt = text.New(pixel.ZV, atlas)
	v.txt.Color = colornames.Cornsilk
	fmt.Fprintln(v.txt, "Victory!")
}

func (v *Victory) KeyAction(key pixelgl.Button) {
	switch key {
	case pixelgl.KeyEscape: // from controller
		v.game.Notify(events.STAGEVENT_QUIT)
	}
}
