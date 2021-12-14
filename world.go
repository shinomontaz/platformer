package main

import (
	"github.com/faiface/pixel"
	"github.com/faiface/pixel/imdraw"
)

type world struct {
	platforms []*platform
	// enemies   []*Enemy
	// currenm   []*Enemy
	visible []Objecter
}

// type Enemy struct {
// 	p  phys
// 	a  Hero
// 	ai Ai
// }

func NewWorld(w, h float64) *world {
	wrld := world{
		platforms: make([]*platform, 0),
		//		enemies:   make([]*Enemy, 0),
	}
	return &wrld
}

// func NewEnemy(cfg config.Enemy, wcfg config.World) *Enemy {
// 	e := Enemy{}

// 	e.p = NewPhys()
// 	e.p.rect = pixel.R(cfg.Coords[0], cfg.Coords[1], cfg.Coords[0]+cfg.Width/2, cfg.Coords[1]+cfg.Height*0.75)
// 	e.p.runSpeed = cfg.Run
// 	e.p.walkSpeed = cfg.Walk
// 	e.p.jumpSpeed = wcfg.Gravity * 50
// 	e.p.gravity = wcfg.Gravity

// 	e.a = Hero{
// 		phys:  &e.p,
// 		rect:  pixel.R(cfg.Coords[0], cfg.Coords[1], cfg.Coords[0]+cfg.Width, cfg.Coords[1]+cfg.Height),
// 		anims: make(map[string]*Anim, 0),
// 		pos:   pixel.V(0.0, 0.0),
// 		dir:   1.0,
// 	}

// 	for _, anim := range *cfg.Anims {
// 		e.a.SetAnim(anim.Name, anim.File, anim.Frames)
// 	}

// 	// add Ai

// 	e.ai = Ai{
// 		pers: &e.p,
// 	}

// 	return &e
// }

func (w *world) Update(rect pixel.Rect) {
	// update viewport and detect visible objects to draw only them
	w.visible = make([]Objecter, 0, 10)
	for _, p := range w.platforms {
		if rect.Intersects(p.rect) {
			w.visible = append(w.visible, p)
		}
	}

	// w.currenm = make([]*Enemy, 0)
	// for _, e := range w.enemies {
	// 	if rect.Intersects(e.p.rect) {
	// 		w.currenm = append(w.currenm, e)
	// 	}
	// }
}

func (w *world) GetObjects() []Objecter {
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

	// for _, e := range w.enemies {
	// 	e.a.draw(t)
	// }

	imd.Draw(t)
}
