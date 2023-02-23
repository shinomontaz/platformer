package common

import (
	"math"

	"github.com/shinomontaz/pixel"
)

// sweptAABB implemented here
//Box b1, Box b2, float& normalx, float& normaly
func collide(p, q pixel.Rect, v pixel.Vec) (float64, pixel.Vec) {
	n := pixel.ZV
	var xInvEntry, yInvEntry, xInvExit, yInvExit float64
	if v.X > 0 {
		xInvEntry = q.Min.X - p.Max.X
		xInvExit = q.Max.X - p.Min.X
	} else {
		xInvEntry = q.Max.X - p.Min.X
		xInvExit = q.Min.X - p.Max.X
	}

	if v.Y > 0 {
		yInvEntry = q.Min.Y - p.Max.Y
		yInvExit = q.Max.Y - p.Min.Y
	} else {
		yInvEntry = q.Max.Y - p.Min.Y
		yInvExit = q.Min.Y - p.Max.Y
	}

	var xEntry, yEntry, xExit, yExit float64

	if v.X == 0 {
		xEntry = math.Inf(-1)
		xExit = math.Inf(1)
	} else {
		xEntry = xInvEntry / v.X
		xExit = xInvExit / v.X
	}

	if v.Y == 0 {
		yEntry = math.Inf(-1)
		yExit = math.Inf(1)
	} else {
		yEntry = yInvEntry / v.Y
		yExit = yInvExit / v.Y
	}

	entryTime := math.Max(xEntry, yEntry)
	exitTime := math.Min(xExit, yExit)

	if entryTime > exitTime || (xEntry < 0 && yEntry < 0) || xEntry > 1 || yEntry > 1 {
		return 1, n
	} else {
		if xEntry > yEntry {
			if xInvEntry < 0 {
				n.X = 1
			} else {
				n.X = -1
			}
		} else {
			if yInvEntry < 0 {
				n.Y = 1
			} else {
				n.Y = -1
			}
		}
	}

	return entryTime, n
}

func Isinbox(p, q pixel.Rect) bool {
	//return !(b1.x + b1.w < b2.x || b1.x > b2.x + b2.w || b1.y + b1.h < b2.y || b1.y > b2.y + b2.h);
	return !(p.Max.X < q.Min.X || p.Min.X > q.Max.X || p.Max.Y < q.Min.Y || p.Min.Y > q.Max.Y)
}

func Broadbox(r pixel.Rect, v pixel.Vec) pixel.Rect {
	var minx, miny, w, h float64
	if v.X > 0 {
		minx = r.Min.X
		w = v.X + r.W()
	} else {
		minx = r.Min.X + v.X
		w = r.W() - v.X
	}
	if v.Y > 0 {
		miny = r.Min.Y
		h = v.Y + r.H()
	} else {
		miny = r.Min.Y + v.Y
		h = r.H() - v.Y
	}

	return pixel.R(minx, miny, minx+w, miny+h)
}
