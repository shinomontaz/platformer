package particles

import (
	"platformer/common"

	"github.com/shinomontaz/pixel"
	"github.com/shinomontaz/pixel/pixelgl"
)

const DEFAULT_SIZE = 2

var (
	canvas    *pixelgl.Canvas
	particles []particle
	idx       int
	batch     *pixel.Batch
	colors    []uint8
	colormap  map[uint32]int
	grav      float64
)

func Init(maxnum int) {
	particles = make([]particle, 0, maxnum)
	for i := 0; i < maxnum; i++ {
		particles = append(particles, particle{active: false})
	}

	colormap = make(map[uint32]int)

	canvas = pixelgl.NewCanvas(pixel.R(0, 0, 1, 1))
	batch = pixel.NewBatch(&pixel.TrianglesData{}, canvas)
}

func addColorToBatch(color uint32) {
	pos := color | 0xFF
	if _, ok := colormap[pos]; !ok {
		r := color >> 24 & 0xFF
		g := color >> 16 & 0xFF
		b := color >> 8 & 0xFF
		colors = append(colors, uint8(r), uint8(g), uint8(b), 255.0)
		colormap[pos] = (len(colors) / 4) - 1

		canvas.SetBounds(pixel.R(0, 0, float64(len(colors)/4), 1)) // grow canvas
		canvas.SetPixels(colors)
	}
}

func SetGravity(g float64) {
	grav = g
}

func AddBlood(pos, f pixel.Vec) {
	r := 175 + common.GetRandInt()*5
	g := 10 + common.GetRandInt()*2
	b := 10 + common.GetRandInt()*2
	a := 255
	//	col := colornames.Red

	col := uint32(r&0xFF<<24 | g&0xFF<<16 | b&0xFF<<8 | a&0xFF)
	addColorToBatch(col)

	idx++
	if idx >= len(particles) {
		idx = 0
	}

	p := particle{
		active:   true,
		color:    col,
		pos:      pos,
		ttl:      2*common.GetRandFloat() + 0.5,
		force:    f,
		size:     DEFAULT_SIZE,
		mass:     0.1,
		rigidity: 0.1,
	}
	p.Init()
	particles[idx] = p
}

func Update(dt float64, objs []common.Objecter) {
	batch.Clear()
	sprite := pixel.NewSprite(canvas, pixel.R(0, 0, 1, 1))
	for i := range particles {
		if particles[i].active {
			particles[i].update(dt, objs)

			color := particles[i].color
			r := color >> 24 & 0xFF
			g := color >> 16 & 0xFF
			b := color >> 8 & 0xFF
			//			a := color & 0xFF // need for transparent particles

			pos := r&0xFF<<24 | g&0xFF<<16 | b&0xFF<<8 | 0xFF
			sprite.Set(canvas, pixel.R(float64(colormap[pos]), 0, float64(colormap[pos]+1), 1))
			sprite.Draw(batch, pixel.IM.Scaled(pixel.ZV, particles[i].size).Moved(particles[i].pos)) //
		}
	}
}

func Draw(t pixel.Target) {
	batch.Draw(t)
	// imd := imdraw.New(nil)
	// for i := range particles {
	// 	if particles[i].active {
	// 		particles[i].draw(imd)
	// 	}
	// }
	// imd.Draw(t)
}
