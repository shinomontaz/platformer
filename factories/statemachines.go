package factories

import (
	"platformer/actor/state"
	"platformer/actor/statemachine"
	"platformer/events"
)

func Machine(name string) *statemachine.Machine {
	switch name {
	case "player":
		return newPlayer()
	case "deceased":
		return newDeceased()
	default:
		return newEnemy()
	}
	return nil
}

func newPlayer() *statemachine.Machine {
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

	return &m
}

func newEnemy() *statemachine.Machine {
	m := statemachine.New()

	m.Set(state.MELEE, statemachine.Transition{})
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
			events.CTRL: state.MELEE,
		},
	})

	return &m
}

func newDeceased() *statemachine.Machine {
	m := statemachine.New()

	m.Set(state.CAST, statemachine.Transition{})
	m.Set(state.MELEE, statemachine.Transition{})
	m.Set(state.STAND, statemachine.Transition{
		List: map[int]int{
			events.CAST: state.CAST,
			events.CTRL: state.MELEE,
		},
	})
	m.Set(state.IDLE, statemachine.Transition{
		List: map[int]int{
			events.CAST: state.CAST,
			events.CTRL: state.MELEE,
		},
	})
	m.Set(state.WALK, statemachine.Transition{
		List: map[int]int{
			events.CAST: state.CAST,
			events.CTRL: state.MELEE,
		},
	})

	return &m
}
