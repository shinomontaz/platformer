package particles

import (
	"platformer/common"

	"github.com/shinomontaz/pixel"
	"github.com/shinomontaz/pixel/imdraw"
	"golang.org/x/image/colornames"
)

var (
	//	canvas    *pixelgl.Canvas
	particles []particle
	idx       int
	// batch     *pixel.Batch
	// colors    []uint8
	// colormap  map[uint32]int
	grav float64
)

func Init(maxnum int) {
	particles = make([]particle, 0, maxnum)
	for i := 0; i < maxnum; i++ {
		particles = append(particles, particle{active: false})
	}

	// colormap = make(map[uint32]int)

	// canvas = pixelgl.NewCanvas(pixel.R(0, 0, 1, 1)) // Max seems to be 2^14 per row
	// batch = pixel.NewBatch(&pixel.TrianglesData{}, canvas)
}

// func addColorToBatch(color uint32) {
// 	pos := color | 0xFF
// 	if _, ok := colormap[pos]; !ok {
// 		r := color >> 24 & 0xFF
// 		g := color >> 16 & 0xFF
// 		b := color >> 8 & 0xFF
// 		colors = append(colors, uint8(r), uint8(g), uint8(b), 255.0)
// 		colormap[pos] = (len(colors) / 4) - 1

// 		canvas.SetBounds(pixel.R(0, 0, float64(len(colors)/4), 1)) // grow canvas
// 		canvas.SetPixels(colors)
// 	}
// }

func SetGravity(g float64) {
	grav = g
}

func AddBlood(pos, f pixel.Vec) {
	// r := 175 + common.GetRandInt()*5
	// g := 10 + common.GetRandInt()*2
	// b := 10 + common.GetRandInt()*2
	// a := 255
	col := colornames.Red

	idx++
	if idx >= len(particles) {
		idx = 0
	}

	//	force := pixel.V(5+common.GetRandFloat()*5, 5+common.GetRandFloat()*5)

	particles[idx] = newParticle(pos, f, 0.1, common.GetRandFloat()*3, col)
}

func Update(dt float64, objs []common.Objecter) {
	for i := range particles {
		if particles[i].active {
			particles[i].update(dt, objs)
		}
	}
}

func Draw(t pixel.Target) {
	imd := imdraw.New(nil)
	for i := range particles {
		if particles[i].active {
			particles[i].draw(imd)
		}
	}
	imd.Draw(t)
}
