package statemachine

type Machine struct {
	tr    map[int]Transition
	Start int
}

func New() Machine {
	return Machine{
		tr: map[int]Transition{},
	}
}

func (m *Machine) Set(state int, tr Transition) {
	m.tr[state] = tr
}

func (m *Machine) GetStates() []int {
	keys := make([]int, len(m.tr))

	i := 0
	for k := range m.tr {
		keys[i] = k
		i++
	}

	return keys
}

func (m *Machine) GetTransition(state int) Transition {
	if trans, ok := m.tr[state]; ok {
		return trans
	}
	return Transition{}
}

type Transition struct {
	List map[int]int
}
