package config

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"os"
)

var (
	MainConfig   Main
	PlayerConfig Player
	WorldConfig  World
)

func init() {
	var byteValue []byte
	fconfig, err := os.Open("config/config.json")
	if err != nil {
		log.Fatal(err)
	}
	defer fconfig.Close()

	byteValue, _ = ioutil.ReadAll(fconfig)
	json.Unmarshal(byteValue, &MainConfig)

	fworld, err := os.Open(fmt.Sprintf("config/%s", MainConfig.WorldCfg))
	if err != nil {
		log.Fatal(err)
	}
	defer fworld.Close()

	byteValue, _ = ioutil.ReadAll(fworld)
	json.Unmarshal(byteValue, &WorldConfig)

	fenemy, err := os.Open("config/player2.json")
	if err != nil {
		log.Fatal(err)
	}
	defer fenemy.Close()

	// install enemies config
	var EnemyConfig Player
	byteValue, _ = ioutil.ReadAll(fenemy)
	json.Unmarshal(byteValue, &EnemyConfig)
	for _, e := range WorldConfig.Enemies {
		e.Anims = &EnemyConfig.Anims
	}
	fplayer, err := os.Open(fmt.Sprintf("config/%s", MainConfig.PlayerCfg))
	if err != nil {
		log.Fatal(err)
	}
	defer fplayer.Close()

	byteValue, _ = ioutil.ReadAll(fplayer)
	json.Unmarshal(byteValue, &PlayerConfig)
}
