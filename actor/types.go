package actor

import "github.com/faiface/pixel"

type Animater interface {
	//	GetAnims() map[string]Animation
	GetSprite(name string, idx int) *pixel.Sprite
}

type Animation interface {
	GetSprite(idx int) *pixel.Sprite
}

type Physicer interface {
	GetVel() *pixel.Vec
	Update(dt float64, v pixel.Vec)
}

type ActorStater interface {
	GetId() int
	Start()
	Update(dt float64)
	Notify(e int, v *pixel.Vec)
	GetSprite() *pixel.Sprite
}

type AnimStater interface {
	GetId() int
	Start()
	Update(dt float64)
	Notify(e int)
}
