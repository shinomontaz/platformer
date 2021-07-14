package config

type Main struct {
	TestFlag  bool    `json:"TestFlag"`
	Assets    string  `json:"Assets"`
	Width     float64 `json:"Width"`
	Height    float64 `json:"Height"`
	PlayerCfg string  `json:"PlayerCfg"`
	WorldCfg  string  `json:"WorldCfg"`
}

type Player struct {
	Anims []struct {
		Name   string `json:"name"`
		File   string `json:"file"`
		Frames []int  `json:"frames"`
	} `json:"anims"`
	Width  float64 `json:"width"`
	Height float64 `json:"height"`
	Mass   int     `json:"mass"`
	Walk   float64 `json:"walk"`
	Run    float64 `json:"run"`
}

type World struct {
	Platforms [][]float64 `json:"platforms"`
	Width     float64     `json:"width"`
	Heigth    float64     `json:"heigth"`
	Gravity   float64     `json:"gravity"`
}
