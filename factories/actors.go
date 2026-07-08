package factories

import (
	"platformer/actor"
	"platformer/animation"
	"platformer/config"

	"github.com/shinomontaz/pixel"
)

func NewActor(prof config.Profile, w actor.Worlder) *actor.Actor {
	st := Machine(prof.Type)
	playerRect := pixel.R(0, 0, prof.Width, prof.Height)
	return actor.New(w, animation.Get(prof.Type), playerRect,
		actor.WithRun(prof.Run),
		actor.WithWalk(prof.Walk),
		actor.WithHP(prof.Hp),
		actor.WithStrength(prof.Strength),
		actor.WithMass(prof.Mass),
		actor.WithJump(prof.Jump),
		actor.WithStatemachine(st),
		actor.WithAnimDir(prof.Dir),
		actor.WithPortrait(prof.Portrait),
		actor.WithSound(config.Sounds[prof.Type].List),
		actor.WithSkills(prof.Skills),
		actor.WithPhrases(prof.Phrases),
		actor.WithBody(prof.Body),
		actor.WithScore(prof.Hp), // for now!
		actor.WithWaterResistance(prof.WaterResistance),
		actor.WithHurtReaction(prof.HurtReaction),
	)
}
