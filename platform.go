package main

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

type platform struct {
	rect  pixel.Rect
	color color.Color
}

func NewPlatform(rect pixel.Rect) *platform {
	p := platform{rect: rect}
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

func (p *platform) Rect() *pixel.Rect {
	return &p.rect
}

func (p *platform) Draw(imd *imdraw.IMDraw) {
	imd.Color = p.color
	imd.Push(p.rect.Min, p.rect.Max)
	imd.Rectangle(0)
}
