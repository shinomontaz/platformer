package loot

import (
	"fmt"
	"platformer/animation"
	"platformer/common"
	"platformer/config"

	"github.com/shinomontaz/pixel"
)

var list []common.SimpleObjecter
var profiles map[string]config.Profile
var grav float64

func Init(w Worlder, loots map[string]config.Profile) {
	profiles = loots
	grav = w.GetGravity()
	list = make([]common.SimpleObjecter, 0)
}

func AddCoin(pos, vel pixel.Vec) {
	var prof config.Profile
	var ok bool
	if prof, ok = profiles["coin"]; !ok {
		fmt.Println("coin profile in loot not found!")
		panic("!")
	}
	lootRect := pixel.R(0, 0, prof.Width, prof.Height)
	l := New(animation.Get(prof.Type), lootRect,
		WithAnimDir(prof.Dir),
		WithSound(config.Sounds[prof.Type].List),
		WithGravity(grav),
		WithVelocity(vel),
	)

	l.Move(pos)
	add(l)
}

func AddKey(pos, vel pixel.Vec) {
	var prof config.Profile
	var ok bool
	if prof, ok = profiles["key"]; !ok {
		fmt.Println("key profile in loot not found!")
		panic("!")
	}
	lootRect := pixel.R(0, 0, prof.Width, prof.Height)
	l := New(animation.Get(prof.Type), lootRect,
		WithAnimDir(prof.Dir),
		WithSound(config.Sounds[prof.Type].List),
		WithGravity(grav),
		WithVelocity(vel),
	)

	l.Move(pos)
	add(l)
}

func add(l common.SimpleObjecter) {
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

func Draw(t pixel.Target) {
	for _, l := range list {
		l.Draw(t)
	}
}
