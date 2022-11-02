package gamestate

const (
	NORMAL = iota
	DEAD
	MENU
	DIALOG
)

type Gamer interface {
	SetState(int)
	Notify(e int)
}
