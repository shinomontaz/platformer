package stages

import (
	"fmt"
	"platformer/common"
	"platformer/events"

	"github.com/shinomontaz/pixel"

	"github.com/shinomontaz/pixel/pixelgl"
	"github.com/shinomontaz/pixel/text"
	"golang.org/x/image/colornames"
)

type Loading struct {
	Common
	assetloader *common.Loader
	txt         *text.Text
}

func NewLoading(f Inform, l *common.Loader) *Loading {
	return &Loading{
		Common: Common{
			id:       LOADING,
			done:     make(chan struct{}),
			inform:   f,
			eventMap: make(map[int]int),
		},
		assetloader: l,
	}
}

func (l *Loading) Init() {
	fnt := common.GetFont("regular")
	atlas := text.NewAtlas(fnt, text.ASCII)

	l.txt = text.New(pixel.ZV, atlas)
	l.txt.Color = colornames.Red
	fmt.Fprintln(l.txt, "Loading...")

	l.isReady = true
}

func (l *Loading) Run(win *pixelgl.Window, dt float64) {
	select {
	case <-l.done:
		l.Notify(events.STAGEVENT_DONE)
	default:
		win.Clear(colornames.Black)
		l.txt.Draw(win, pixel.IM.Moved(win.Bounds().Center()))
	}
}
