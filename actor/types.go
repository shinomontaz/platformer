package actor

import "github.com/faiface/pixel"

type Animater interface {
	GetAnims() map[string]Animation
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
	SetAnim(name string, anim Animation)
}

type AnimStater interface {
	GetId() int
	Start()
	Update(dt float64)
	Notify(e int)
}
