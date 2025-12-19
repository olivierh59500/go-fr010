package main

import "time"

type TextureManager struct {
	Textures []byte
	blendmap []byte
	perlin   *Perlin
}

func NewTextureManager(number int) *TextureManager {
	tm := &TextureManager{
		Textures: make([]byte, number*256*256),
		blendmap: make([]byte, 256*256),
		perlin:   NewPerlin(),
	}

	o := 0
	for y := -128; y < 128; y++ {
		for x := -128; x < 128; x++ {
			r1 := float32(absMax(x, y)) / 128.0
			r2 := float32(x*x+y*y) / 16000.0
			r := r1 * r2
			k := float32(1.0 - r*r*r*r)
			if k < 0 {
				k = 0
			}
			if k > 1 {
				k = 1
			}
			tm.blendmap[o] = byte(k * 255)
			o++
		}
	}

	for i := 0; i < number; i++ {
		ot := i * 256 * 256
		t := int(time.Now().UnixMilli())
		for v := 0; v < 256; v++ {
			for u := 0; u < 256; u++ {
				om := u + (v << 8)
				n := int(tm.perlin.GetI3DI(uint32((u+t*2)<<8), uint32((v+t)<<8), uint32(t))) * 8
				b := int(tm.blendmap[om]) * n
				b >>= 11
				if b > 0xff {
					b = 0xff
				}
				tm.Textures[ot+u+(v<<8)] = byte(b)
			}
		}
	}

	return tm
}

func (tm *TextureManager) GetTexture(number int) []byte {
	start := number * 256 * 256
	end := start + 256*256
	return tm.Textures[start:end]
}
