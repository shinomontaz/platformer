package stages

import (
	"github.com/shinomontaz/pixel"
	"github.com/shinomontaz/pixel/pixelgl"
)

type Job func()
type Inform func(e int)

type StageOpt func(s Stager)

type Stager interface {
	GetID() int
	Run(win *pixelgl.Window, dt float64)
	Start()
	Stop()
	SetJob(j Job)
	Init()
	GetNext(event int) (int, bool)
	SetNext(event, id int)
	SetUp(opts ...StageOpt)
}

type Gamestater interface {
	GetId() int
	Start()
	Listen(e int, v pixel.Vec)
	Update(dt float64)
	Draw(win *pixelgl.Window)
}
