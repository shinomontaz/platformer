package projectiles

import (
	"platformer/activities"
	"platformer/animation"
	"platformer/common"
	"platformer/common/rand"

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

type SkillData struct {
	Name     string
	Strength int
	Pos      pixel.Vec
	Ttl      float64
	Trace    float64
}

// AddProjectile adds a projectile of kind t, attached to owner and
// in position of owner.GetRect().Center()
func AddProjectile(sd SkillData, dir float64, owner common.Actorer) {
	var f pixel.Vec
	g := grav
	switch sd.Name {
	case "sling":
		f = pixel.V(160*dir, 80+rand.Float64()*120)
	case "crossbow":
		f = pixel.V(150*dir, 8+rand.Float64()*12)
		g = grav / 9
	case "shamanbolt":
		f = pixel.V(75*dir, 0)
		g = grav / 1000
	case "ancientbolt":
		f = pixel.V(90*dir, 0)
		g = grav / 1000
	case "deathbolt":
		f = pixel.V(130*dir, 0)
		g = grav / 1000
	}

	pos := owner.GetRect().Center().Add(pixel.V(sd.Pos.X*dir, sd.Pos.Y))
	p := projectile{
		pos:      pos,
		force:    f,
		mass:     0.001,
		rigidity: 0,
		friction: 1,
		dir:      dir,
		anim:     animation.Get(sd.Name),
		strength: float64(sd.Strength),
		owner:    owner,
		ttl:      sd.Ttl,
		trace:    sd.Trace,
	}
	p.Init(g)

	hb := activities.AddStrike(
		owner,
		p.rect,
		sd.Strength,
		pixel.ZV,
		activities.WithTTL(-1),
	)

	p.hb = hb
	projectiles = append(projectiles, p)
}

func Update(dt float64, objs, spec []common.O) {
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
