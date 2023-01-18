package common

import "github.com/shinomontaz/pixel"

type Subscriber interface {
	//	GetId() int
	Listen(e int, v pixel.Vec)
}

const (
	GROUND = iota
	BARRIER
	WATER
)

type Actorer interface {
	GetRect() pixel.Rect
	GetPos() pixel.Vec
	GetId() int
	Hit(vec pixel.Vec, power int)
	Update(dt float64, visiblePhys []Objecter)
	UpdateSpecial(dt float64, visibleSpec []Objecter)
	Draw(t pixel.Target)
	GetHp() int
	IsGround() bool
}
