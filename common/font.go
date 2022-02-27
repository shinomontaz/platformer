package common

import (
	"io/ioutil"

	"github.com/golang/freetype/truetype"
	"golang.org/x/image/font"
)

var fonts map[string]font.Face

func init() {
	fonts = make(map[string]font.Face)
	bold, err := loadFont("assets/fonts/GravityBold8.ttf", 8)
	if err != nil {
		panic(err)
	}
	fonts["bold"] = bold

	reg, err := loadFont("assets/fonts/GravityRegular5.ttf", 8)
	if err != nil {
		panic(err)
	}
	fonts["regular"] = reg

	menu, err := loadFont("assets/fonts/capture.ttf", 20)
	if err != nil {
		panic(err)
	}
	fonts["menu"] = menu

	menusmall, err := loadFont("assets/fonts/capture.ttf", 12)
	if err != nil {
		panic(err)
	}
	fonts["menusmall"] = menusmall

}

func GetFont(s string) font.Face {
	return fonts[s]
}

func loadFont(path string, size float64) (font.Face, error) {
	b, err := ioutil.ReadFile(path)
	if err != nil {
		return nil, err
	}

	f, err := truetype.Parse(b)
	if err != nil {
		return nil, err
	}
	face := truetype.NewFace(f, &truetype.Options{
		Size:              size,
		Hinting:           font.HintingFull,
		GlyphCacheEntries: 1,
	})
	return face, nil
}
