package main

import (
	"github.com/faiface/pixel"
	"github.com/faiface/pixel/pixelgl"
)

const (
	STATE_FREE = iota
	STATE_ATTACK
	STATE_DEAD
)

type ControllerStater interface {
	Process(dt float64, win *pixelgl.Window)
	Start()
}

type FreeState struct {
	id int
	pc *PlayerController
}

func (s *FreeState) Start() {
	s.pc.SetCmd(NOACTION)
}

func (s *FreeState) Process(dt float64, win *pixelgl.Window) {
	vec := pixel.ZV

	if win.Pressed(pixelgl.KeyLeftShift) {
		if win.Pressed(pixelgl.KeyLeft) {
			vec.X--
		} else if win.Pressed(pixelgl.KeyRight) {
			vec.X++
		}
	} else if !win.Pressed(pixelgl.KeyLeftShift) { // running
		if win.Pressed(pixelgl.KeyLeft) {
			vec.X -= 2.0
		} else if win.Pressed(pixelgl.KeyRight) {
			vec.X += 2.0
		}
	}

	if win.Pressed(pixelgl.KeyUp) {
		vec.Y++
	}

	s.pc.SetVec(vec)
	s.pc.SetCmd(NOACTION)

	if win.Pressed(pixelgl.KeyLeftControl) {
		s.pc.SetState(STATE_ATTACK)
	}
}

type AttackState struct {
	id        int
	pc        *PlayerController
	time      float64
	timelimit float64
}

func (s *AttackState) Start() {
	s.time = 0.0
	s.pc.SetCmd(STRIKE)
}

func (s *AttackState) Process(dt float64, win *pixelgl.Window) {
	if s.time > s.timelimit {
		s.pc.SetState(STATE_FREE)
		return
	}

	s.time += dt
	s.pc.SetVec(pixel.ZV)
	s.pc.SetCmd(STRIKE)
}

type DeadState struct {
	id int
	pc *PlayerController
}

func (s *DeadState) Start() {
	s.pc.SetCmd(NOACTION)
}

func (s *DeadState) Process(dt float64, win *pixelgl.Window) {
	s.pc.SetVec(pixel.ZV)
	s.pc.SetCmd(NOACTION)
}
