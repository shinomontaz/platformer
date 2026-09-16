package factories

import (
	"platformer/actor"
	"platformer/animation"
	"platformer/config"

	"github.com/shinomontaz/pixel"
)

var bboxCache = make(map[string]pixel.Rect)

func correctBbox(prof config.Profile) pixel.Rect {
	bbox, ok := bboxCache[prof.Type]
	if ok {
		return bbox
	}

	sheet, frame := animation.Get(prof.Type).GetSprite("stand", 0)
	w, h := int(frame.W()), int(frame.H())
	minX, minY := w, h
	maxX, maxY := 0, 0

	pd := pixel.PictureDataFromPicture(sheet)
	for y := range h {
		for x := range w {
			sx := int(frame.Min.X) + x
			sy := int(frame.Min.Y) + y
			pidx := sy*pd.Stride + sx
			if pd.Pix[pidx].A > 0 {
				if x < minX {
					minX = x
				}
				if x > maxX {
					maxX = x
				}
				if y < minY {
					minY = y
				}
				if y > maxY {
					maxY = y
				}
			}
		}
	}

	return pixel.R(float64(minX), float64(minY), float64(maxX+1), float64(maxY+1))
}

func NewActor(prof config.Profile, w actor.Worlder) *actor.Actor {
	st := Machine(prof.Type)
	bbox := correctBbox(prof)
	playerRect := pixel.R(0, 0, prof.Width, prof.Height)

	return actor.New(w, animation.Get(prof.Type), playerRect,
		actor.WithBBox(bbox),
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
