package fr010

import (
	"encoding/binary"
	"math"
)

type binReader struct {
	data []byte
	pos  int
}

func newBinReader(data []byte) *binReader {
	return &binReader{data: data}
}

func (r *binReader) readU8() byte {
	b := r.data[r.pos]
	r.pos++
	return b
}

func (r *binReader) readI16() int16 {
	v := int16(binary.LittleEndian.Uint16(r.data[r.pos:]))
	r.pos += 2
	return v
}

func (r *binReader) readU16() uint16 {
	v := binary.LittleEndian.Uint16(r.data[r.pos:])
	r.pos += 2
	return v
}

func (r *binReader) readF32() float32 {
	bits := binary.LittleEndian.Uint32(r.data[r.pos:])
	r.pos += 4
	return math.Float32frombits(bits)
}

func (r *binReader) readF32Slice(count int) []float32 {
	out := make([]float32, count)
	for i := 0; i < count; i++ {
		out[i] = r.readF32()
	}
	return out
}
