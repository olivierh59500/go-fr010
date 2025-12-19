package main

import "math"

const (
	perlinYWrap  = 16
	perlinYWrapB = 4
	perlinTWrap  = 256
	perlinTWrapB = 8
	perlinOctave = 2
	perlinIBuf   = 1 << 16
)

func fsc(i float32) float32 {
	return 0.5 * (1 - float32(math.Cos(float64(i*3.1415926535))))
}

type Perlin struct {
	buffer [4096]float32
	ucbuf  [perlinIBuf]byte
	wtb    [256]byte
}

func NewPerlin() *Perlin {
	p := &Perlin{}
	for i := 0; i < len(p.buffer); i++ {
		p.buffer[i] = randFloat()
	}
	for i := 0; i < len(p.ucbuf); i++ {
		p.ucbuf[i] = randByte()
	}
	for i := 0; i < 256; i++ {
		p.wtb[i] = byte(255 * fsc(float32(i)/256.0))
	}
	return p
}

func (p *Perlin) Get(x float32) float32 {
	if x < 0 {
		x = -x
	}
	var r float32
	ampl := float32(0.5)
	xi := int(x)
	xf := x - float32(xi)

	for i := 0; i < perlinOctave; i++ {
		n1 := p.buffer[xi&4095]
		n1 += fsc(xf) * (p.buffer[(xi+1)&4095] - n1)
		r += n1 * ampl
		ampl *= 0.5
		xi <<= 1
		xf *= 2
		if xf >= 1.0 {
			xi++
			xf -= 1.0
		}
	}
	return r
}

func (p *Perlin) Get2D(x, y float32) float32 {
	if x < 0 {
		x = -x
	}
	if y < 0 {
		y = -y
	}
	xi := int(x)
	yi := int(y)
	var r float32
	ampl := float32(0.5)
	xf := x - float32(xi)
	yf := y - float32(yi)

	xf, yf = 0, 0

	for i := 0; i < perlinOctave; i++ {
		of := xi + (yi << perlinYWrapB)
		rxf := fsc(xf)

		n1 := p.buffer[of&4095]
		n1 += rxf * (p.buffer[(of+1)&4095] - n1)
		n2 := p.buffer[(of+perlinYWrap)&4095]
		n2 += rxf * (p.buffer[(of+perlinYWrap+1)&4095] - n2)
		n1 += fsc(yf) * (n2 - n1)

		r += n1 * ampl
		ampl *= 0.5
		xf *= 2
		xi <<= 1
		yf *= 2
		yi <<= 1
		if xf >= 1.0 {
			xi++
			xf -= 1.0
		}
		if yf >= 1.0 {
			yi++
			yf -= 1.0
		}
	}
	return r
}

func (p *Perlin) Get3D(x, y, t float32) float32 {
	if x < 0 {
		x = -x
	}
	if y < 0 {
		y = -y
	}
	if t < 0 {
		t = -t
	}

	xi := int(x)
	yi := int(y)
	ti := int(t)
	xf := x - float32(xi)
	yf := y - float32(yi)
	tf := t - float32(ti)
	var r float32
	ampl := float32(0.5)

	for i := 0; i < perlinOctave; i++ {
		of := xi + (yi << perlinYWrapB) + (ti << perlinTWrapB)
		rxf := fsc(xf)
		ryf := fsc(yf)

		n1 := p.buffer[of&4095]
		n1 += rxf * (p.buffer[(of+1)&4095] - n1)
		n2 := p.buffer[(of+perlinYWrap)&4095]
		n2 += rxf * (p.buffer[(of+perlinYWrap+1)&4095] - n2)
		n1 += ryf * (n2 - n1)

		of += perlinTWrap
		n2 = p.buffer[of&4095]
		n2 += rxf * (p.buffer[(of+1)&4095] - n2)
		n3 := p.buffer[(of+perlinYWrap)&4095]
		n3 += rxf * (p.buffer[(of+perlinYWrap+1)&4095] - n3)
		n2 += ryf * (n3 - n2)

		n1 += fsc(tf) * (n2 - n1)

		r += n1 * ampl
		ampl *= 0.5
		xi <<= 1
		xf *= 2
		yi <<= 1
		yf *= 2
		ti <<= 1
		tf *= 2
		if xf >= 1.0 {
			xi++
			xf -= 1.0
		}
		if yf >= 1.0 {
			yi++
			yf -= 1.0
		}
		if tf >= 1.0 {
			ti++
			tf -= 1.0
		}
	}

	return r
}

func (p *Perlin) GetI3DI(xfp, yfp, tfp uint32) byte {
	r := 0
	amplShift := 17
	for amplShift < 25 {
		xindex := (xfp >> 16) & 0xffff
		yindex := ((yfp >> 16) & 0xffff) << 5
		tindex := ((tfp >> 16) & 0xffff) << 10
		xf := int(p.wtb[(xfp>>8)&0xff])
		yf := int(p.wtb[(yfp>>8)&0xff])
		tf := int(p.wtb[(tfp>>8)&0xff])
		off := int(xindex + yindex + tindex)

		n1 := int(p.ucbuf[off&0xffff])
		n2 := int(p.ucbuf[(off+1)&0xffff])
		nh1 := (n1 << 8) + (n2-n1)*xf
		n1 = int(p.ucbuf[(off+32)&0xffff])
		n2 = int(p.ucbuf[(off+33)&0xffff])
		nh2 := (n1 << 8) + (n2-n1)*xf
		nv1 := ((nh1 << 8) + (nh2-nh1)*yf) >> 8

		n1 = int(p.ucbuf[(off+1024)&0xffff])
		n2 = int(p.ucbuf[(off+1025)&0xffff])
		nh1 = (n1 << 8) + (n2-n1)*xf
		n1 = int(p.ucbuf[(off+1056)&0xffff])
		n2 = int(p.ucbuf[(off+1057)&0xffff])
		nh2 = (n1 << 8) + (n2-n1)*xf
		nv2 := ((nh1 << 8) + (nh2-nh1)*yf) >> 8

		n := ((nv1 << 8) + (nv2-nv1)*tf) >> amplShift
		r += n
		xfp <<= 1
		yfp <<= 1
		tfp <<= 1
		amplShift++
	}

	return byte(r)
}

func (p *Perlin) GetI3D(x, y, t float32) byte {
	xfp := uint32(x * float32((1<<16)-1))
	yfp := uint32(y * float32((1<<16)-1))
	tfp := uint32(t * float32((1<<16)-1))
	return p.GetI3DI(xfp, yfp, tfp)
}
