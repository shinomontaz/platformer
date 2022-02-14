package magic

import (
	"platformer/config"

	"platformer/animation"

	"github.com/faiface/pixel"
)

var w Worlder

func SetWorld(w Worlder) {
	w = w
}

func createspell(s string, source, target pixel.Vec) Speller {
	switch s {
	case "deathstrike":
		return &Deathstrike{
			anim:   animation.Get("deathstrike"),
			w:      w,
			power:  1,
			source: source,
			target: target,
			ttl:    3.0,
		}
	default:
		return nil
	}
}

type Magic struct {
	list   map[string]struct{}
	spells []Speller
}

func (m *Magic) AddSpell(s string) {
	// check against a list of spells
	// config.Spells map[string]Spellprofile

	if _, ok := config.Spells[s]; ok {
		m.list[s] = struct{}{}
	}
}

func (m *Magic) GetSpells() {
}

func (m *Magic) Cast(s string, source, target pixel.Vec) {
	if _, ok := m.list[s]; !ok {
		return
	}
	m.spells = append(m.spells, createspell(s, source, target))
}

func (m *Magic) Update(dt float64) {
	i := 0
	for _, s := range m.spells {
		s.Update(dt)
		if !s.IsFinished() {
			m.spells[i] = s
			i++
		}
	}

	m.spells = m.spells[:i]
}

func (m *Magic) Draw(t pixel.Target) {
	for _, s := range m.spells {
		s.Draw(t)
	}
}
