package main

import (
	"errors"
	"testing"

	"github.com/hajimehoshi/ebiten/v2"
)

type stubGame struct {
	draws int
	err   error
}

func (g *stubGame) Update() error {
	return g.err
}

func (g *stubGame) Draw(*ebiten.Image) {
	g.draws++
}

func (*stubGame) Layout(int, int) (int, int) {
	return 640, 480
}

func TestDrawOnUpdateGame(t *testing.T) {
	inner := &stubGame{}
	game := newDrawOnUpdateGame(inner)

	game.Draw(nil)
	game.Draw(nil)
	if inner.draws != 1 {
		t.Fatalf("initial draws = %d, want 1", inner.draws)
	}

	if err := game.Update(); err != nil {
		t.Fatal(err)
	}
	game.Draw(nil)
	game.Draw(nil)
	if inner.draws != 2 {
		t.Fatalf("draws after update = %d, want 2", inner.draws)
	}

	inner.err = errors.New("stop")
	if err := game.Update(); !errors.Is(err, inner.err) {
		t.Fatalf("Update error = %v, want %v", err, inner.err)
	}
	game.Draw(nil)
	if inner.draws != 2 {
		t.Fatalf("draws after failed update = %d, want 2", inner.draws)
	}
}
