package sound

import (
	"log"
	"os"
	"platformer/config"
	"time"

	"github.com/faiface/beep"
	"github.com/faiface/beep/effects"
	"github.com/faiface/beep/mp3"
	"github.com/faiface/beep/speaker"
	"github.com/faiface/beep/wav"
	"github.com/faiface/pixel"
)

type Sound struct {
	buff *beep.Buffer
}

type PosEffect struct {
	s   Sound
	pos pixel.Vec
}

var (
	listener     pixel.Vec
	soundeffects map[string]Sound
	music        map[string]Sound

	currEffects []PosEffect

	volMusic   *effects.Volume
	volEffects *effects.Volume
)

func init() {
	music = make(map[string]Sound)
	soundeffects = make(map[string]Sound)
	currEffects = make([]PosEffect, 0)

	volMusic = &effects.Volume{
		Base:   2,
		Volume: -5,
		Silent: false,
	}
	volEffects = &effects.Volume{
		Base:   2,
		Volume: 0,
		Silent: false,
	}

	// read effects
	// read music
	f, err := os.Open("assets/sounds/music/fluffily-11859.mp3")
	if err != nil {
		log.Fatal(err)
	}

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
				f, err := os.Open(ef)
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

				soundeffects[ef] = Sound{
					buff: buff,
				}

			}
		}
	}
}

func Update(pos pixel.Vec) {
	listener = pos
}

func AddEffect(name string, pos pixel.Vec) {
	l := pixel.L(listener, pos).Len()
	if l > 500 {
		return
	}
	// currEffects = append(currEffects, PosEffect{
	// 	s:   soundeffects[name],
	// 	pos: pos,
	// })

	bfr := soundeffects[name].buff
	eft := bfr.Streamer(0, bfr.Len())

	volEffects.Streamer = eft
	//	volEffects.Volume -= l / 50

	speaker.Play(volEffects)
}

func PlayMusic(name string) {
	bfr := music[name].buff
	mus := bfr.Streamer(0, bfr.Len())
	ctrl := &beep.Ctrl{Streamer: beep.Loop(-1, mus), Paused: false}

	volMusic.Streamer = ctrl

	speaker.Play(volMusic)
}
