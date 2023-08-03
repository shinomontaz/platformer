package controller

import "platformer/bindings"

func WithBinding(b *bindings.Bindings) Option {
	return func(c *Controller) {
		for key, val := range b.List() {
			c.currBind[key] = val
		}
	}
}
