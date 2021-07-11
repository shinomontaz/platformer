package main

import (
	"fmt"

	"platformer/config"

	"github.com/faiface/pixel/pixelgl"
)

func init() {
	// create world
	// create physics
	// create hero
}

func gameLoop() {

}

func main() {
	fmt.Println(config.PlayerConfig)
	pixelgl.Run(gameLoop)
}
