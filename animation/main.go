package animation

import (
	"image/png"
	"os"

	"github.com/faiface/pixel"
)

type Anim struct {
	sheet  pixel.Picture
	frames []pixel.Rect
}

type Anims struct {
	items map[string]*Anim
	rect  pixel.Rect
}

func New(rect pixel.Rect) *Anims {
	return &Anims{
		rect:  rect,
		items: make(map[string]*Anim),
	}
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
	}

	return nil
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
