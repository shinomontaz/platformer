package loot

import (
	"platformer/common"

	"github.com/shinomontaz/pixel"
)

var list []common.SimpleObjecter

func Add(l common.SimpleObjecter) {
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
