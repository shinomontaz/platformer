package gamestate

import (
	"fmt"
	"platformer/ai"
	"platformer/background"
	"platformer/common"
	"platformer/config"
	"platformer/controller"
	"platformer/creatures"
	"platformer/events"
	"platformer/factories"
	"platformer/inventory"
	"platformer/score"
	"platformer/world"

	"github.com/shinomontaz/pixel"
	"github.com/shinomontaz/pixel/pixelgl"
	"github.com/shinomontaz/pixel/text"
	"golang.org/x/image/colornames"
)

type Victory struct {
	Common
	win     *pixelgl.Window
	ctrl    *controller.Controller
	loader  *common.Loader
	w       *world.World
	initPos pixel.Vec
	//	pic     *pixel.Sprite
	txt    *text.Text
	txt2   *text.Text
	ailist []*ai.Ai
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
		ailist: make([]*ai.Ai, 0),
	}

	v.ctrl.AddKeyListener(v) // to listen ESCAPE  keyborad event

	return v
}

func (v *Victory) Update(dt float64) {
	if dt > 0 {
		v.ctrl.Update()
		v.w.Update(v.win.Bounds(), v.initPos, dt)
		for _, a := range v.ailist {
			a.Update(dt)
		}
	}
}

func (v *Victory) Draw(win *pixelgl.Window) {
	camPos := v.initPos.Add(pixel.V(0, 150))
	cntr := win.Bounds().Center()

	v.w.Draw(win, v.initPos, camPos, cntr)

	bounds := win.Bounds()
	//	picbounds := v.pic.Frame()
	//	v.pic.Draw(win, pixel.IM.Moved(bounds.Center()).ScaledXY(bounds.Center(), pixel.V(bounds.W()/picbounds.W(), bounds.H()/picbounds.H())))
	v.txt.Draw(win, pixel.IM.Moved(bounds.Center()).Moved(pixel.Vec{-50, 200}))
	v.txt2.Draw(win, pixel.IM.Moved(bounds.Center()).Moved(pixel.Vec{-50, 100}))
}

func (v *Victory) GetId() int {
	return v.id
}

func (v *Victory) Start() {
	w, err := world.New("victory.tmx", v.win.Bounds(), world.WithLoader(v.loader))
	if err != nil {
		panic(err)
	}
	v.w = w
	v.initPos = w.GetViewport().Center()

	crtrs := creatures.New()
	list := v.w.GetMetas()
	for _, o := range list {
		if o.Class == "npc" {
			npc := factories.NewActor(config.Profiles[o.Name], w)
			npc.Move(pixel.V(o.X, o.Y))

			ai_type := o.Properties.GetString("ai")
			if ai_type != "" {
				v.ailist = append(v.ailist, factories.NewAi(ai_type, npc, w))
			}

			dir := o.Properties.GetFloat("dir")
			if dir != 0 {
				npc.SetDir(dir)
			}

			crtrs.AddNpc(npc)
		}
	}

	w.SetCreatures(crtrs)

	b := background.New(v.initPos, v.win.Bounds().Moved(pixel.V(0, 150)), v.loader, "victory.png")
	v.w.SetBackground(b)

	// get image
	// pic, err := v.loader.LoadPicture("victory.png")
	// if err != nil {
	// 	panic(err)
	// }

	// v.pic = pixel.NewSprite(pic, pixel.R(0, 0, pic.Bounds().W(), pic.Bounds().H()))

	fnt := common.GetFont("regular32")
	atlas := text.NewAtlas(fnt, text.ASCII)
	v.txt = text.New(pixel.ZV, atlas)
	v.txt.Color = colornames.Cornsilk
	fmt.Fprintln(v.txt, "Victory!")

	scoretxt := fmt.Sprintf("Killed monsters score: %d", score.Count)
	scoretxt = fmt.Sprintf("%s\nCoins collected: %d", scoretxt, inventory.GetCoins()*10)

	fnt = common.GetFont("regular24")
	atlas = text.NewAtlas(fnt, text.ASCII)
	v.txt2 = text.New(pixel.ZV, atlas)
	v.txt2.Color = colornames.Cadetblue
	fmt.Fprintln(v.txt2, scoretxt)

	fmt.Println("victory screen prepared")
}

func (v *Victory) KeyAction(key pixelgl.Button) {
	switch key {
	case pixelgl.KeyEscape: // from controller
		v.game.Notify(events.STAGEVENT_QUIT)
	}
}
