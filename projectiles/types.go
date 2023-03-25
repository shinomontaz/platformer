package projectiles

import "github.com/shinomontaz/pixel"

type Animater interface {
	GetSprite(name string, idx int) (pixel.Picture, pixel.Rect)
	GetLen(name string) int
}
