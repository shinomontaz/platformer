package gamestate

import (
	"image/color"
	"platformer/actor"
	"platformer/controller"
	"platformer/dialogs"
	"platformer/ui"
	"platformer/world"

	"github.com/shinomontaz/pixel/pixelgl"

	"github.com/shinomontaz/pixel"
)

var fragSource2 = `
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

type Dialog struct {
	Common
	win     *pixelgl.Window
	ctrl    *controller.Controller
	cnv     *pixelgl.Canvas
	uTime   float32
	currDlg *dialogs.Dialog
}

func NewDialog(game Gamer, u *ui.Ui, w *world.World, hero *actor.Actor, win *pixelgl.Window) *Dialog {
	d := &Dialog{
		Common: Common{
			game:    game,
			id:      DIALOG,
			u:       u,
			w:       w,
			hero:    hero,
			lastPos: pixel.ZV,
		},
		win:  win,
		ctrl: controller.New(win, true),
	}

	d.cnv = pixelgl.NewCanvas(win.Bounds())
	d.cnv.SetSmooth(true)
	d.cnv.SetUniform("uTime", &d.uTime)
	d.cnv.SetFragmentShader(fragSource2)

	d.ctrl.AddListener(d) // to listen ESCAPE  keyborad event

	return d
}

func (d *Dialog) Update(dt float64) {
	d.currDlg = dialogs.GetActive()
	if d.currDlg == nil {
		d.game.SetState(NORMAL)
		return
	}

	if dt > 0 {
		d.ctrl.Update()
		//		d.w.Update(d.currBounds, dt)
		d.uTime = float32(dt)
	}
}

func (d *Dialog) Draw(win *pixelgl.Window) {
	camPos := d.lastPos.Add(pixel.V(0, 150))

	d.cnv.Clear(color.RGBA{0, 0, 0, 1})

	d.w.Draw(d.cnv, d.lastPos, camPos, d.cnv.Bounds().Center())
	//	d.u.Draw(d.cnv)
	d.cnv.Draw(win, pixel.IM.Moved(win.Bounds().Center()))

	// draw dialog on center of window, no shaders upplied to it, only to background
	d.currDlg.Draw(win)
}

func (d *Dialog) GetId() int {
	return d.id
}

func (d *Dialog) Start() {
	d.lastPos = d.hero.GetPos()

	d.currDlg = dialogs.GetActive()
	d.currDlg.Start(d.win.Bounds())
	// hero => get active dialog
	// prepare options and dialog
}

func (d *Dialog) KeyEvent(key pixelgl.Button) {
	// if up or down - handle just here, otherwise make item handle it
	switch key {
	case pixelgl.KeyUp:
		d.currDlg.UpdateAnswer(-1)
	case pixelgl.KeyDown:
		d.currDlg.UpdateAnswer(+1)
	case pixelgl.KeyEnter:
		d.currDlg.Action()
	case pixelgl.KeyEscape:
		d.game.SetState(NORMAL)
	}
}

// func (d *Dialog) Listen(e int, v pixel.Vec) {
// 	// if up or down - handle just here, otherwise make item handle it
// 	if v.Y > 0 {
// 		d.currDlg.UpdateAnswer(-1)
// 	}
// 	if v.Y < 0 {
// 		d.currDlg.UpdateAnswer(+1)
// 	}
// 	if e == events.INTERACT {
// 		d.currDlg.Action()
// 	}
// 	if e == events.ESCAPE {
// 		d.game.SetState(NORMAL)
// 	}
// }
