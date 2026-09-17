package fr010

import (
	"fmt"
	"io"
	"sync"

	"github.com/olivierh59500/ym-player/pkg/stsound"
)

const ymBufferSize = 4096

// YMPlayer adapts StSound's mono samples to the stereo PCM stream expected by
// Ebitengine. Its working buffer is reused for every Read.
type YMPlayer struct {
	player *stsound.StSound
	buffer []int16
	mutex  sync.Mutex
	loop   bool
}

func NewYMPlayer(data []byte, sampleRate int, loop bool) (*YMPlayer, error) {
	player := stsound.CreateWithRate(sampleRate)
	if err := player.LoadMemory(data); err != nil {
		player.Destroy()
		return nil, fmt.Errorf("load YM data: %w", err)
	}

	player.SetLoopMode(loop)
	player.SetLowpassFilter(true)

	return &YMPlayer{
		player: player,
		buffer: make([]int16, ymBufferSize),
		loop:   loop,
	}, nil
}

func (y *YMPlayer) Read(p []byte) (int, error) {
	y.mutex.Lock()
	defer y.mutex.Unlock()

	if y.player == nil {
		return 0, io.EOF
	}

	// Ebitengine consumes signed 16-bit little-endian stereo frames.
	frameCount := len(p) / 4
	processed := 0
	var readErr error
	for processed < frameCount {
		chunkSize := min(frameCount-processed, len(y.buffer))
		if !y.player.Compute(y.buffer[:chunkSize], chunkSize) && !y.loop {
			clear(p[processed*4 : frameCount*4])
			readErr = io.EOF
			break
		}

		for i, mono := range y.buffer[:chunkSize] {
			sample := mono / 2
			offset := (processed + i) * 4
			low := byte(sample)
			high := byte(sample >> 8)
			p[offset] = low
			p[offset+1] = high
			p[offset+2] = low
			p[offset+3] = high
		}
		processed += chunkSize
	}

	return frameCount * 4, readErr
}

func (y *YMPlayer) Close() error {
	if y == nil {
		return nil
	}

	y.mutex.Lock()
	defer y.mutex.Unlock()
	if y.player != nil {
		y.player.Destroy()
		y.player = nil
	}
	return nil
}
