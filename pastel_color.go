package main

type PastelColor struct {
	ColorMap [256]uint32
}

func NewPastelColor(v Vector) *PastelColor {
	pc := &PastelColor{}
	for i := 0; i < 256; i++ {
		r := 255 - int(float32(i)*v.X)
		g := 255 - int(float32(i)*v.Y)
		b := 255 - int(float32(i)*v.Z)
		if r < 0 {
			r = 0
		}
		if g < 0 {
			g = 0
		}
		if b < 0 {
			b = 0
		}
		pc.ColorMap[i] = uint32((r << 16) | (g << 8) | b)
	}
	return pc
}
