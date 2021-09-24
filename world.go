package main

import (
	"platformer/config"

	"github.com/faiface/pixel"
	"github.com/faiface/pixel/imdraw"
)

type world struct {
	platforms []*platform
	enemies   []*Enemy
	visible   []Objecter
}

type Enemy struct {
	p phys
	a Hero
}

func NewWorld(w, h float64) *world {
	wrld := world{
		platforms: make([]*platform, 0),
		enemies:   make([]*Enemy, 0),
	}
	return &wrld
}

func NewEnemy(cfg config.Enemy, wcfg config.World) *Enemy {
	e := Enemy{}

	e.p = phys{
		rect:      pixel.R(cfg.Coords[0], cfg.Coords[1], cfg.Coords[0]+cfg.Width, cfg.Coords[1]+cfg.Height),
		runSpeed:  cfg.Run,
		walkSpeed: cfg.Walk,
		jumpSpeed: wcfg.Gravity * 50,
		ground:    true,
		gravity:   wcfg.Gravity,
	}

	e.a = Hero{
		phys:  &e.p,
		rect:  e.p.rect,
		anims: make(map[string]*Anim, 0),
		pos:   pixel.V(0.0, 0.0),
		dir:   1.0,
	}

	for _, anim := range *cfg.Anims {
		e.a.SetAnim(anim.Name, anim.File, anim.Frames)
	}

	return &e
}

func (w *world) Update(rect pixel.Rect) {
	// update viewport and detect visible objects to draw only them
	w.visible = make([]Objecter, 0, 10)
	for _, p := range w.platforms {
		if rect.Intersects(p.rect) {
			w.visible = append(w.visible, p)
		}
	}

	// for _, e := range w.enemies {
	// 	e.p.Update()
	// }

	// for _, e := range w.enemies {
	// 	if rect.Intersects(e.p.rect) {
	// 		w.visible = append(w.visible, e)
	// 	}
	// }
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

	for _, e := range w.enemies {
		e.a.draw(t)
	}

	imd.Draw(t)
}
