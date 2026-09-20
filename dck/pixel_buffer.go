package fr010

import "unsafe"

const (
	opaqueAlpha = uint32(0xff000000)
	opaqueWhite = uint32(0xffffffff)
)

// newPixelBuffer returns two views of the same little-endian storage. The
// rasterizer writes packed ABGR uint32 values and Ebitengine reads RGBA bytes,
// avoiding a full-frame color conversion before every upload.
func newPixelBuffer(pixelCount int) ([]uint32, []byte) {
	pixels := make([]uint32, pixelCount)
	if len(pixels) == 0 {
		return pixels, nil
	}
	rgba := unsafe.Slice((*byte)(unsafe.Pointer(unsafe.SliceData(pixels))), pixelCount*4)
	return pixels, rgba
}

func packRGB(r, g, b int) uint32 {
	return opaqueAlpha | uint32(b)<<16 | uint32(g)<<8 | uint32(r)
}
