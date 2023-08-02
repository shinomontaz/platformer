package menu

import (
	"fmt"
	"platformer/config"

	"github.com/shinomontaz/pixel/pixelgl"

	"github.com/shinomontaz/pixel"
	"github.com/shinomontaz/pixel/text"
)

var soundmenu *Menu

func NewSound(r pixel.Rect, atlas *text.Atlas, opts ...MenuOption) *Menu {
	if soundmenu != nil {
		soundmenu.rect = r
		for _, o := range opts {
			o(soundmenu)
		}
		return soundmenu
	}

	soundmenu = New(r)
	for _, o := range opts {
		o(soundmenu)
	}
	txt := text.New(pixel.V(0, 0), atlas)

	it := NewItem(fmt.Sprintf("%v: %-10v", "Master", config.Opts.Volumes["main"]), txt,
		WithHandle(func(b pixelgl.Button) {
			switch b {
			case pixelgl.KeyLeft:
				if config.Opts.Volumes["main"] > 0 {
					config.Opts.Volumes["main"] -= 10
				}
			case pixelgl.KeyRight:
				if config.Opts.Volumes["main"] < 100 {
					config.Opts.Volumes["main"] += 10
				}
			}
			soundmenu.UpdateSelectedItemText(fmt.Sprintf("%v: %-10v", "Master", config.Opts.Volumes["main"]))
		}),
	)
	soundmenu.AddItem(it)
	it.Select(true)

	txt = text.New(pixel.V(0, 0), atlas)
	it = NewItem(fmt.Sprintf("%v: %-10v", "Music", config.Opts.Volumes["music"]), txt,
		WithHandle(func(b pixelgl.Button) {
			switch b {
			case pixelgl.KeyLeft:
				if config.Opts.Volumes["music"] > 0 {
					config.Opts.Volumes["music"] -= 10
				}
			case pixelgl.KeyRight:
				if config.Opts.Volumes["music"] < 100 {
					config.Opts.Volumes["music"] += 10
				}
			}
			soundmenu.UpdateSelectedItemText(fmt.Sprintf("%v: %-10v", "Music", config.Opts.Volumes["music"]))
		}),
	)
	soundmenu.AddItem(it)

	txt = text.New(pixel.V(0, 0), atlas)
	it = NewItem(fmt.Sprintf("%v: %-10v", "Actions", config.Opts.Volumes["actions"]), txt,
		WithHandle(func(b pixelgl.Button) {
			switch b {
			case pixelgl.KeyLeft:
				if config.Opts.Volumes["actions"] > 0 {
					config.Opts.Volumes["actions"] -= 10
				}
			case pixelgl.KeyRight:
				if config.Opts.Volumes["actions"] < 100 {
					config.Opts.Volumes["actions"] += 10
				}
			}
			soundmenu.UpdateSelectedItemText(fmt.Sprintf("%v: %-10v", "Actions", config.Opts.Volumes["actions"]))
		}),
	)
	soundmenu.AddItem(it)

	txt = text.New(pixel.V(0, 0), atlas)
	it = NewItem("Quit", txt, WithAction(func() {
		if soundmenu.onquit != nil {
			soundmenu.onquit()
		}
	}))
	soundmenu.AddItem(it)
	return soundmenu
}
