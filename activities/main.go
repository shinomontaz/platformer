package activities

import (
	"github.com/shinomontaz/pixel"
)

var listeners []Actor = make([]Actor, 0)

func Init(en []Actor) {
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
