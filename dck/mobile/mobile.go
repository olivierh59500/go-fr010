// Package mobile exposes FR-010 to ebitenmobile.
package mobile

import (
	enginemobile "github.com/hajimehoshi/ebiten/v2/mobile"
	"github.com/olivierh59500/fr010/dck"
)

func init() {
	game, err := fr010.NewGame()
	if err != nil {
		panic(err)
	}
	enginemobile.SetGame(game)
}

// Dummy forces gomobile to include this package in the Android binding.
func Dummy() {}
