package actor

type Animater interface {
}

type Physicer interface {
}

type ActorStater interface {
	Start()
	Update(dt float64)
	Notify(e int)
}

type Commander interface {
}
