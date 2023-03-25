package projectiles

import (
	"platformer/activities"
	"platformer/animation"
	"platformer/common"

	"github.com/shinomontaz/pixel"
)

var (
	projectiles []projectile
	grav        float64
)

func Init(g float64) {
	projectiles = make([]projectile, 0)
	grav = g
}

func AddProjectile(t string, pos, f pixel.Vec, strength float64, dir float64, owner common.Actorer) {
	p := projectile{
		pos:      pos,
		force:    f,
		mass:     0.1,
		rigidity: 0,
		friction: 1,
		dir:      dir,
		anim:     animation.Get(t),
		strength: strength,
		owner:    owner,
		size:     8,
	}
	p.Init()
	hb := activities.AddStrike(
		owner,
		p.rect,
		int(strength),
		pixel.ZV,
		activities.WithTTL(-1),
	)
	p.hb = hb
	projectiles = append(projectiles, p)
}

func Update(dt float64, objs, spec []common.Objecter) {
	i := 0
	for _, pr := range projectiles {
		if pr.active {
			projectiles[i] = pr
			projectiles[i].update(dt, objs, spec)
			i++
		}
	}
	projectiles = projectiles[:i]
}

func Draw(t pixel.Target) {
	for i := range projectiles {
		projectiles[i].draw(t)
	}
}
