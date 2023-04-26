package animation

import (
	"platformer/common"

	"github.com/shinomontaz/pixel"
)

var anims map[string]*Anims

func init() {
	anims = make(map[string]*Anims)
}

type Anim struct {
	sheet  pixel.Picture
	frames []pixel.Rect
	sprite *pixel.Sprite
}

type Anims struct {
	items  map[string]*Anim
	groups map[string]map[string]*Anim
	rect   pixel.Rect
	sprite *pixel.Sprite
	m      [4]float64 // margin
}

var loader *common.Loader

func Init(l *common.Loader) {
	loader = l
}

func Load(cfg AnimatingConfig) {
	animRect := pixel.R(0, 0, cfg.W(), cfg.H())
	a := New(animRect, cfg.M())
	names, files, frames := cfg.Get()
	for i := 0; i < len(names); i++ {
		an, err := a.PrepareAnim(names[i], files[i], frames[i])
		if err == nil {
			a.items[names[i]] = an
		}
	}

	grnames, grfiles, grframes := cfg.GetGroups()
	for grname, animnames := range grnames {
		if _, ok := a.groups[grname]; !ok {
			a.groups[grname] = make(map[string]*Anim)
			for i := 0; i < len(animnames); i++ {
				an, err := a.PrepareAnim(animnames[i], grfiles[grname][i], grframes[grname][i])
				if err == nil {
					a.groups[grname][animnames[i]] = an
				}
			}
		}
	}

	anims[cfg.N()] = a
}

func Get(name string) *Anims {
	return anims[name]
}

func GetGroup(name string) *Anims {
	return anims[name]
}

func New(rect pixel.Rect, margin [4]float64) *Anims {
	return &Anims{
		rect:   rect,
		items:  make(map[string]*Anim),
		groups: make(map[string]map[string]*Anim),
		sprite: pixel.NewSprite(nil, pixel.Rect{}),
		m:      margin,
	}
}

func (a *Anims) GetGroupSprite(group, name string, num int) (pixel.Picture, pixel.Rect) {
	_, ok := a.groups[group]
	if !ok {
		return a.GetSprite("idle", num)
	}

	idx := num % len(a.groups[group][name].frames)
	return a.groups[group][name].sheet, a.groups[group][name].frames[idx]
}

func (a *Anims) GetSprite(name string, num int) (pixel.Picture, pixel.Rect) {
	_, ok := a.items[name]
	if !ok {
		name = "idle" // fallback animation
		num = 1
	}

	idx := num % len(a.items[name].frames)
	return a.items[name].sheet, a.items[name].frames[idx]
}

func (a *Anims) GetLen(name string) int {
	return a.items[name].GetLen()
}

func (a *Anims) GetGroupLen(name string) int {
	return len(a.groups[name])
}

func (a *Anims) PrepareAnim(name, file string, frames []int) (*Anim, error) {
	spritesheet, err := loader.LoadPicture(file)
	if err != nil {
		return nil, err
	}

	frs := make([]pixel.Rect, 0, frames[0])

	frameWidth := a.rect.W()
	frameHeight := a.rect.H()
	for x := 0.0; x+frameWidth <= spritesheet.Bounds().Max.X; x += frameWidth {
		x += a.m[0]
		frs = append(frs, pixel.R(
			x,
			a.m[3],
			x+frameWidth,
			a.m[3]+frameHeight,
		))
		x += a.m[2]
	}

	return &Anim{
		sheet:  spritesheet,
		frames: frs[frames[1]:frames[2]],
		sprite: pixel.NewSprite(nil, pixel.Rect{}),
	}, nil
}

func (a *Anim) GetFrames() []pixel.Rect {
	return a.frames
}

func (a *Anim) GetLen() int {
	return len(a.frames)
}

func (a *Anim) GetSprite(idx int) *pixel.Sprite {
	a.sprite.Set(a.sheet, a.frames[idx])
	return a.sprite
}
