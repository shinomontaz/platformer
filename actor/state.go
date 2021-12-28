package actor

const (
	STATE_FREE = iota
	STATE_ATTACK
	STATE_HIT
	STATE_DEAD
)

type CommonState struct {
	id int
	a  *Actor
	//	animations map[string]Animation
	anims Animater
}

func (s *CommonState) GetId() int {
	return s.id
}
