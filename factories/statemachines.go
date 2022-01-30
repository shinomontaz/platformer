package factories

import (
	"platformer/actor/state"
	"platformer/actor/statemachine"
	"platformer/events"
)

func NewPlayer() statemachine.Machine {
	m := statemachine.New()

	m.Set(state.MELEE, statemachine.Transition{})
	m.Set(state.MELEEMOVE, statemachine.Transition{})
	m.Set(state.STAND, statemachine.Transition{
		List: map[int]int{
			events.CTRL: state.MELEE,
		},
	})
	m.Set(state.IDLE, statemachine.Transition{
		List: map[int]int{
			events.CTRL: state.MELEE,
		},
	})
	m.Set(state.WALK, statemachine.Transition{
		List: map[int]int{
			events.CTRL: state.MELEEMOVE,
		},
	})
	// m.Set(state.RUN, statemachine.Transition{
	// 	List: map[int]int{
	// 		events.CTRL: state.MELEEMOVE,
	// 	},
	// })

	return m
}
