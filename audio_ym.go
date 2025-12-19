package main

import (
	"sync"

	"github.com/olivierh59500/ym-player/pkg/audio"
	"github.com/olivierh59500/ym-player/pkg/stsound"
)

const (
	ymSampleRate = 44100
	ymBufferSize = 2048
)

type YMPlayer struct {
	st        *stsound.StSound
	output    audio.Output
	buffer    []int16
	stop      chan struct{}
	done      chan struct{}
	closeOnce sync.Once
}

func NewYMPlayer(data []byte) (*YMPlayer, error) {
	st := stsound.CreateWithRate(ymSampleRate)
	if err := st.LoadMemory(data); err != nil {
		st.Destroy()
		return nil, err
	}
	st.SetLoopMode(false)
	st.SetLowpassFilter(true)

	var out audio.Output
	out, err := audio.NewStreamingOtoOutput()
	if err != nil {
		out, err = audio.NewFallbackOutput()
		if err != nil {
			st.Destroy()
			return nil, err
		}
	}
	if err := out.Open(ymSampleRate, 1, ymBufferSize); err != nil {
		st.Destroy()
		return nil, err
	}

	p := &YMPlayer{
		st:     st,
		output: out,
		buffer: make([]int16, ymBufferSize),
		stop:   make(chan struct{}),
		done:   make(chan struct{}),
	}

	go p.loop()
	return p, nil
}

func (p *YMPlayer) loop() {
	defer close(p.done)
	p.st.Play()
	for {
		select {
		case <-p.stop:
			return
		default:
		}

		if !p.st.Compute(p.buffer, len(p.buffer)) {
			return
		}

		_ = p.output.Write(p.buffer)
	}
}

func (p *YMPlayer) PosMS() int {
	if p == nil || p.st == nil {
		return 0
	}
	return int(p.st.GetPos())
}

func (p *YMPlayer) Close() {
	if p == nil {
		return
	}
	p.closeOnce.Do(func() {
		close(p.stop)
		<-p.done
		_ = p.output.Close()
		p.st.Destroy()
	})
}
