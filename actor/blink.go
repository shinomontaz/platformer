package actor

import (
	"log"
	"math"
	"platformer/common/rand"

	"github.com/shinomontaz/pixel"
)

type span struct{ a, b float64 }

func cut(in []span, a, b float64) []span {
	out := make([]span, 0, len(in)+1)
	for _, s := range in {
		if b <= s.a || a >= s.b {
			out = append(out, s)
			continue
		}
		if a > s.a {
			out = append(out, span{s.a, a})
		}
		if b < s.b {
			out = append(out, span{b, s.b})
		}
	}
	return out
}

func (a *Actor) blinkSpot(avoid pixel.Vec, maxRise float64) (pixel.Vec, bool) {
	const margin = 6.0
	feet := a.phys.GetRect()          // коллизионный прямоугольник — им стоим на площадке
	offY := feet.Min.Y - a.rect.Min.Y // насколько ноги выше низа отрисовки
	halfW := feet.W() / 2

	if !a.hasHome {
		log.Printf("[blink] no home platform yet")
		return pixel.ZV, false
	}
	plat := a.homePlat

	lo, hi := plat.Min.X+halfW+margin, plat.Max.X-halfW-margin
	if lo >= hi {
		log.Printf("[blink] platform too narrow %v", plat)
		return pixel.ZV, false
	}
	free := []span{{lo, hi}}

	band := pixel.R(plat.Min.X, plat.Max.Y, plat.Max.X, plat.Max.Y+maxRise+feet.H())
	for _, o := range a.currObjs {
		pr := o.Rect()
		if pr == plat || !pr.Intersects(band) {
			continue
		}
		free = cut(free, pr.Min.X-halfW-margin, pr.Max.X+halfW+margin)
	}

	keep := halfW + 60
	free = cut(free, avoid.X-keep, avoid.X+keep)
	if len(free) == 0 {
		log.Printf("[blink] no free span on %v (hero x=%.0f)", plat, avoid.X)
		return pixel.ZV, false
	}

	total := 0.0
	for _, s := range free {
		total += s.b - s.a
	}
	pick := rand.Float64() * total
	x := free[len(free)-1].b
	for _, s := range free {
		l := s.b - s.a
		if pick <= l {
			x = s.a + pick
			break
		}
		pick -= l
	}

	minRise := feet.H() / 2
	rise := minRise + rand.Float64()*math.Max(0, maxRise-minRise)

	// ноги окажутся на plat.Max.Y+rise; Teleport ставит центр a.rect
	rectMinY := plat.Max.Y + rise - offY
	dx := a.rect.Center().X - feet.Center().X // bbox может быть несимметричным
	return pixel.V(x+dx, rectMinY+a.rect.H()/2), true
}

func (a *Actor) Blink(avoid pixel.Vec, maxRise float64) (pixel.Vec, bool) {
	spot, ok := a.blinkSpot(avoid, maxRise)
	if !ok {
		return pixel.ZV, false
	}
	moved := spot.Sub(a.rect.Center())
	a.Teleport(spot)
	return moved, true
}
