package background

import (
	"image/png"
	"math"
	"os"

	"github.com/faiface/pixel"
)

type Back struct {
	width    float64
	height   float64
	p        pixel.Picture
	viewport pixel.Rect
	vector1  pixel.Vec
	vector2  pixel.Vec
	part1    *pixel.Sprite
	part2    *pixel.Sprite
	pos      pixel.Vec
	steps    int
}

func New(start pixel.Vec, viewport pixel.Rect, path string) *Back {
	width := viewport.W()
	height := viewport.H()
	x, y := viewport.Min.X, viewport.Min.Y
	b := Back{
		width:    width,
		height:   height,
		viewport: viewport, //pixel.R(0, 0, width, height),

		vector1: pixel.V(x+width/2+1, y+height/2),
		vector2: pixel.V(x+3*width/2, y+height/2),

		pos:   start,
		steps: 0,
	}
	bg, err := loadPicture(path)
	if err != nil {
		panic(err)
	}

	b.p = bg

	b.part1 = pixel.NewSprite(b.p, pixel.R(0, 0, 600, 450))
	b.part2 = pixel.NewSprite(b.p, pixel.R(600, 0, 1200, 450))

	return &b
}

func (b *Back) Draw(t pixel.Target, pos pixel.Vec, cam pixel.Vec) {
	x, _ := pos.XY()
	cam.Y = 0
	steps := int(math.Abs(x / b.width))

	if steps != b.steps {
		b.vector1, b.vector2 = b.vector2, b.vector1

		b.steps = steps
	}

	x = math.Mod(x, b.width)

	b.pos = pixel.V(x, 0)

	b.part1.Draw(t, pixel.IM.Moved(b.vector1.Sub(b.pos).Sub(cam)))
	b.part2.Draw(t, pixel.IM.Moved(b.vector2.Sub(b.pos).Sub(cam)))
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
