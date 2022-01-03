package world

import (
	"image/color"
	"math"
	"math/rand"
	"time"

	"github.com/faiface/pixel"
	"github.com/faiface/pixel/imdraw"
)

func init() {
	rand.Seed(time.Now().UnixNano())
}

type Platform struct {
	rect  pixel.Rect
	color color.Color
}

func NewPlatform(rect pixel.Rect) *Platform {
	p := Platform{rect: rect}
	p.color = random_color()
	return &p
}

func random_color() pixel.RGBA {
again:
	r := rand.Float64()
	g := rand.Float64()
	b := rand.Float64()
	len := math.Sqrt(r*r + g*g + b*b)
	if len == 0 {
		goto again
	}
	return pixel.RGB(r/len, g/len, b/len)
}

func (p *Platform) Rect() *pixel.Rect {
	return &p.rect
}

func (p *Platform) Pixels() []uint32 {
	return nil
}

func (p *Platform) Draw(imd *imdraw.IMDraw) {
	imd.Color = p.color
	imd.Push(p.rect.Min, p.rect.Max)
	imd.Rectangle(0)
}

func (p *Platform) Hit(pos, vec pixel.Vec, power int) {
	// do nothing
}
