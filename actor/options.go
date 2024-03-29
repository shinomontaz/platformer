package actor

import (
	"platformer/actor/statemachine"
	"platformer/common"
	"platformer/config"
	"platformer/particles"

	"github.com/shinomontaz/pixel"
)

type Option func(*Actor)

func WithRun(speed float64) Option {
	return func(a *Actor) {
		a.runspeed = speed
	}
}

func WithWalk(speed float64) Option {
	return func(a *Actor) {
		a.walkspeed = speed
	}
}

func WithJump(force float64) Option {
	return func(a *Actor) {
		a.jumpforce = force
	}
}

func WithStatemachine(sm *statemachine.Machine) Option {
	return func(a *Actor) {
		a.sm = sm
	}
}

func WithAnimDir(animdir float64) Option {
	return func(a *Actor) {
		a.animdir = animdir
	}
}

func WithHP(hp int) Option {
	return func(a *Actor) {
		a.hp = hp
	}
}

func WithScore(v int) Option {
	return func(a *Actor) {
		a.score = v
	}
}

func WithWaterResistance(r bool) Option {
	return func(a *Actor) {
		a.iswaterresistant = r
	}
}

func WithStrength(strength int) Option {
	return func(a *Actor) {
		a.strength = strength
	}
}

func WithMass(mass float64) Option {
	return func(a *Actor) {
		a.mass = mass
	}
}

func WithAttackrange(attackrange float64) Option {
	return func(a *Actor) {
		a.attackrange = attackrange
	}
}

func WithPortrait(path string) Option {
	return func(a *Actor) {
		if path == "" {
			return
		}
		prt, err := loader.LoadPicture(path)
		if err != nil {
			panic(err)
		}

		a.portrait = pixel.NewSprite(prt, pixel.R(0, 0, prt.Bounds().W(), prt.Bounds().H()))
	}
}

func WithSound(seffects []config.Soundeffect) Option {
	return func(a *Actor) {
		for _, se := range seffects {
			a.sounds[se.Type] = soundeffect{
				List: se.List,
			}
		}
	}
}

func WithSkills(skills []config.Skill) Option {
	return func(a *Actor) {
		for _, s := range skills {
			a.skills = append(a.skills, NewSkill(s))
		}
	}
}

func WithPhrases(key string) Option {
	return func(a *Actor) {
		a.phrasesClass = key
	}
}

func WithBody(body string) Option {
	return func(a *Actor) {
		if body == "human" {
			a.onhit = func(pos, vel pixel.Vec) {
				n := 10 * common.GetRandInt()
				for i := 0; i < n; i++ {
					particles.AddBlood(pos, pixel.V(vel.X+common.GetRandFloat()*1000-500, common.GetRandFloat()*10000))
				}
			}
		}

	}
}
