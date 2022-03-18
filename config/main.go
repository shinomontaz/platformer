package config

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"os"

	"github.com/faiface/pixel"
)

var (
	MainConfig Main
	//	PlayerConfig Player
	AnimConfig []Anims
	Sounds     map[string]Soundprofile
	Profiles   map[string]Profile
	Spells     map[string]Spellprofile
	Opts       Options
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

	fanims, err := os.Open(fmt.Sprintf("config/%s", MainConfig.AllAnims))
	if err != nil {
		log.Fatal(err)
	}
	defer fanims.Close()

	byteValue, _ = ioutil.ReadAll(fanims)
	json.Unmarshal(byteValue, &AnimConfig)

	fprofiles, err := os.Open(fmt.Sprintf("config/%s", MainConfig.Profiles))
	if err != nil {
		log.Fatal(err)
	}
	defer fprofiles.Close()

	sliceProfiles := make([]Profile, 0)
	byteValue, _ = ioutil.ReadAll(fprofiles)
	json.Unmarshal(byteValue, &sliceProfiles)

	Profiles = make(map[string]Profile)
	for _, pr := range sliceProfiles {
		Profiles[pr.Type] = pr
	}

	sprofiles, err := os.Open(fmt.Sprintf("config/%s", MainConfig.Sounds))
	if err != nil {
		log.Fatal(err)
	}
	defer sprofiles.Close()

	sliceSprofiles := make([]Soundprofile, 0)
	byteValue, _ = ioutil.ReadAll(sprofiles)
	json.Unmarshal(byteValue, &sliceSprofiles)

	Sounds = make(map[string]Soundprofile)
	for _, pr := range sliceSprofiles {
		Sounds[pr.Type] = pr
	}

	fspells, err := os.Open(fmt.Sprintf("config/%s", MainConfig.Spells))
	if err != nil {
		log.Fatal(err)
	}
	defer fspells.Close()

	byteValue, _ = ioutil.ReadAll(fspells)
	json.Unmarshal(byteValue, &Spells)

	foptions, err := os.Open("config/options.json")
	if err != nil {
		log.Fatal(err)
	}
	defer foptions.Close()

	byteValue, _ = ioutil.ReadAll(foptions)
	json.Unmarshal(byteValue, &Opts)
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

func (a *Anims) GetGroups() (map[string][]string, map[string][]string, map[string][][]int) {
	group_names := make(map[string][]string)
	group_files := make(map[string][]string)
	group_frames := make(map[string][][]int)

	for _, gr := range a.Groups {
		if _, ok := group_names[gr.Name]; !ok {
			group_names[gr.Name] = make([]string, 0)
			group_files[gr.Name] = make([]string, 0)
			group_frames[gr.Name] = make([][]int, 0)
		}
		for _, an := range gr.List {
			group_names[gr.Name] = append(group_names[gr.Name], an.Name)
			group_files[gr.Name] = append(group_files[gr.Name], an.File)
			group_frames[gr.Name] = append(group_frames[gr.Name], an.Frames)
		}
	}

	return group_names, group_files, group_frames
}

func (sp *Spellprofile) GetHitbox() pixel.Rect {
	return pixel.R(sp.Hitbox[0], sp.Hitbox[1], sp.Hitbox[2], sp.Hitbox[3])
}

func (sp *Spellprofile) GetSound() string {
	return sp.Sound
}
