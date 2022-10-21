package stages

import (
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
