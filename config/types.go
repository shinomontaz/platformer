package config

type Main struct {
	TestFlag bool    `json:"TestFlag"`
	Assets   string  `json:"Assets"`
	Width    float64 `json:"Width"`
	Height   float64 `json:"Height"`
	WorldCfg string  `json:"WorldCfg"`
	AllAnims string  `json:"AllAnims"`
	Profiles string  `json:"Profiles"`
	Loots    string  `json:"Loots"`
	Sounds   string  `json:"sounds"`
	Spells   string  `json:"spells"`
}

type Anim struct {
	Name   string `json:"name"`
	File   string `json:"file"`
	Frames []int  `json:"frames"`
}

type Player struct {
	Anim   string  `json:"anim"`
	Width  float64 `json:"width"`
	Height float64 `json:"height"`
	Margin float64 `json:"margin"`
	Mass   int     `json:"mass"`
	Walk   float64 `json:"walk"`
	Run    float64 `json:"run"`
}

type AnGroup struct {
	Name string `json:"name"`
	List []Anim `json:"group"`
}

type Anims struct {
	Name   string    `json:"type"`
	Width  int       `json:"width"`
	Height int       `json:"height"`
	Margin []int     `json:"margin"`
	List   []Anim    `json:"list"`
	Groups []AnGroup `json:"groups"`
}

type Skill struct {
	Type   string    `json:"type"`
	Name   string    `json:"name"`
	Min    float64   `json:"min"`
	Max    float64   `json:"max"`
	Weight int       `json:"weight"`
	Ttl    float64   `json:"ttl"`
	Hitbox []float64 `json:"hitbox"`
}
type Profile struct {
	Type     string  `json:"type"`
	Anim     string  `json:"anim"`
	Portrait string  `json:"portrait"`
	Body     string  `json:"body"`
	Dir      float64 `json:"dir"`
	Hp       int     `json:"hp"`
	Strength int     `json:"strength"`
	Width    float64 `json:"width"`
	Height   float64 `json:"height"`
	Margin   float64 `json:"margin"`
	Jump     float64 `json:"jump"`
	Mass     float64 `json:"mass"`
	Walk     float64 `json:"walk"`
	Run      float64 `json:"run"`
	Skills   []Skill `json:"skills"`
	Phrases  string  `json:"phrases"`
	Dialog   int     `json:"dialog"`
	Only     int     `json:"only"`
}

type Soundeffect struct {
	Type string   `json:"type"`
	List []string `json:"list"`
}
type Soundprofile struct {
	Type string        `json:"type"`
	List []Soundeffect `json:"list"`
}

type Spellprofile struct {
	Type   string    `json:"type"`
	Anim   string    `json:"anim"`
	Sound  string    `json:"sound"`
	Hitbox []float64 `json:"hitbox"`
}

type Options struct {
	Vsync        bool               `json:"Vsync"`
	Fullscreen   bool               `json:"Fullscreen"`
	WindowHeight float64            `json:"WindowHeight"`
	WindowWidth  float64            `json:"WindowWidth"`
	CurrentMap   int                `json:"CurrentMap"`
	Volumes      map[string]float64 `json:"Volumes"`
}
