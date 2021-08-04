package main

import (
	"image/png"
	"math"
	"os"

	"github.com/faiface/pixel"
)

const (
	WALKING = iota
	RUNNING
	JUMPING
	STANDING
	FIRING
	DYING
	DEAD
)

type Anim struct {
	sheet  pixel.Picture
	frames []pixel.Rect
}
type Hero struct {
	phys *phys
	rect pixel.Rect
	pos  pixel.Vec

	state   int
	counter float64

	sheet  pixel.Picture
	frame  pixel.Rect
	sprite *pixel.Sprite
	anims  map[string]*Anim
	dir    float64
}

func (h *Hero) SetAnim(name, file string, frames []int) error {
	spritesheet, err := loadPicture(file)
	if err != nil {
		return err
	}

	frs := make([]pixel.Rect, 0, frames[0])

	frameWidth := h.rect.W()
	for x := 0.0; x+frameWidth <= spritesheet.Bounds().Max.X; x += frameWidth {
		frs = append(frs, pixel.R(
			x,
			0,
			x+frameWidth,
			spritesheet.Bounds().H(),
		))
	}

	h.anims[name] = &Anim{
		sheet:  spritesheet,
		frames: frs[frames[1]:frames[2]],
	}

	return nil
}

func (h *Hero) getPos() pixel.Vec {
	return h.pos
}

func (h *Hero) Update(dt float64) {
	// detect state
	h.counter += dt

	var newState int
	switch {
	case !h.phys.ground:
		newState = JUMPING
	case h.phys.vel.Len() == 0:
		newState = STANDING
	case h.phys.vel.Len() == h.phys.walkSpeed:
		newState = WALKING
	case h.phys.vel.Len() == h.phys.runSpeed:
		newState = RUNNING
	}

	// make state transision

	if h.state != newState {
		h.state = newState
		h.counter = 0
	}

	switch h.state {
	case STANDING:
		h.frame = h.anims["stand"].frames[0]
		h.sheet = h.anims["stand"].sheet
	case WALKING:
		i := int(math.Floor(h.counter / 0.1)) // MAGIC CONST!
		h.frame = h.anims["walk"].frames[i%len(h.anims["walk"].frames)]
		h.sheet = h.anims["walk"].sheet
	case RUNNING:
		i := int(math.Floor(h.counter / 0.1)) // MAGIC CONST!
		h.frame = h.anims["run"].frames[i%len(h.anims["run"].frames)]
		h.sheet = h.anims["run"].sheet
	case JUMPING:
		speed := h.phys.vel.Y
		i := int((-speed/h.phys.jumpSpeed + 1) / 2 * float64(len(h.anims["jump"].frames)))
		if i < 0 {
			i = 0
		}
		if i >= len(h.anims["jump"].frames) {
			i = len(h.anims["jump"].frames) - 1
		}
		h.frame = h.anims["jump"].frames[i]
		h.sheet = h.anims["jump"].sheet
	}

	if h.phys.vel.X != 0 {
		if h.phys.vel.X > 0 {
			h.dir = +1
		} else {
			h.dir = -1
		}
	}

	h.pos = h.pos.Add(h.phys.vel.Scaled(dt))
}

func (h *Hero) draw(t pixel.Target) {
	if h.sprite == nil {
		h.sprite = pixel.NewSprite(nil, pixel.Rect{})
	}
	// draw the correct frame with the correct position and direction
	h.sprite.Set(h.sheet, h.frame)
	h.sprite.Draw(t, pixel.IM.
		ScaledXY(pixel.ZV, pixel.V(
			h.phys.rect.W()/h.sprite.Frame().W(),
			h.phys.rect.H()/h.sprite.Frame().H(),
		)).
		ScaledXY(pixel.ZV, pixel.V(h.dir, 1)).
		Moved(h.phys.rect.Center()),
	)
}

func loadPicture(path string) (pixel.Picture, error) {
	file, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer file.Close()
	img, err := png.Decode(file)
	if err != nil {
		return nil, err
	}

	return pixel.PictureDataFromImage(img), nil
}
