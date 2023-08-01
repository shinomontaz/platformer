package gamestate

import (
	"image/color"
	"platformer/actor"
	"platformer/controller"
	"platformer/world"

	"github.com/shinomontaz/pixel/pixelgl"

	"github.com/shinomontaz/pixel"
)

var fragSource = `
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

type Dead struct {
	Common
	win   *pixelgl.Window
	ctrl  *controller.Controller
	cnv   *pixelgl.Canvas
	uTime float32
}

func NewDead(game Gamer, w *world.World, hero *actor.Actor, win *pixelgl.Window) *Dead {
	d := &Dead{
		Common: Common{
			game:    game,
			id:      DEAD,
			w:       w,
			hero:    hero,
			lastPos: pixel.ZV,
		},
		win:  win,
		ctrl: controller.New(win, false),
	}

	d.cnv = pixelgl.NewCanvas(win.Bounds())
	d.cnv.SetSmooth(true)
	d.cnv.SetUniform("uTime", &d.uTime)
	d.cnv.SetFragmentShader(fragSource)

	d.ctrl.AddListener(d) // to listen ESCAPE  keyborad event

	return d
}

func (d *Dead) Update(dt float64) {
	if dt > 0 {
		d.ctrl.Update()
		d.w.Update(d.currBounds, d.lastPos, dt)
		d.uTime += float32(dt)
	}
}

func (d *Dead) Draw(win *pixelgl.Window) {
	camPos := d.lastPos.Add(pixel.V(0, 150))

	d.cnv.Clear(color.RGBA{0, 0, 0, 1})

	d.w.Draw(d.cnv, d.lastPos, camPos, d.cnv.Bounds().Center())
	d.cnv.Draw(win, pixel.IM.Moved(win.Bounds().Center()))
}

func (d *Dead) GetId() int {
	return d.id
}

func (d *Dead) Start() {
	d.currBounds = d.w.GetViewport()
	d.lastPos = d.hero.GetPos()
}

func (d *Dead) KeyEvent(key pixelgl.Button) {
	switch key {
	case pixelgl.KeyEscape: // from controller
		d.game.SetState(MENU)
	}
}
