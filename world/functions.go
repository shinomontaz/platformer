package world

import (
	"errors"
	"image/color"

	"github.com/shinomontaz/pixel"
)

func tileIDToCoord(tID int, numColumns int, numRows int) (x int, y int) {
	x = tID % numColumns
	y = numRows - (tID / numColumns) - 1
	return
}

func indexToGamePos(idx int, width int, height int) pixel.Vec {
	gamePos := pixel.V(
		float64(idx%width),
		float64(height)-float64(idx/width)-1,
	)
	return gamePos
}

func Box(from, to pixel.Vec) pixel.Rect {
	var minx, miny, maxx, maxy float64
	if from.X > to.X {
		minx = to.X
		maxx = from.X
	} else {
		minx = from.X
		maxx = to.X
	}

	if from.Y > to.Y {
		miny = to.Y
		maxy = from.Y
	} else {
		miny = from.Y
		maxy = to.Y
	}

	return pixel.R(minx, miny, maxx, maxy)
}

var errInvalidFormat = errors.New("invalid format")

func ParseHexColor(s string) (c color.RGBA, err error) {
	c.A = 0xff

	if s[0] != '#' {
		return c, errInvalidFormat
	}

	hexToByte := func(b byte) byte {
		switch {
		case b >= '0' && b <= '9':
			return b - '0'
		case b >= 'a' && b <= 'f':
			return b - 'a' + 10
		case b >= 'A' && b <= 'F':
			return b - 'A' + 10
		}
		err = errInvalidFormat
		return 0
	}

	switch len(s) {
	case 7:
		c.R = hexToByte(s[1])<<4 + hexToByte(s[2])
		c.G = hexToByte(s[3])<<4 + hexToByte(s[4])
		c.B = hexToByte(s[5])<<4 + hexToByte(s[6])
	case 4:
		c.R = hexToByte(s[1]) * 17
		c.G = hexToByte(s[2]) * 17
		c.B = hexToByte(s[3]) * 17
	default:
		err = errInvalidFormat
	}
	return
}
