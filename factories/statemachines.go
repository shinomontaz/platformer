package factories

import (
	"platformer/actor/state"
	"platformer/actor/statemachine"
	"platformer/events"
)

func NewPlayer() statemachine.Machine {
	m := statemachine.New()

	m.Set(state.ATTACK, statemachine.Transition{})
	m.Set(state.STAND, statemachine.Transition{
		List: map[int]int{
			events.CTRL: state.ATTACK,
		},
	})
	m.Set(state.WALK, statemachine.Transition{
		List: map[int]int{
			events.CTRL: state.ATTACK,
		},
	})
	m.Set(state.RUN, statemachine.Transition{
		List: map[int]int{
			events.CTRL: state.ATTACK,
		},
	})

	return m
}
