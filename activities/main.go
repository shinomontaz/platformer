package activities

import (
	"platformer/common"

	"github.com/shinomontaz/pixel"
)

var listeners []common.Actorer = make([]common.Actorer, 0)

func Init(en []common.Actorer) {
	listeners = en
}

func Alert(rect pixel.Rect) {
	// for _, en := range listeners {
	// 	if rect.Contains(en.GetPos()) {
	// 		a := ai.GetByObj(en)
	// 		if a != nil {
	// 			a.Listen(events.ALERT, rect.Center())
	// 		}
	// 	}
	// }
}
