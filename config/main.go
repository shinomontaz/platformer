package config

import (
	"encoding/json"
	"io"
	"log"
	"os"
	"platformer/common"

	"github.com/shinomontaz/pixel"
)

var (
	MainConfig Main
	//	PlayerConfig Player
	AnimConfig []Anims
	Sounds     map[string]Soundprofile
	Profiles   map[string]Profile
	Loots      map[string]Profile
	Spells     map[string]Spellprofile
	Opts       Options
	runtimecfg *os.File
)

func Init(loader *common.Loader, rtcfg *os.File) error {
	runtimecfg = rtcfg
	var byteValue []byte
	fconfig, err := loader.Open("config.json")
	if err != nil {
		return err
	}
	defer fconfig.Close()

	byteValue, _ = io.ReadAll(fconfig)
	json.Unmarshal(byteValue, &MainConfig)

	//	fanims, err := os.Open(fmt.Sprintf("config/%s", MainConfig.AllAnims))
	fanims, err := loader.Open(MainConfig.AllAnims)
	if err != nil {
		return err
	}
	defer fanims.Close()

	byteValue, _ = io.ReadAll(fanims)
	json.Unmarshal(byteValue, &AnimConfig)

	fprofiles, err := loader.Open(MainConfig.Profiles)
	if err != nil {
		return err
	}
	defer fprofiles.Close()

	sliceProfiles := make([]Profile, 0)
	byteValue, _ = io.ReadAll(fprofiles)
	json.Unmarshal(byteValue, &sliceProfiles)

	Profiles = make(map[string]Profile)
	for _, pr := range sliceProfiles {
		Profiles[pr.Type] = pr
	}

	floots, err := loader.Open(MainConfig.Loots)
	if err != nil {
		return err
	}
	defer floots.Close()
	sliceLoots := make([]Profile, 0)
	byteValue, _ = io.ReadAll(floots)
	json.Unmarshal(byteValue, &sliceLoots)
	Loots = make(map[string]Profile)
	for _, pr := range sliceLoots {
		Loots[pr.Type] = pr
	}

	sprofiles, err := loader.Open(MainConfig.Sounds)
	if err != nil {
		return err
	}
	defer sprofiles.Close()

	sliceSprofiles := make([]Soundprofile, 0)
	byteValue, _ = io.ReadAll(sprofiles)
	json.Unmarshal(byteValue, &sliceSprofiles)

	Sounds = make(map[string]Soundprofile)
	for _, pr := range sliceSprofiles {
		Sounds[pr.Type] = pr
	}

	fspells, err := loader.Open(MainConfig.Spells)
	if err != nil {
		return err
	}
	defer fspells.Close()

	byteValue, _ = io.ReadAll(fspells)
	json.Unmarshal(byteValue, &Spells)

	byteValue, _ = io.ReadAll(runtimecfg)
	json.Unmarshal(byteValue, &Opts)
	return nil
}

func SaveRuntime() {
	json, err := json.Marshal(Opts)
	if err != nil {
		panic("Failed to save configuration 1")
	}
	runtimecfg.Truncate(0)
	if _, err = runtimecfg.WriteAt(json, 0); err != nil {
		log.Fatal(err)
		//		panic("Failed to save configuration 2")
	}
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
