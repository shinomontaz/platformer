package common

import (
	"github.com/shinomontaz/pixel"
	"github.com/shinomontaz/pixel/pixelgl"
)

type Subscriber interface {
	//	GetId() int
	Listen(e int, v pixel.Vec)
}

type KeySubscriber interface {
	KeyEvent(keys pixelgl.Button)
}

const (
	GROUND = iota
	BARRIER
	WATER
	OBJECT
)

type Animater interface {
	GetSprite(name string, idx int) (pixel.Picture, pixel.Rect)
	GetGroupSprite(group, name string, idx int) (pixel.Picture, pixel.Rect)
	GetGroupLen(name string) int
	GetLen(name string) int
}

type SimpleObjecter interface {
	GetRect() pixel.Rect
	GetPos() pixel.Vec
	GetId() int
	Update(dt float64, visiblePhys []Objecter)
	Draw(t pixel.Target)
	IsGround() bool
}

type SpecialObjecter interface {
	SimpleObjecter
	UpdateSpecial(dt float64, visibleSpec []Objecter)
}

type Interactor interface {
	Interact()
	OnInteract()
}

type Actorer interface {
	SpecialObjecter
	Interactor
	Hit(vec pixel.Vec, power int)
	GetHp() int
}
