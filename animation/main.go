package animation

import (
	"image/png"
	"os"

	"github.com/faiface/pixel"
)

type Anim struct {
	sheet  pixel.Picture
	frames []pixel.Rect
	sprite *pixel.Sprite
}

type Anims struct {
	items  map[string]*Anim
	rect   pixel.Rect
	sprite *pixel.Sprite
}

func New(rect pixel.Rect) *Anims {
	return &Anims{
		rect:   rect,
		items:  make(map[string]*Anim),
		sprite: pixel.NewSprite(nil, pixel.Rect{}),
	}
}

// func (a *Anims) GetAnims() map[string]*Anim {
// 	return a.items
// }

func (a *Anims) GetSprite(name string, idx int) *pixel.Sprite {
	//	return a.items[name]
	a.sprite.Set(a.items[name].sheet, a.items[name].frames[idx])

	return a.sprite
}

func (a *Anims) SetAnim(name, file string, frames []int) error {
	spritesheet, err := loadPicture(file)
	if err != nil {
		return err
	}

	frs := make([]pixel.Rect, 0, frames[0])

	frameWidth := a.rect.W()
	for x := 0.0; x+frameWidth <= spritesheet.Bounds().Max.X; x += frameWidth {
		frs = append(frs, pixel.R(
			x,
			0,
			x+frameWidth,
			spritesheet.Bounds().H(),
		))
	}

	a.items[name] = &Anim{
		sheet:  spritesheet,
		frames: frs[frames[1]:frames[2]],
		sprite: pixel.NewSprite(nil, pixel.Rect{}),
	}

	return nil
}

func (a *Anim) GetFrames() []pixel.Rect {
	return a.frames
}

func (a *Anim) GetSprite(idx int) *pixel.Sprite {
	a.sprite.Set(a.sheet, a.frames[idx])
	return a.sprite
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
