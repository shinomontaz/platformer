package sound

import (
	"fmt"
	"log"
	"platformer/common"
	"platformer/config"
	"time"

	"github.com/faiface/beep"
	"github.com/faiface/beep/effects"
	"github.com/faiface/beep/mp3"
	"github.com/faiface/beep/speaker"
	"github.com/faiface/beep/wav"
	"github.com/shinomontaz/pixel"
)

type Sound struct {
	buff *beep.Buffer
}

type PosEffect struct {
	vol *effects.Volume
	pos pixel.Vec
}

var (
	listener     pixel.Vec
	soundeffects map[string]Sound
	music        map[string]Sound

	currEffects []PosEffect
	currMusic   *effects.Volume

	volBase    float64
	volMusic   *effects.Volume
	volEffects *effects.Volume
)

var sbrs []common.Subscriber

func AddSubscriber(sbr common.Subscriber) {
	sbrs = append(sbrs, sbr)
}

func Notify(e int) {
	for _, s := range sbrs {
		s.Listen(e, pixel.ZV)
	}
}

// main, music, actions in [-100, 100]
func SetVolumes(main, music, actions float64) {
	volBase = main/100 - 0.5

	volMusic = &effects.Volume{
		Base:   2,
		Volume: volBase + music/10,
	}
	if music/100 <= -1 || volBase == -0.5 {
		volMusic.Silent = true
	}
	volEffects = &effects.Volume{
		Base:   2,
		Volume: volBase + actions/10,
	}

	if actions/100 <= -1 || volBase == -0.5 {
		volEffects.Silent = true
	}

	if currMusic != nil {
		currMusic.Volume = volMusic.Volume
		currMusic.Silent = volMusic.Silent
	}
}

func Init(loader *common.Loader) {
	music = make(map[string]Sound)
	soundeffects = make(map[string]Sound)
	currEffects = make([]PosEffect, 0)
	sbrs = make([]common.Subscriber, 0)

	// read effects
	// read music
	f, err := loader.Open("sounds/music/fluffily-11859.mp3")
	if err != nil {
		log.Fatal(err)
	}
	defer f.Close()

	streamer, format, err := mp3.Decode(f)
	if err != nil {
		log.Fatal(err)
	}
	buff := beep.NewBuffer(format)
	buff.Append(streamer)
	streamer.Close()

	speaker.Init(format.SampleRate, format.SampleRate.N(time.Second/10))

	music["main"] = Sound{
		buff: buff,
	}

	for _, sp := range config.Sounds {
		for _, eff := range sp.List {
			for _, ef := range eff.List {
				if _, ok := soundeffects[ef]; ok {
					continue
				}
				// read soundeffect
				f, err := loader.Open(ef)
				if err != nil {
					log.Fatal(err)
				}

				streamer, _, err := wav.Decode(f)
				if err != nil {
					log.Fatal(err)
				}
				buff := beep.NewBuffer(format)
				buff.Append(streamer)
				streamer.Close()
				f.Close()

				soundeffects[ef] = Sound{
					buff: buff,
				}

			}
		}
	}
}

func Update(pos pixel.Vec) {
	listener = pos
	i := 0
	for _, ce := range currEffects {
		l := pixel.L(listener, ce.pos).Len()
		if l > 500 {
			continue
		}
		currEffects[i] = ce
		ce.vol.Volume = volEffects.Volume - l/50
		ce.vol.Silent = volEffects.Silent
		i++
	}
	currEffects = currEffects[:i]
}

func AddEffect(name string, pos pixel.Vec) {
	l := pixel.L(listener, pos).Len()
	if l > 500 {
		return
	}

	bfr := soundeffects[name].buff
	eft := bfr.Streamer(0, bfr.Len())

	vol := &effects.Volume{
		Base:     volEffects.Base,
		Streamer: eft,
		Volume:   volEffects.Volume - l/50, //-l / 50,
		Silent:   volEffects.Silent,
	}

	fmt.Println("AddEffect", vol)

	pe := PosEffect{
		vol,
		pos,
	}
	currEffects = append(currEffects, pe)

	speaker.Play(vol)
}

func PlauseMusic() {
	//	ctrl.Paused
}

func PlayMusic(name string) {
	speaker.Clear()
	bfr := music[name].buff
	mus := bfr.Streamer(0, bfr.Len())
	ctrl := &beep.Ctrl{Streamer: beep.Loop(-1, mus), Paused: false}

	currMusic = &effects.Volume{
		Base:     volMusic.Base,
		Streamer: ctrl,
		Volume:   volMusic.Volume,
		Silent:   volMusic.Silent,
	}

	speaker.Play(currMusic)
}
