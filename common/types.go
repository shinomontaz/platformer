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
	KeyAction(keys pixelgl.Button)
}

type EventSubscriber interface {
	EventAction(e int)
}

const (
	GROUND int = iota
	BARRIER
	WATER
	OBJECT
)

type Animater interface {
	GetSprite(name string, idx int) (pixel.Picture, pixel.Rect)
	GetLen(name string) int
}

type SimpleObjecter interface {
	GetRect() pixel.Rect
	GetPos() pixel.Vec
	GetId() int
	Update(dt float64, visiblePhys []O)
	Draw(t pixel.Target)
	IsGround() bool
}

type SpecialObjecter interface {
	SimpleObjecter
	UpdateSpecial(dt float64, visibleSpec []O)
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
	GetMass() float64
}

type O interface {
	Rect() pixel.Rect
	Id() int
}

type Objecter = O
