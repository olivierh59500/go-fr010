package fr010

import (
	"github.com/olivierh59500/democonstructionkit/sound"
	"testing"
)

func BenchmarkMusicRead4096(b *testing.B) {
	player, err := sound.Open("soundtrack.ym", ymData, sound.Options{SampleRate: sampleRate, PCMFormat: sound.PCM16, Loop: true, Gain: .5})
	if err != nil {
		b.Fatal(err)
	}
	b.Cleanup(func() { _ = player.Close() })
	buffer := make([]byte, 4096*4)

	b.SetBytes(int64(len(buffer)))
	b.ReportAllocs()
	b.ResetTimer()
	for b.Loop() {
		if _, err := player.Read(buffer); err != nil {
			b.Fatal(err)
		}
	}
}

func BenchmarkRenderFrame(b *testing.B) {
	game, err := NewGame()
	if err != nil {
		b.Fatal(err)
	}
	frames := []int{10000, 30000, 50000, 75000, 100000, 130000, 170000, 180000}

	b.ReportAllocs()
	b.ResetTimer()
	for i := 0; b.Loop(); i++ {
		game.renderFrame(frames[i%len(frames)])
	}
}

func BenchmarkRenderFrameByScene(b *testing.B) {
	frames := []struct {
		name  string
		frame int
	}{
		{name: "intro", frame: 10000},
		{name: "scene1", frame: 30000},
		{name: "scene2", frame: 50000},
		{name: "scene3", frame: 75000},
		{name: "scene4", frame: 100000},
		{name: "greetings", frame: 130000},
		{name: "scene6", frame: 170000},
		{name: "credits", frame: 180000},
	}

	for _, frame := range frames {
		b.Run(frame.name, func(b *testing.B) {
			game, err := NewGame()
			if err != nil {
				b.Fatal(err)
			}
			b.ReportAllocs()
			b.ResetTimer()
			for b.Loop() {
				game.renderFrame(frame.frame)
			}
		})
	}
}
