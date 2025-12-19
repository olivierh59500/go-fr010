package main

import (
	"log"

	"github.com/hajimehoshi/ebiten/v2"
)

const (
	screenWidth  = 640
	screenHeight = 480
)

func main() {
	game, err := NewGame()
	if err != nil {
		log.Fatal(err)
	}
	defer game.Close()

	ebiten.SetWindowSize(screenWidth, screenHeight)
	ebiten.SetWindowTitle("FR-010 (Go/Ebiten)")
	if err := ebiten.RunGame(game); err != nil && err != ebiten.Termination {
		log.Fatal(err)
	}
}
