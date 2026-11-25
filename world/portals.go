package world

import (
	"log"
	"platformer/common"

	"github.com/shinomontaz/pixel"
)

type Portal struct {
	Name        string // o.Name из TMX — ключ
	Rect        pixel.Rect
	Destination string          // свойство "destination", "" = вылазка без возврата
	allow       bool            // свойство "allow", default false
	anim        common.Animater // animation.Get("portal") в init-ветке
	counter     float64
}

func (w *World) AllowPortal(name string) {
	if p, ok := w.portals[name]; ok {
		p.allow = true
	}
}

func (w *World) updatePortals(dt float64) {
	hero := w.crtrs.GetHero()
	if hero == nil {
		return
	}
	hr := hero.GetRect()

	for _, p := range w.portals {
		p.counter += dt
		if !p.allow || p.Destination == "" || !hr.Intersects(p.Rect) {
			continue
		}
		dst, ok := w.portals[p.Destination]
		if !ok {
			log.Printf("[portal] %q: destination %q not found", p.Name, p.Destination)
			continue
		}

		hero.Teleport(dst.Rect.Center())

		break
	}
}

func (w *World) drawPortals(t pixel.Target) {
	for _, p := range w.portals {
		if p.anim == nil || !p.allow {
			continue // анимации "portal" нет в anims.json — просто не рисуем
		}

		pic, rect := p.anim.GetSprite("idle", int(p.counter/0.15))
		sprite := pixel.NewSprite(pic, rect)
		sprite.Draw(t, pixel.IM.
			ScaledXY(pixel.ZV, pixel.V(
				p.Rect.W()/sprite.Frame().W(),
				p.Rect.H()/sprite.Frame().H(),
			)).
			Moved(p.Rect.Center()),
		)
	}
}
