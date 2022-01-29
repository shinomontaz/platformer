package config

type Main struct {
	TestFlag  bool    `json:"TestFlag"`
	Assets    string  `json:"Assets"`
	Width     float64 `json:"Width"`
	Height    float64 `json:"Height"`
	PlayerCfg string  `json:"PlayerCfg"`
	WorldCfg  string  `json:"WorldCfg"`
	AllAnims  string  `json:"AllAnims"`
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
	Margin int       `json:"margin"`
	List   []Anim    `json:"list"`
	Groups []AnGroup `json:"groups"`
}
