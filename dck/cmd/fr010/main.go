package main

import (
	"log"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/olivierh59500/fr010/dck"
)

const (
	windowWidth  = 640
	windowHeight = 480
)

func main() {
	game, err := fr010.NewGame()
	if err != nil {
		log.Fatal(err)
	}
	defer game.Close()

	ebiten.SetWindowSize(windowWidth, windowHeight)
	ebiten.SetWindowResizingMode(ebiten.WindowResizingModeEnabled)
	ebiten.SetWindowTitle("FR-010 (Go/Ebitengine)")
	ebiten.SetScreenClearedEveryFrame(false)
	if err := ebiten.RunGame(newDrawOnUpdateGame(game)); err != nil && err != ebiten.Termination {
		log.Fatal(err)
	}
}
