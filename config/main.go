package config

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"os"
)

var (
	MainConfig   Main
	PlayerConfig Player
	WorldConfig  World
	AnimConfig   []Anims
)

func init() {
	var byteValue []byte
	fconfig, err := os.Open("config/config.json")
	if err != nil {
		log.Fatal(err)
	}
	defer fconfig.Close()

	byteValue, _ = ioutil.ReadAll(fconfig)
	json.Unmarshal(byteValue, &MainConfig)

	fworld, err := os.Open(fmt.Sprintf("config/%s", MainConfig.WorldCfg))
	if err != nil {
		log.Fatal(err)
	}
	defer fworld.Close()

	byteValue, _ = ioutil.ReadAll(fworld)
	json.Unmarshal(byteValue, &WorldConfig)

	fanims, err := os.Open(fmt.Sprintf("config/%s", MainConfig.AllAnims))
	if err != nil {
		log.Fatal(err)
	}
	defer fanims.Close()

	byteValue, _ = ioutil.ReadAll(fanims)
	json.Unmarshal(byteValue, &AnimConfig)

	fplayer, err := os.Open(fmt.Sprintf("config/%s", MainConfig.PlayerCfg))
	if err != nil {
		log.Fatal(err)
	}
	defer fplayer.Close()

	byteValue, _ = ioutil.ReadAll(fplayer)
	json.Unmarshal(byteValue, &PlayerConfig)
}

func (a *Anims) W() float64 {
	return float64(a.Width)
}

func (a *Anims) H() float64 {
	return float64(a.Height)
}

func (a *Anims) M() float64 {
	return float64(a.Margin)
}

func (a *Anims) N() string {
	return a.Name
}

func (a *Anims) Get() ([]string, []string, [][]int) {
	names := make([]string, 0)
	files := make([]string, 0)
	frames := make([][]int, 0)

	for _, an := range a.List {
		names = append(names, an.Name)
		files = append(files, an.File)
		frames = append(frames, an.Frames)
	}

	return names, files, frames
}
