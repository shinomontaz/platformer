package controller

const (
	E_MOVE = iota
	E_CTRL
	E_SHIFT
	E_ESCAPE
)

type Subscriber interface {
	GetId() int
	Notify(e int)
}
