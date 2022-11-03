package gamestate

import (
	"platformer/actor"
	"platformer/controller"
	"platformer/events"
	"platformer/world"

	"github.com/shinomontaz/pixel/pixelgl"

	"github.com/shinomontaz/pixel"
)

var fragSource2 = `
#version 330 core

in vec2 vTexCoords;
out vec4 fragColor;

uniform sampler2D uTexture;
uniform vec4 uTexBounds;

// custom uniforms
uniform float uTime;

void main() {
    vec2 t = vTexCoords / uTexBounds.zw;
	vec3 influence = texture(uTexture, t).rgb;
	float uSpeed = 5.0;

    if (influence.r + influence.g + influence.b > 0.3) {
		t.y += cos(t.x * 40.0 + (uTime * uSpeed))*0.005;
		t.x += cos(t.y * 40.0 + (uTime * uSpeed))*0.01;
	}

    vec3 col = texture(uTexture, t).rgb;
	fragColor = vec4(col * vec3(0.6, 0.6, 1.2),1.0);
}
`

type Dialog struct {
	Common
	win   *pixelgl.Window
	ctrl  *controller.Controller
	cnv   *pixelgl.Canvas
	uTime float32
}

func NewDialog(game Gamer, w *world.World, hero *actor.Actor, win *pixelgl.Window) *Dialog {
	d := &Dialog{
		Common: Common{
			game:    game,
			id:      DIALOG,
			w:       w,
			hero:    hero,
			lastPos: pixel.ZV,
		},
		win:  win,
		ctrl: controller.New(win),
	}

	d.cnv = pixelgl.NewCanvas(win.Bounds())
	d.cnv.SetSmooth(true)
	d.cnv.SetUniform("uTime", &d.uTime)
	d.cnv.SetFragmentShader(fragSource2)

	d.ctrl.AddListener(d) // to listen ESCAPE  keyborad event

	return d
}

func (d *Dialog) Update(dt float64) {
	if dt > 0 {
		d.ctrl.Update()
		d.w.Update(d.currBounds, dt)
		d.uTime = float32(dt)
	}
}

func (d *Dialog) Draw(win *pixelgl.Window) {
	camPos := d.lastPos.Add(pixel.V(0, 150))

	d.w.Draw(d.cnv, d.lastPos, camPos, d.cnv.Bounds().Center())
	d.cnv.Draw(win, pixel.IM.Moved(win.Bounds().Center()))
	// draw dialog on center of window, no shaders upplied to it, only to background
}

func (d *Dialog) GetId() int {
	return d.id
}

func (d *Dialog) Start() {
	d.currBounds = d.w.GetViewport()
	d.lastPos = d.hero.GetPos()
}

func (d *Dialog) Listen(e int, v pixel.Vec) {
	switch e {
	case events.ESCAPE: // from controller
		d.game.SetState(NORMAL)
	}
}
