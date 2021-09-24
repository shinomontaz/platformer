package config

type Main struct {
	TestFlag  bool    `json:"TestFlag"`
	Assets    string  `json:"Assets"`
	Width     float64 `json:"Width"`
	Height    float64 `json:"Height"`
	PlayerCfg string  `json:"PlayerCfg"`
	WorldCfg  string  `json:"WorldCfg"`
}

type Anim struct {
	Name   string `json:"name"`
	File   string `json:"file"`
	Frames []int  `json:"frames"`
}

type Player struct {
	Anims  []Anim  `json:"anims"`
	Width  float64 `json:"width"`
	Height float64 `json:"height"`
	Mass   int     `json:"mass"`
	Walk   float64 `json:"walk"`
	Run    float64 `json:"run"`
}

type Enemy struct {
	Coords [2]float64 `json:"coords"`
	Width  float64    `json:"width"`
	Height float64    `json:"height"`
	Mass   int        `json:"mass"`
	Walk   float64    `json:"walk"`
	Run    float64    `json:"run"`
	Type   int        `json:"type"`
	Anims  *[]Anim
}

type World struct {
	Platforms [][]float64 `json:"platforms"`
	Enemies   []*Enemy    `json:"enemies"`
	Width     float64     `json:"width"`
	Heigth    float64     `json:"heigth"`
	Gravity   float64     `json:"gravity"`
}
