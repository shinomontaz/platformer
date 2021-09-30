package main

import (
	"github.com/faiface/pixel"
	"github.com/faiface/pixel/pixelgl"
)

type PlayerController struct {
	vec    pixel.Vec
	cmd    int // command action
	ground bool
	state  ControllerStater
	states map[int]ControllerStater
}

func NewController() *PlayerController {
	ctrl := &PlayerController{}

	sFree := &FreeState{
		id: STATE_FREE,
		pc: ctrl,
	}
	sAttack := &AttackState{
		id:        STATE_ATTACK,
		pc:        ctrl,
		timelimit: 0.5,
	}
	sDead := &DeadState{
		id: STATE_DEAD,
		pc: ctrl,
	}

	ctrl.states = map[int]ControllerStater{STATE_FREE: sFree, STATE_ATTACK: sAttack, STATE_DEAD: sDead}
	ctrl.state = sFree

	return ctrl
}

func (pc *PlayerController) SetState(id int) {
	pc.state = pc.states[id]
	pc.state.Start()
}

func (pc *PlayerController) SetVec(vec pixel.Vec) {
	pc.vec = vec
}

func (pc *PlayerController) SetCmd(cmd int) {
	pc.cmd = cmd
}

func (pc *PlayerController) SetGround(g bool) {
	pc.ground = g
}

func (pc *PlayerController) Update(dt float64, win *pixelgl.Window) {
	pc.state.Process(dt, win)
}
