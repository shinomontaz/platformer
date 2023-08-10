package gamestate

const (
	NORMAL = iota
	DEAD
	MENU
	DIALOG
	VICTORY
)

type Gamer interface {
	SetState(int)
	Notify(e int)
}
