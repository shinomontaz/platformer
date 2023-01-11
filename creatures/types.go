package creatures

import (
	tmx "github.com/lafriks/go-tiled"
)

type Worlder interface {
	GetMetas() []*tmx.Object
}
