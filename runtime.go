package main

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"os"
	"platformer/bindings"
	"platformer/common"
	"platformer/config"
	"platformer/sound"

	"github.com/shinomontaz/pixel/pixelgl"

	"github.com/shinomontaz/pixel"
)

var (
	isfullscreen bool
	winHeight    float64
	winWidth     float64
	cfgloader    *common.Loader
	assetloader  *common.Loader
	startConfig  StartConfig
	runtimecfg   *os.File
)

type StartConfig struct {
	TestFlag bool   `json:"TestFlag"`
	Assets   string `json:"Assets"`
	Zip      string `json:"Zip"`
	Configs  string `json:"Configs"`
}

// get window mode, sound volumes and menu/game mode
// currBounds
func initRuntime() {

	var byteValue []byte
	fconfig, err := os.Open("start.cfg")
	if err != nil {
		log.Fatal(err)
	}
	defer fconfig.Close()

	byteValue, _ = ioutil.ReadAll(fconfig)
	json.Unmarshal(byteValue, &startConfig)

	// init configs

	if startConfig.Zip == "" {
		cfgloader = common.NewLoader(startConfig.Configs)
		assetloader = common.NewLoader(startConfig.Assets)
	} else {
		cfgloader = common.NewLoader(startConfig.Configs, common.WithZip(startConfig.Zip))
		assetloader = common.NewLoader(startConfig.Assets, common.WithZip(startConfig.Zip))
	}

	runtimecfg, err = os.OpenFile("runtime.cfg", os.O_RDWR, 0755)
	if err != nil {
		log.Fatal(err)
	}

	isdebug = startConfig.TestFlag

	err = config.Init(cfgloader, runtimecfg)
	if err != nil {
		log.Fatal(err)
	}

	winWidth = config.Opts.WindowWidth
	winHeight = config.Opts.WindowHeight

	currBounds = pixel.R(0, 0, winWidth, winHeight)

	common.InitFont(assetloader)
	bindings.Init()
	//	magic.Init(assetloader)
}

func initScreen(win *pixelgl.Window) {
	if config.Opts.Fullscreen {
		win.SetMonitor(pixelgl.PrimaryMonitor())
		win.SetBounds(currBounds)
		win.SetCursorVisible(false)
	} else {
		win.SetMonitor(nil)
		win.SetBounds(currBounds)
		win.SetCursorVisible(true)
	}
}

func initSound() {
	sound.SetVolumes(config.Opts.Volumes["main"], (config.Opts.Volumes["music"] - 100), (config.Opts.Volumes["actions"] - 100))
}
