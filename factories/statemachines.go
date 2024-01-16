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
	case "goblin_mage":
		return newMage()
	case "goblin_range":
		fallthrough
	case "goblin_range2":
		return newRanged()
	default:
		return newEnemy()
	}
	return nil
}

func newPlayer() *statemachine.Machine {
	m := statemachine.New()
	m.Set(state.DEADSUNK, statemachine.Transition{})
	m.Set(state.MELEE, statemachine.Transition{})
	m.Set(state.MELEEMOVE, statemachine.Transition{})
	m.Set(state.FISHING, statemachine.Transition{})
	m.Set(state.ROLL, statemachine.Transition{})

	m.Set(state.STAND, statemachine.Transition{
		List: map[int]int{
			events.MELEE:     state.MELEE,
			events.MELEEMOVE: state.MELEEMOVE,
			events.INTERACT:  state.INTERACT,
		},
	})
	m.Set(state.IDLE, statemachine.Transition{
		List: map[int]int{
			events.MELEE:     state.MELEE,
			events.MELEEMOVE: state.MELEEMOVE,
			events.INTERACT:  state.INTERACT,
		},
	})
	m.Set(state.WALK, statemachine.Transition{
		List: map[int]int{
			events.ROLL:      state.ROLL,
			events.MELEEMOVE: state.MELEEMOVE,
			events.INTERACT:  state.INTERACT,
		},
	})

	return &m
}

func newEnemy() *statemachine.Machine {
	m := statemachine.New()

	m.Set(state.MELEE, statemachine.Transition{})
	m.Set(state.MELEEMOVE, statemachine.Transition{})
	m.Set(state.SWIM, statemachine.Transition{})
	m.Set(state.STAND, statemachine.Transition{
		List: map[int]int{
			events.MELEE:     state.MELEE,
			events.MELEEMOVE: state.MELEEMOVE,
		},
	})
	m.Set(state.IDLE, statemachine.Transition{
		List: map[int]int{
			events.MELEE:     state.MELEE,
			events.MELEEMOVE: state.MELEEMOVE,
		},
	})
	m.Set(state.WALK, statemachine.Transition{
		List: map[int]int{
			events.MELEEMOVE: state.MELEEMOVE,
		},
	})

	return &m
}

func newDeceased() *statemachine.Machine {
	m := statemachine.New()

	m.Set(state.CAST, statemachine.Transition{})
	m.Set(state.MELEE, statemachine.Transition{})
	m.Set(state.RANGED, statemachine.Transition{})
	m.Set(state.STAND, statemachine.Transition{
		List: map[int]int{
			events.CAST:   state.CAST,
			events.RANGED: state.RANGED,
			events.MELEE:  state.MELEE,
		},
	})
	m.Set(state.IDLE, statemachine.Transition{
		List: map[int]int{
			events.CAST:   state.CAST,
			events.RANGED: state.RANGED,
			events.MELEE:  state.MELEE,
		},
	})
	m.Set(state.WALK, statemachine.Transition{
		List: map[int]int{
			events.CAST:  state.CAST,
			events.MELEE: state.MELEE,
		},
	})

	return &m
}

func newMage() *statemachine.Machine {
	m := statemachine.New()

	m.Set(state.CAST, statemachine.Transition{})
	m.Set(state.STAND, statemachine.Transition{
		List: map[int]int{
			events.CAST: state.CAST,
		},
	})
	m.Set(state.IDLE, statemachine.Transition{
		List: map[int]int{
			events.CAST: state.CAST,
		},
	})

	return &m
}

func newRanged() *statemachine.Machine {
	m := statemachine.New()
	m.Set(state.RANGED, statemachine.Transition{})
	m.Set(state.STAND, statemachine.Transition{
		List: map[int]int{
			events.RANGED: state.RANGED,
		},
	})
	m.Set(state.IDLE, statemachine.Transition{
		List: map[int]int{
			events.RANGED: state.RANGED,
		},
	})

	return &m
}
