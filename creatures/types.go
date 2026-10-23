package creatures

import (
	tmx "github.com/lafriks/go-tiled"
	"github.com/shinomontaz/pixel"
)

type Worlder interface {
	GetMetas() []*tmx.Object
}

type object interface {
	Rect() pixel.Rect
	Id() int
}
