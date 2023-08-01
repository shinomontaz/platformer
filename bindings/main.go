package bindings

import (
	"fmt"
	"log"
	"platformer/config"

	"github.com/shinomontaz/pixel/pixelgl"
)

type Bindings struct {
	list    map[int]pixelgl.Button
	actions map[pixelgl.Button]int
}

const (
	CTRL = 1 << iota
	ENTER
	SHIFT
	LEFT
	RIGHT
	UP
	DOWN
)

var Default, Active *Bindings

var KeyActions = []string{
	"ctrl",
	"enter",
	"shift",
	"left",
	"right",
	"up",
	"down",
}

var ActionId = []int{
	CTRL,
	ENTER,
	SHIFT,
	LEFT,
	RIGHT,
	UP,
	DOWN,
}

// var KeyAction = map[int]string{
// 	CTRL:  "ctrl",
// 	ENTER: "enter",
// 	SHIFT: "shift",
// 	LEFT:  "left",
// 	RIGHT: "right",
// 	UP:    "up",
// 	DOWN:  "down",
// }

var KeyAction = map[string]int{
	"ctrl":  CTRL,
	"enter": ENTER,
	"shift": SHIFT,
	"left":  LEFT,
	"right": RIGHT,
	"up":    UP,
	"down":  DOWN,
}

var KeyActionNames = map[int]string{
	CTRL:  "Strike",
	ENTER: "Interact",
	SHIFT: "Toggle run",
	LEFT:  "Left",
	RIGHT: "Right",
	UP:    "Jump",
	DOWN:  "Crunch",
}

func Init() { // setup default and active bindings
	Default = NewDefaultBindings()

	if len(config.Opts.Bindings) == 0 {
		Active = NewDefaultBindings()
	} else {
		fmt.Println("config.Opts.Bindings not empty ", config.Opts.Bindings)
		Active = &Bindings{
			list: map[int]pixelgl.Button{
				CTRL:  pixelgl.Button(config.Opts.Bindings["ctrl"]),
				ENTER: pixelgl.Button(config.Opts.Bindings["enter"]),
				SHIFT: pixelgl.Button(config.Opts.Bindings["shift"]),
				LEFT:  pixelgl.Button(config.Opts.Bindings["left"]),
				RIGHT: pixelgl.Button(config.Opts.Bindings["right"]),
				UP:    pixelgl.Button(config.Opts.Bindings["up"]),
				DOWN:  pixelgl.Button(config.Opts.Bindings["down"]),
			},
			actions: map[pixelgl.Button]int{
				pixelgl.Button(config.Opts.Bindings["ctrl"]):  CTRL,
				pixelgl.Button(config.Opts.Bindings["enter"]): ENTER,
				pixelgl.Button(config.Opts.Bindings["shift"]): SHIFT,
				pixelgl.Button(config.Opts.Bindings["left"]):  LEFT,
				pixelgl.Button(config.Opts.Bindings["right"]): RIGHT,
				pixelgl.Button(config.Opts.Bindings["up"]):    UP,
				pixelgl.Button(config.Opts.Bindings["down"]):  DOWN,
			},
		}
	}
}

func (b *Bindings) List() map[int]pixelgl.Button {
	return b.list
}

func (b *Bindings) GetBinding(event int) pixelgl.Button {
	if bind, ok := b.list[event]; ok {
		return bind
	}
	log.Fatalf("no binding for event %s", event)
	return 0
}

// returns ActionId(int) for provided pressed button
func (b *Bindings) GetAction(key pixelgl.Button) int {
	if action, ok := b.actions[key]; ok {
		return action
	}
	log.Fatalf("no action for key %s", key)
	return -1
}

func (b *Bindings) SetBind(key pixelgl.Button, event int) {
	b.list[event] = key
}

func (b *Bindings) Save() {
	config.Opts.Bindings = map[string]int{
		"ctrl":  int(b.list[CTRL]),
		"enter": int(b.list[ENTER]),
		"shift": int(b.list[SHIFT]),
		"left":  int(b.list[LEFT]),
		"right": int(b.list[RIGHT]),
		"up":    int(b.list[UP]),
		"down":  int(b.list[DOWN]),
	}
}

func NewDefaultBindings() *Bindings {
	return &Bindings{
		list: map[int]pixelgl.Button{
			CTRL:  pixelgl.KeyLeftControl,
			ENTER: pixelgl.KeyEnter,
			SHIFT: pixelgl.KeyLeftShift,
			LEFT:  pixelgl.KeyLeft,
			RIGHT: pixelgl.KeyRight,
			UP:    pixelgl.KeyUp,
			DOWN:  pixelgl.KeyDown,
		},
		actions: map[pixelgl.Button]int{
			pixelgl.KeyLeftControl: CTRL,
			pixelgl.KeyEnter:       ENTER,
			pixelgl.KeyLeftShift:   SHIFT,
			pixelgl.KeyLeft:        LEFT,
			pixelgl.KeyRight:       RIGHT,
			pixelgl.KeyUp:          UP,
			pixelgl.KeyDown:        DOWN,
		},
	}
}
