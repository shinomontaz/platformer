package common

import (
	"io/ioutil"

	"github.com/golang/freetype/truetype"
	"golang.org/x/image/font"
)

var fonts map[string]font.Face

var loader *Loader

func InitFont(l *Loader) {
	loader = l
	fonts = make(map[string]font.Face)
	bold, err := loadFont("fonts/GravityBold8.ttf", 8)
	if err != nil {
		panic(err)
	}
	fonts["bold"] = bold

	reg, err := loadFont("fonts/GravityRegular5.ttf", 8)
	if err != nil {
		panic(err)
	}
	fonts["regular"] = reg

	menu, err := loadFont("fonts/capture.ttf", 20)
	if err != nil {
		panic(err)
	}
	fonts["menu"] = menu

	menusmall, err := loadFont("fonts/capture.ttf", 12)
	if err != nil {
		panic(err)
	}
	fonts["menusmall"] = menusmall

}

func GetFont(s string) font.Face {
	return fonts[s]
}

func loadFont(path string, size float64) (font.Face, error) {
	file, err := loader.Open(path)
	if err != nil {
		return nil, err
	}
	b, err := ioutil.ReadAll(file)
	if err != nil {
		return nil, err
	}
	file.Close()

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
