package animation

import (
	"image/png"
	"os"

	"github.com/faiface/pixel"
)

var anims map[string]*Anims

func init() {
	anims = make(map[string]*Anims)
}

type Anim struct {
	sheet  pixel.Picture
	frames []pixel.Rect
	sprite *pixel.Sprite
}

type Anims struct {
	items  map[string]*Anim
	rect   pixel.Rect
	sprite *pixel.Sprite
	m      float64
}

func Load(cfg AnimatingConfig) {
	animRect := pixel.R(0, 0, cfg.W(), cfg.H())
	a := New(animRect, cfg.M())
	names, files, frames := cfg.Get()
	for i := 0; i < len(names); i++ {
		a.SetAnim(names[i], files[i], frames[i])
	}
	anims[cfg.N()] = a
}

func Get(name string) *Anims {
	return anims[name]
}

func New(rect pixel.Rect, margin float64) *Anims {
	return &Anims{
		rect:   rect,
		items:  make(map[string]*Anim),
		sprite: pixel.NewSprite(nil, pixel.Rect{}),
		m:      margin,
	}
}

func (a *Anims) GetSprite(name string, num int) (pixel.Picture, pixel.Rect) {
	_, ok := a.items[name]
	if !ok {
		name = "idle" // fallback animation
		num = 1
	}

	idx := num % len(a.items[name].frames)
	return a.items[name].sheet, a.items[name].frames[idx]
}

func (a *Anims) GetLen(name string) int {
	return a.items[name].GetLen()
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
		x += a.m
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

func (a *Anim) GetLen() int {
	return len(a.frames)
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
