package factories

import (
	"platformer/actor"
	"platformer/animation"
	"platformer/common"
	"platformer/config"

	"github.com/shinomontaz/pixel"
)

type Worlder interface {
	GetGravity() float64
	GetQt() *common.Quadtree // temporary solution, we will check collision in the world ?
	AddStrike(owner *actor.Actor, r pixel.Rect, power int, speed pixel.Vec)
	AddSpell(owner *actor.Actor, t pixel.Vec, spell string)
	AddAlert(pos pixel.Vec, force float64)
	IsSee(from, to pixel.Vec) bool
	GetHero() *actor.Actor
}

func NewActor(prof config.Profile, w Worlder) *actor.Actor {
	st := Machine(prof.Type)
	playerRect := pixel.R(0, 0, prof.Width, prof.Height)
	return actor.New(w, animation.Get(prof.Type), playerRect,
		actor.WithRun(prof.Run),
		actor.WithWalk(prof.Walk),
		actor.WithHP(prof.Hp),
		actor.WithStrength(prof.Strength),
		actor.WithJump(prof.Jump),
		actor.WithStatemachine(st),
		actor.WithAnimDir(prof.Dir),
		actor.WithPortrait(prof.Portrait),
		actor.WithSound(config.Sounds[prof.Type].List),
		actor.WithSkills(prof.Skills),
	)
}
