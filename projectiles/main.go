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

// AddProjectile adds a projectile of kind t, attached to owner and
// in position of owner.GetRect().Center()
func AddProjectile(t string, strength float64, dir float64, owner common.Actorer) {
	var f pixel.Vec
	g := grav
	switch t {
	case "sling":
		f = pixel.V(10000*dir, 5000+common.GetRandFloat()*7500)
	case "crossbow":
		f = pixel.V(9000*dir, 500+common.GetRandFloat()*750)
		g = grav / 9
	case "deathbolt":
		f = pixel.V(8000*dir, 0)
		g = grav / 1000
	}

	pos := owner.GetRect().Center()
	p := projectile{
		pos:      pos,
		force:    f,
		mass:     0.001,
		rigidity: 0,
		friction: 1,
		dir:      dir,
		anim:     animation.Get(t),
		strength: strength,
		owner:    owner,
		size:     8,
	}
	p.Init(g)
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
		if pr.hb.IsHitted() {
			pr.active = false
		}
		if pr.active {
			projectiles[i] = pr
			projectiles[i].update(dt, objs, spec)
			i++
		}
	}
	for j := i; j < len(projectiles); j++ {
		projectiles[j].hb.SetTtl(0)
	}
	projectiles = projectiles[:i]
}

func Draw(t pixel.Target) {
	for i := range projectiles {
		projectiles[i].draw(t)
	}
}
