package config

type Main struct {
	TestFlag  bool   `json:"TestFlag"`
	Assets    string `json:"Assets"`
	Width     int    `json:"Width"`
	Height    int    `json:"Height"`
	PlayerCfg string `json:"PlayerCfg"`
	WorldCfg  string `json:"WorldCfg"`
}

type Player struct {
	Anims []struct {
		Name   string `json:"name"`
		File   string `json:"file"`
		Frames []int  `json:"frames"`
	} `json:"anims"`
	Width  int `json:"width"`
	Height int `json:"height"`
	Mass   int `json:"mass"`
	Walk   int `json:"walk"`
	Run    int `json:"run"`
}

type World struct {
	Platforms [][]int `json:"platforms"`
	Width     int     `json:"width"`
	Heigth    int     `json:"heigth"`
	Gravity   int     `json:"gravity"`
}
