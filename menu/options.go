package menu

import (
	"github.com/shinomontaz/pixel"
	"github.com/shinomontaz/pixel/pixelgl"
)

type MenuOption func(*Menu)
type ItemOption func(*Item)

func WithQuit(quit func()) MenuOption {
	return func(m *Menu) {
		m.onquit = quit
	}
}

func WithLogo(pic *pixel.Sprite) MenuOption {
	return func(m *Menu) {
		m.logo = pic
	}
}

func WithTitle(title string) MenuOption {
	return func(m *Menu) {
		m.title = title
	}
}

func WithHandle(f func(pixelgl.Button)) ItemOption {
	return func(i *Item) {
		i.handle = f
	}
}

func WithAction(f func()) ItemOption {
	return func(i *Item) {
		i.action = f
	}
}
