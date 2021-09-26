package main

import (
	"fmt"
	"image/png"
	"math"
	"math/rand"
	"os"

	"github.com/faiface/pixel"
)

const (
	WALKING = iota
	RUNNING
	JUMPING
	STANDING
	IDLE
	FIRING
	DYING
	DEAD
)

const (
	NOACTION = iota
	STRIKE
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

	sheet      pixel.Picture
	frame      pixel.Rect
	sprite     *pixel.Sprite
	anims      map[string]*Anim
	dir        float64
	attackCode string
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

	h.attackCode = "attack1"

	return nil
}

func (h *Hero) getPos() pixel.Vec {
	return h.pos
}

func (h *Hero) Update(dt float64, cmd int) {
	// detect state
	h.counter += dt

	var newState int
	switch {
	case !h.phys.ground:
		newState = JUMPING
	case cmd == STRIKE:
		newState = FIRING
	case h.phys.vel.Len() == 0:
		newState = STANDING
	case h.phys.vel.Len() == h.phys.walkSpeed:
		newState = WALKING
	case h.phys.vel.Len() == h.phys.runSpeed:
		newState = RUNNING
	}

	if h.state == IDLE { // make idle animation
		newState = IDLE
	}
	if h.counter > 5.0 && h.state == STANDING && (newState == h.state || newState == STANDING) { // make idle animation
		newState = IDLE
	}

	// make state transision
	if h.state != newState {
		h.state = newState
		h.counter = 0
	}

	switch h.state {
	case IDLE:
		i := int(math.Floor(h.counter / 0.1)) // MAGIC CONST!

		if i > len(h.anims["idle"].frames) {
			h.state = STANDING
		}

		h.frame = h.anims["idle"].frames[i%len(h.anims["idle"].frames)]
		h.sheet = h.anims["idle"].sheet
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
	case FIRING:
		i := int(math.Floor(h.counter / 0.1)) // h.counter stands for frame rate change in animation
		if i == 0 || i%len(h.anims[h.attackCode].frames) == 0 {
			h.attackCode = h.selectAttack(h.phys.vel)
		}
		h.frame = h.anims[h.attackCode].frames[i%len(h.anims[h.attackCode].frames)]
		h.sheet = h.anims[h.attackCode].sheet
	}

	if h.phys.vel.X != 0 {
		if h.phys.vel.X > 0 {
			h.dir = +1
		} else {
			h.dir = -1
		}
	}

	//	h.pos = h.phys.rect.Center()
	h.pos = h.phys.rect.Min
	//	h.rect = h.rect.Moved(pixel.Vec{h.phys.rect.Center().X - h.phys.rect.W()/2, h.phys.rect.Center().Y - h.phys.rect.H()/2})

	h.rect = pixel.R(h.phys.rect.Center().X-h.phys.rect.W()/2, h.phys.rect.Center().Y-h.phys.rect.H()/2, h.phys.rect.Center().X-h.phys.rect.W()/2+h.rect.W(), h.phys.rect.Center().Y-h.phys.rect.H()/2+h.rect.H())
}

func (h *Hero) selectAttack(move pixel.Vec) string {
	if move.Len() == 0 {
		return fmt.Sprintf("attack%d", rand.Intn(2)+1)
	}
	return "attack3"
}

func (h *Hero) draw(t pixel.Target) {
	if h.sprite == nil {
		h.sprite = pixel.NewSprite(nil, pixel.Rect{})
	}
	// draw the correct frame with the correct position and direction
	h.sprite.Set(h.sheet, h.frame)
	h.sprite.Draw(t, pixel.IM.
		ScaledXY(pixel.ZV, pixel.V(
			h.rect.W()/h.sprite.Frame().W(),
			h.rect.H()/h.sprite.Frame().H(),
		)).
		ScaledXY(pixel.ZV, pixel.V(h.dir, 1)).
		Moved(h.rect.Center()),
	)
}

func (h *Hero) Hit(pos, vec pixel.Vec, power int) {
	// make him suffer

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
