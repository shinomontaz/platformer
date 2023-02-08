package factories

import (
	"platformer/animation"
	"platformer/config"
	"platformer/loot"

	"github.com/shinomontaz/pixel"
)

func NewLoot(prof config.Profile, w Worlder) *loot.Loot {
	lootRect := pixel.R(0, 0, prof.Width, prof.Height)
	return loot.New(w, animation.Get(prof.Type), lootRect,
		loot.WithAnimDir(prof.Dir),
		loot.WithSound(config.Sounds[prof.Type].List),
	)
}
