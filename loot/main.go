package loot

import (
	"fmt"
	"platformer/animation"
	"platformer/common"
	"platformer/config"

	"github.com/shinomontaz/pixel"
)

var (
	list []*Loot

	profiles  map[string]config.Profile
	grav      float64
	collected []*Loot
)

func Init(w Worlder, loots map[string]config.Profile) {
	profiles = loots
	grav = w.GetGravity()
	list = make([]*Loot, 0)
	collected = make([]*Loot, 0)
}

func AddCoin(pos, vel pixel.Vec) {
	var prof config.Profile
	var ok bool
	if prof, ok = profiles["coin"]; !ok {
		fmt.Println("coin profile in loot not found!")
		panic("!")
	}
	lootRect := pixel.R(0, 0, prof.Width, prof.Height)
	l := New(COIN, animation.Get(prof.Type), lootRect,
		WithAnimDir(prof.Dir),
		WithSound(config.Sounds[prof.Type].List),
		WithGravity(grav),
		WithVelocity(vel),
		WithPortrait(prof.Portrait),
	)

	l.Move(pos)
	add(l)
	fmt.Println("add coin", pos)
}

func AddKey(pos, vel pixel.Vec) {
	var prof config.Profile
	var ok bool
	if prof, ok = profiles["key"]; !ok {
		fmt.Println("key profile in loot not found!")
		panic("!")
	}
	lootRect := pixel.R(0, 0, prof.Width, prof.Height)
	l := New(KEY, animation.Get(prof.Type), lootRect,
		WithAnimDir(prof.Dir),
		WithSound(config.Sounds[prof.Type].List),
		WithGravity(grav),
		WithVelocity(vel),
		WithPortrait(prof.Portrait),
	)

	l.Move(pos)
	add(l)
}

func add(l *Loot) {
	list = append(list, l)
}

func Remove(l common.SimpleObjecter) {
	i := 0
	for _, ll := range list {
		if ll.GetId() != l.GetId() {
			list[i] = ll
			i++
		}
	}

	list = list[:i]
}

func Update(dt float64, visiblePhys []common.Objecter) {
	for _, l := range list {
		l.Update(dt, visiblePhys)
	}
}

func Collect(hero common.Actorer) []*Loot {
	collected = collected[:0]
	if len(list) == 0 {
		return nil
	}
	i := 0
	for _, l := range list {
		if (l.vel.X == 0 && l.vel.Y == 0) && l.rect.Intersects(hero.GetRect()) {
			collected = append(collected, l)
			continue
		}
		i++
	}
	list = list[:i]
	return collected
}

func Draw(t pixel.Target) {
	for _, l := range list {
		l.Draw(t)
	}
}
