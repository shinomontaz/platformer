package factories

import (
	"platformer/ai"
)

func NewAi(t string, obj ai.Manageder, w Worlder) {
	switch t {
	case "deceased":
		ai.NewMage(obj, w)
	default:
		ai.NewCommon(obj, w)
	}
}
