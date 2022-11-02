package gamestate

import (
	"platformer/actor"
	"platformer/ui"
	"platformer/world"

	"github.com/shinomontaz/pixel"
)

type Common struct {
	id            int
	game          Gamer
	currBounds    pixel.Rect
	u             *ui.Ui
	w             *world.World
	hero          *actor.Actor
	initialCenter pixel.Vec
	lastPos       pixel.Vec

	camPos   pixel.Vec
	deltaVec pixel.Vec
}
