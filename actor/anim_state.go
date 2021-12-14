package actor

import "platformer/controller"

const (
	ASTATE_STAND = iota
	ASTATE_IDLE
	ASTATE_WALK
	ASTATE_RUN
	ASTATE_JUMP
	ASTATE_ATTACK
)

type AnimStateStand struct {
	id int
	s  ActorStater
}

func (s *AnimStateStand) GetId() int {
	return s.id
}

func (s *AnimStateStand) Start() {
}

func (s *AnimStateStand) Update(dt float64) {
	// detect animation and make it idle in case of long inactivity
	// make movements
}

func (s *AnimStateStand) Notify(e int) {
}

type AnimStateIdle struct {
	id int
	s  ActorStater
}

func (s *AnimStateIdle) GetId() int {
	return s.id
}

func (s *AnimStateIdle) Start() {
}

func (s *AnimStateIdle) Update(dt float64) {
	// detect animation and make it idle in case of long inactivity
	// make movements
}

func (s *AnimStateIdle) Notify(e int) {
}

type AnimStateWalk struct {
	id    int
	s     ActorStater
	shift bool
}

func (s *AnimStateWalk) GetId() int {
	return s.id
}

func (s *AnimStateWalk) Start() {
}

func (s *AnimStateWalk) Update(dt float64) {
	// detect animation and make it idle in case of long inactivity
	// make movements
	s.shift = false
}

func (s *AnimStateWalk) Notify(e int) {
	if e == controller.E_SHIFT {
		s.shift = true
	}
}

type AnimStateJump struct {
	id int
	s  ActorStater
}

func (s *AnimStateJump) GetId() int {
	return s.id
}

func (s *AnimStateJump) Start() {
}

func (s *AnimStateJump) Update(dt float64) {
	// detect animation and make it idle in case of long inactivity
	// make movements
}

func (s *AnimStateJump) Notify(e int) {
}
