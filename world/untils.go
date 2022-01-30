package world

import (
	"fmt"
	"image/png"
	"os"

	"github.com/faiface/pixel"
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

func loadSprite(path string) (*pixel.Sprite, *pixel.PictureData) {
	f, err := os.Open(path)
	if err != nil {
		fmt.Println(path)
		panic(err)
	}

	img, err := png.Decode(f)
	if err != nil {
		panic(err)
	}

	pd := pixel.PictureDataFromImage(img)
	return pixel.NewSprite(pd, pd.Bounds()), pd
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
