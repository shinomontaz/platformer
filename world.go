package main

import (
	"github.com/faiface/pixel"
	"github.com/faiface/pixel/imdraw"
)

type world struct {
	platforms []*platform
}

func NewWorld(w, h float64) *world {
	wrld := world{
		platforms: make([]*platform, 0),
	}
	return &wrld
}

func (w *world) Draw(t pixel.Target) {
	imd := imdraw.New(nil)
	for _, p := range w.platforms {
		p.draw(imd)
	}
	imd.Draw(t)
}
