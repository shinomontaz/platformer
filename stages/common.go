package stages

import (
	"platformer/events"

	"github.com/shinomontaz/pixel/pixelgl"
	"golang.org/x/image/colornames"
)

type Common struct {
	id       int
	done     chan struct{}
	isReady  bool
	isActive bool
	inform   Inform
	j        Job
	eventMap map[int]int
}

func (s *Common) GetID() int {
	return s.id
}

func (s *Common) Init() {
	s.isReady = true
}

func (s *Common) SetUp(opts ...StageOpt) {
	for _, opt := range opts {
		opt(s)
	}
}

func (s *Common) Start() {
	if !s.isReady {
		return
	}
	if s.j != nil {
		go func() {
			s.j()
			s.done <- struct{}{}
		}()
	}

	s.isActive = true
}

func (s *Common) Stop() {
	if !s.isReady {
		return
	}
	s.isActive = false
}

func (s *Common) SetJob(j Job) {
	s.j = j
}

func (s *Common) Notify(event int) {
	s.inform(event)
}

func (s *Common) Run(win *pixelgl.Window, dt float64) {
	select {
	case <-s.done:
		s.Notify(events.STAGEVENT_DONE)
	default:
		win.Clear(colornames.Black)
	}
}

func (s *Common) SetNext(event, id int) {
	s.eventMap[event] = id
}

func (s *Common) GetNext(event int) (int, bool) {
	next, ok := s.eventMap[event]
	return next, ok
}
