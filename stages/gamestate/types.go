package gamestate

import "platformer/ai"

const (
	NORMAL = iota
	DEAD
	MENU
	DIALOG
	VICTORY
)

type Gamer interface {
	SetState(int)
	GetAiList() []*ai.Ai
	Notify(e int)
}
