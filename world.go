package main

import (
	"github.com/faiface/pixel"
	"github.com/faiface/pixel/imdraw"
)

type world struct {
	platforms []*platform
	visible   []Objecter
}

func NewWorld(w, h float64) *world {
	wrld := world{
		platforms: make([]*platform, 0),
	}
	return &wrld
}

func (w *world) Update(rect pixel.Rect) {
	// update viewport and detect visible objects to draw only them
	w.visible = make([]Objecter, 0, 10)
	for _, p := range w.platforms {
		if rect.Intersects(p.rect) {
			w.visible = append(w.visible, p)
		}
	}
}

func (w *world) Objects() []Objecter {
	// return visible platforms only
	return w.visible
}

func (w *world) Draw(t pixel.Target) {
	imd := imdraw.New(nil)
	// for _, p := range w.platforms {
	// 	p.draw(imd)
	// }
	for _, p := range w.visible {
		p.Draw(imd)
	}

	imd.Draw(t)
}
