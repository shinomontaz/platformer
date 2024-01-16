package creatures

import (
	"github.com/shinomontaz/pixel"

	"platformer/activities"
	"platformer/actor"
	"platformer/common"
)

type List struct {
	npcs    []common.Actorer
	enemies []common.Actorer
	hero    common.Actorer
}

func New() *List {
	return &List{
		enemies: make([]common.Actorer, 0),
		npcs:    make([]common.Actorer, 0),
	}
}

// func reset_lists() {
// 	enemies = make([]common.Actorer, 0)
// 	npcs = make([]common.Actorer, 0)
// }

func (l *List) AddEnemy(enemy common.Actorer) {
	l.enemies = append(l.enemies, enemy)
}

func (l *List) AddNpc(npc *actor.Actor) {
	l.npcs = append(l.npcs, npc)
}

func (l *List) DeleteNpc(npc *actor.Actor) {
	i := 0
	for _, n := range l.npcs {
		if n == npc {
			continue
		}
		l.npcs[i] = n
		i++
	}
	l.npcs = l.npcs[:i]
}

func (l *List) GetHero() common.Actorer {
	return l.hero
}

func (l *List) SetHero(h common.Actorer) {
	l.hero = h
}

func (l *List) Update(dt float64, visiblePhys, visibleSpec []common.Objecter) {

	if l.hero != nil {
		l.hero.Update(dt, visiblePhys)
		l.hero.UpdateSpecial(dt, visibleSpec)
	}

	for _, en := range l.enemies {
		en.Update(dt, visiblePhys)
	}
	for _, npc := range l.npcs {
		npc.Update(dt, visiblePhys)
		npc.UpdateSpecial(dt, visibleSpec)
	}
	activities.UpdateStrikes(dt, l.enemies, l.hero)
	activities.UpdateInteractions(dt, l.npcs)
}

func (l *List) Draw(t pixel.Target) {
	//	activities.DrawStrikes(t)
	for _, e := range l.enemies {
		e.Draw(t)
	}
	for _, n := range l.npcs {
		n.Draw(t)
	}
	if l.hero != nil {
		l.hero.Draw(t)
	}
}
