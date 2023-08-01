package dialogs

import (
	"encoding/json"
	"fmt"
	"strings"
	"unicode"

	"platformer/actor"

	"github.com/shinomontaz/pixel"
	"github.com/shinomontaz/pixel/imdraw"
	"github.com/shinomontaz/pixel/text"
	"golang.org/x/image/colornames"
)

var margin = 10.0
var portraitwidth = 16.0

type Dialog struct {
	ID            int                   `json:"id"`
	Variants      map[int]DialogVariant `json:"variants"`
	currVariant   int
	currAnswer    int
	imd           *imdraw.IMDraw
	rect          pixel.Rect
	a             *actor.Actor
	maintxtstring string
	maintxt       *text.Text
}

type DialogVariant struct {
	Idx     int            `json:"idx"`
	Text    string         `json:"text"`
	Answers []DialogAnswer `json:"answers"`
}

type DialogAnswer struct {
	Text string `json:"text"`
	Goto int    `json:"goto,omitempty"`
	Code int    `json:"code,omitempty"`
	Exit bool   `json:"exit,omitempty"`
}

func (d *Dialog) UnmarshalJSON(b []byte) error {
	type JSONDialog struct {
		ID       int             `json:"id"`
		Variants []DialogVariant `json:"variants"`
	}
	var jDlg JSONDialog
	err := json.Unmarshal(b, &jDlg)
	if err != nil {
		return err
	}

	d.ID = jDlg.ID
	d.Variants = make(map[int]DialogVariant)
	for _, dv := range jDlg.Variants {
		if d.currVariant == 0 {
			d.currVariant = dv.Idx
			d.currAnswer = 0
		}
		d.Variants[dv.Idx] = dv
	}

	return nil
}

func (d *Dialog) SetVariant(id int) {
	if _, ok := d.Variants[id]; ok {
		d.currVariant = id
		d.currAnswer = 0
	}
	d.prepareText()
}

func (d *Dialog) Start(bounds pixel.Rect) {
	d.rect = pixel.R(0, 0, 300, 240)
	d.rect = d.rect.Moved(bounds.Center().Sub(d.rect.Center()))
	d.imd = imdraw.New(nil)
	d.setImdr()
	d.maintxt = text.New(pixel.V(0, 0), atlas)

	d.prepareText()
}

func (d *Dialog) Draw(t pixel.Target) {
	d.imd.Draw(t)

	pos := pixel.V(d.rect.Min.X+margin+margin+32, d.rect.Max.Y-2*margin)

	maintxtrect := d.maintxt.Bounds().Moved(pos)

	portrait := d.a.GetPortrait()
	portrait.Draw(t, pixel.IM.Moved(pixel.Vec{d.rect.Min.X + margin + portraitwidth, d.rect.Max.Y - margin - portraitwidth}))

	d.maintxt.Draw(t, pixel.IM.Moved(maintxtrect.Min))

	// imd := imdraw.New(nil)
	// vertices := maintxtrect.Vertices()
	// imd.Color = colornames.Red
	// for _, v := range vertices {
	// 	imd.Push(v)
	// }
	// imd.Rectangle(1)
	// imd.Draw(t)

	pointertxt := text.New(pixel.V(0, 0), atlasbig)
	pointertxt.Color = colornames.Aliceblue
	fmt.Fprintln(pointertxt, "->")

	anstxt := text.New(pixel.V(0, 0), atlasbig)
	h := 0.0
	for i, ans := range d.Variants[d.currVariant].Answers {
		anstxt.Clear()
		anstxt.Color = colornames.Whitesmoke
		if i == d.currAnswer {
			anstxt.Color = colornames.Aliceblue
			pointertxt.Draw(t, pixel.IM.Moved(pixel.Vec{maintxtrect.Min.X - 2*margin, maintxtrect.Min.Y - 4*margin - h}))
		}
		fmt.Fprintln(anstxt, ans.Text)
		anstxt.Draw(t, pixel.IM.Moved(pixel.Vec{maintxtrect.Min.X, maintxtrect.Min.Y - 4*margin - h}))
		h += anstxt.Bounds().H() + margin
	}
}

func (d *Dialog) Action() {
	go_to := d.Variants[d.currVariant].Answers[d.currAnswer].Goto
	fmt.Println("dialog action:", go_to)

	code := d.Variants[d.currVariant].Answers[d.currAnswer].Code
	if code > 0 {
		runAction(code, d.a)
	}

	exit := d.Variants[d.currVariant].Answers[d.currAnswer].Exit
	if exit {
		UnsetActive()
		return
	}

	if go_to != 0 {
		d.SetVariant(go_to)
	}
}

func (d *Dialog) UpdateAnswer(i int) {
	d.currAnswer += i
	if d.currAnswer < 0 {
		d.currAnswer = len(d.Variants[d.currVariant].Answers) - 1
	}
	d.currAnswer %= len(d.Variants[d.currVariant].Answers)
	fmt.Println("d.currAnswer: ", d.currAnswer)
}

func (d *Dialog) setImdr() {
	d.imd.Clear()
	d.imd.Color = colornames.Darkslategray
	d.imd.Push(d.rect.Min)
	d.imd.Push(d.rect.Max)
	d.imd.Rectangle(0)

	d.imd.Color = colornames.Darkgray
	d.imd.Push(d.rect.Min.Add(pixel.Vec{3, 3}))
	d.imd.Push(d.rect.Max.Sub(pixel.Vec{3, 3}))
	d.imd.Rectangle(2)
}

func (d *Dialog) prepareText() {
	d.maintxtstring = d.Variants[d.currVariant].Text
	d.maintxt.Color = colornames.Whitesmoke
	fmt.Fprintln(d.maintxt, d.maintxtstring)

	for d.maintxt.Bounds().W() > d.rect.W()-2*margin-(margin+portraitwidth) {
		d.maintxtstring = splitToChunks(d.maintxtstring)
		d.maintxt.Clear()
		fmt.Fprintln(d.maintxt, d.maintxtstring)
	}
}

func splitToChunks(s string) string {
	if len(s) <= 1 {
		return s
	}

	i := strings.Index(s, " ")
	if i == 0 {
		return s
	}

	substrs := strings.Split(s, "\n")
	longest_idx := 0
	longest_len := len(substrs[longest_idx])
	for i := 0; i < len(substrs); i++ {
		curr_len := len(substrs[i])
		if curr_len > longest_len {
			longest_len = curr_len
			longest_idx = i
		}
	}

	fixed := []string{subsplitToChunks(substrs[longest_idx])}
	fixed = append(fixed, substrs[longest_idx+1:]...)
	substrs = append(substrs[:longest_idx], fixed...)

	return strings.Join(substrs, "\n")
}

func subsplitToChunks(s string) string {
	if len(s) <= 1 {
		return s
	}

	i := strings.Index(s, " ")
	if i == 0 {
		return s
	}

	i = strings.Index(s, "\n")
	if i != -1 {
		return s
	}

	currLen := 0
	res := make([]string, 3)
	for i := range s {
		if currLen >= len(s)/2 && unicode.IsSpace(rune(s[i])) {
			res = []string{s[:i], "\n", s[i+1:]}
			break
		}
		currLen++
	}

	return strings.Join(res, "")
}
