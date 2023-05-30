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

	reg8, err := loadFont("fonts/MatchupPro.ttf", 8)
	if err != nil {
		panic(err)
	}
	fonts["regular8"] = reg8

	reg16, err := loadFont("fonts/MatchupPro.ttf", 16)
	if err != nil {
		panic(err)
	}
	fonts["regular16"] = reg16

	reg20, err := loadFont("fonts/MatchupPro.ttf", 20)
	if err != nil {
		panic(err)
	}
	fonts["regular20"] = reg20

	reg24, err := loadFont("fonts/MatchupPro.ttf", 24)
	if err != nil {
		panic(err)
	}
	fonts["regular24"] = reg24

	reg28, err := loadFont("fonts/MatchupPro.ttf", 28)
	if err != nil {
		panic(err)
	}
	fonts["regular28"] = reg28

	reg32, err := loadFont("fonts/MatchupPro.ttf", 32)
	if err != nil {
		panic(err)
	}
	fonts["regular32"] = reg32

	//	menu, err := loadFont("fonts/capture.ttf", 20)
	menu28, err := loadFont("fonts/CompassPro.ttf", 28)

	if err != nil {
		panic(err)
	}
	fonts["menu28"] = menu28

	menu20, err := loadFont("fonts/CompassPro.ttf", 20)
	if err != nil {
		panic(err)
	}
	fonts["menu20"] = menu20

	fancy20, err := loadFont("fonts/NicerNightie.ttf", 20)
	if err != nil {
		panic(err)
	}
	fonts["fancy20"] = fancy20

	fancy32, err := loadFont("fonts/NicerNightie.ttf", 32)
	if err != nil {
		panic(err)
	}
	fonts["fancy32"] = fancy32

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
