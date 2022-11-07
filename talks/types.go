package talks

import (
	"github.com/shinomontaz/pixel"
)

type Worlder interface {
	Alert(rect pixel.Rect)
}
