package talks

import (
	"encoding/json"
	"io"
	"math"
	"platformer/common"

	"github.com/shinomontaz/pixel"
	"golang.org/x/image/colornames"
)

var phrases map[string]Phrase

type Phrase struct {
	Class    string   `json:"class"`
	Variants []string `json:"variants"`
}

func initPhrases(loader *common.Loader) {
	phrases = make(map[string]Phrase)

	tphrases := make([]Phrase, 0)
	phrs, err := loader.Open("phrases.json")
	if err != nil {
		panic(err)
	}
	defer phrs.Close()

	byteValue, _ := io.ReadAll(phrs)
	json.Unmarshal(byteValue, &tphrases)

	for _, p := range tphrases { // TODO: unmarshal to map[string]Phrase
		phrases[p.Class] = p
	}
}

func AddPhrase(pos pixel.Vec, class string) {
	if p, ok := phrases[class]; ok {
		// get random string
		txt := p.Variants[int(math.Round(common.GetRandFloat()*float64(len(p.Variants)-1)))]
		al := addAlert(pos, colornames.Green, txt, 2, 1)
		alerts = append(alerts, al)
	}
}

func GetPhrases(key string) *Phrase {
	if p, ok := phrases[key]; ok {
		return &p
	}
	return nil
}
