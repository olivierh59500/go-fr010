package main

import "math"

const lineDarken = 0x60606060

var textLine bool

type Drawing struct {
	t        float32
	perlin   *Perlin
	vscreen  []uint32
	width    int
	height   int
	pitch    int
	textures *TextureManager

	o_eck int
	u_eck int
	r_eck int
	l_eck int

	rdy int
	ldy int

	rx int
	lx int
	lu int
	lv int

	tdu int
	tdv int

	rdx int
	ldx int
	ldu int
	ldv int

	poly struct {
		x [3]int
		y [3]int
		u [3]int
		v [3]int
	}
}

func NewDrawing(width, height int, textures *TextureManager) *Drawing {
	return &Drawing{
		perlin:   NewPerlin(),
		width:    width,
		height:   height,
		pitch:    width,
		textures: textures,
	}
}

func (d *Drawing) SetBuffer(buf []uint32, pitch int) {
	d.vscreen = buf
	d.pitch = pitch
}

func (d *Drawing) DrawLine(a, b Vector) {
	var i, j int
	delta := a.Sub(b)
	length := float32(math.Sqrt(float64(delta.X*delta.X + delta.Y*delta.Y)))

	subdivs := int(length/20) + 1

	maxd := float32(16.0)
	if length < 32.0 {
		maxd = length * 0.5
	}
	tf := float32(1.0)
	if textLine {
		maxd *= 2.1
		tf *= 5.5
	} else {
		maxd *= 1.6
		tf *= 2.2
	}

	b.X += (d.perlin.Get3D(b.X*0.01, b.Y*0.01, d.t*1.3*tf)-0.5) * 0.5 * maxd
	b.Y += (d.perlin.Get3D(b.Y*0.01, b.X*0.01, d.t*1.4*tf)-0.5) * 0.5 * maxd
	a.X += (d.perlin.Get3D(a.X*0.01, a.Y*0.01, d.t*1.3*tf)-0.5) * 0.5 * maxd
	a.Y += (d.perlin.Get3D(a.Y*0.01, a.X*0.01, d.t*1.4*tf)-0.5) * 0.5 * maxd
	a.X += (d.perlin.Get2D(a.X*0.01, a.Y*0.01)-0.5) * maxd
	a.Y += (d.perlin.Get2D(a.Y*0.01, a.X*0.01)-0.5) * maxd

	for i = 0; i < 4; i++ {
		delta = a.Sub(b)
		a.AddInPlace(delta.MulScalar((d.perlin.Get2D(b.X*0.0404, b.Y*0.00707) - 0.6) * 0.1))
		b.SubInPlace(delta.MulScalar((d.perlin.Get2D(b.X*0.0303, b.Y*0.00909) - 0.6) * 0.1))

		k := a
		step := b.Sub(a).DivScalar(float32(subdivs))

		for j = 0; j < (subdivs+1)/2; j++ {
			ok := k
			k.AddInPlace(step)
			shift := (d.perlin.Get(float32(j)*0.1) - 0.5) * 0.303
			k.X += shift * -step.Y
			k.Y += shift * step.X
			d.DDBresenhamLine(ok.X, ok.Y, k.X, k.Y, lineDarken)
		}
		for j = 0; j < subdivs/2; j++ {
			ok := k
			k.AddInPlace(step)
			shift := (d.perlin.Get(float32(j)*0.1) - 0.5) * 0.303
			k.X += shift * step.Y
			k.Y += shift * -step.X
			d.DDBresenhamLine(ok.X, ok.Y, k.X, k.Y, lineDarken)
		}

		b.X += (d.perlin.Get3D(b.X*0.01, b.Y*0.01, float32(i)*0.1*tf)-0.5) * step.X * 0.01707
		b.Y += (d.perlin.Get3D(b.Y*0.01, b.X*0.01, float32(i)*0.1*tf)-0.5) * step.Y * 0.01707
	}
}

func (d *Drawing) drawFace(f *DrawFaceObj) {
	if d.textures == nil || d.vscreen == nil {
		return
	}

	col := f.col.ColorMap[:]
	tex := d.textures.GetTexture(f.texNum)

	d.poly.x[0] = int(f.x1 * 65536.0)
	d.poly.y[0] = int(f.y1 * 65536.0)
	d.poly.u[0] = int(f.u1 * 65536.0)
	d.poly.v[0] = int(f.v1 * 65536.0)

	d.poly.x[1] = int(f.x2 * 65536.0)
	d.poly.y[1] = int(f.y2 * 65536.0)
	d.poly.u[1] = int(f.u2 * 65536.0)
	d.poly.v[1] = int(f.v2 * 65536.0)

	d.poly.x[2] = int(f.x3 * 65536.0)
	d.poly.y[2] = int(f.y3 * 65536.0)
	d.poly.u[2] = int(f.u3 * 65536.0)
	d.poly.v[2] = int(f.v3 * 65536.0)

	miny := d.poly.y[0]
	maxy := d.poly.y[0]
	d.o_eck = 0
	d.u_eck = 0

	for i := 1; i < 3; i++ {
		if d.poly.y[i] < miny {
			miny = d.poly.y[i]
			d.o_eck = i
		}
		if d.poly.y[i] > maxy {
			maxy = d.poly.y[i]
			d.u_eck = i
		}
	}

	d.l_eck = d.o_eck
	d.r_eck = d.o_eck

	det := (float64(d.poly.x[0]-d.poly.x[2])*float64(d.poly.y[1]-d.poly.y[2]) -
		float64(d.poly.x[1]-d.poly.x[2])*float64(d.poly.y[0]-d.poly.y[2])) / (65536.0 * 65536.0)
	if det == 0.0 {
		return
	}

	inv := 1.0 / det
	y02 := float64(d.poly.y[0]-d.poly.y[2]) / 65536.0
	y12 := float64(d.poly.y[1]-d.poly.y[2]) / 65536.0

	d.tdu = int(((float64(d.poly.u[0]-d.poly.u[2])*y12 - float64(d.poly.u[1]-d.poly.u[2])*y02) * inv))
	d.tdv = int(((float64(d.poly.v[0]-d.poly.v[2])*y12 - float64(d.poly.v[1]-d.poly.v[2])*y02) * inv))

	for {
		if d.r_eck == d.u_eck {
			return
		}
		d.rightSide()
		if d.rdy > 0 {
			break
		}
	}

	for {
		if d.l_eck == d.u_eck {
			return
		}
		d.leftSide()
		if d.ldy > 0 {
			break
		}
	}

	vbuffer := ceil16r(miny) * d.pitch

	for {
		x1 := ceil16r(d.lx)
		w := ceil16r(d.rx) - x1

		if w > 0 {
			prestep := (x1 << 16) - d.lx
			u := d.lu + imul16(prestep, d.tdu)
			v := d.lv + imul16(prestep, d.tdv)
			end := x1 + w

			for x := x1; x < end; x++ {
				tu := uint32(u)
				tv := uint32(v)
				texIndex := int(((tu >> 16) & 0x00ff) | ((tv >> 8) & 0xff00))
				c := col[tex[texIndex]]
				idx := vbuffer + x
				if idx >= 0 && idx < len(d.vscreen) {
					d.vscreen[idx] = c
				}
				u += d.tdu
				v += d.tdv
			}
		}

		vbuffer += d.pitch

		d.rdy--
		if d.rdy <= 0 {
			for d.rdy <= 0 {
				if d.r_eck == d.u_eck {
					return
				}
				d.rightSide()
			}
		} else {
			d.rx += d.rdx
		}

		d.ldy--
		if d.ldy <= 0 {
			for d.ldy <= 0 {
				if d.l_eck == d.u_eck {
					return
				}
				d.leftSide()
			}
		} else {
			d.lx += d.ldx
			d.lu += d.ldu
			d.lv += d.ldv
		}
	}
}

func (d *Drawing) rightSide() {
	index := d.r_eck - 1
	if index < 0 {
		index = 2
	}

	d.rdy = ceil16r(d.poly.y[index]) - ceil16r(d.poly.y[d.r_eck])
	if d.rdy <= 0 {
		d.r_eck = index
		return
	}

	dy := d.poly.y[index] - d.poly.y[d.r_eck]
	if d.rdy > 1 {
		d.rdx = idiv16(d.poly.x[index]-d.poly.x[d.r_eck], dy)
	} else {
		idy := (0x10000 << 14) / dy
		d.rdx = imul14(d.poly.x[index]-d.poly.x[d.r_eck], idy)
	}

	prestep := (ceil16r(d.poly.y[d.r_eck]) << 16) - d.poly.y[d.r_eck]
	d.rx = d.poly.x[d.r_eck] + imul16(prestep, d.rdx)
	d.r_eck = index
}

func (d *Drawing) leftSide() {
	index := d.l_eck + 1
	if index > 2 {
		index = 0
	}

	d.ldy = ceil16r(d.poly.y[index]) - ceil16r(d.poly.y[d.l_eck])
	if d.ldy <= 0 {
		d.l_eck = index
		return
	}

	dy := d.poly.y[index] - d.poly.y[d.l_eck]
	if d.ldy > 1 {
		d.ldx = idiv16(d.poly.x[index]-d.poly.x[d.l_eck], dy)
		d.ldu = idiv16(d.poly.u[index]-d.poly.u[d.l_eck], dy)
		d.ldv = idiv16(d.poly.v[index]-d.poly.v[d.l_eck], dy)
	} else {
		idy := (0x10000 << 14) / dy
		d.ldx = imul14(d.poly.x[index]-d.poly.x[d.l_eck], idy)
		d.ldu = imul14(d.poly.u[index]-d.poly.u[d.l_eck], idy)
		d.ldv = imul14(d.poly.v[index]-d.poly.v[d.l_eck], idy)
	}

	prestep := (ceil16r(d.poly.y[d.l_eck]) << 16) - d.poly.y[d.l_eck]
	d.lx = d.poly.x[d.l_eck] + imul16(prestep, d.ldx)
	d.lu = d.poly.u[d.l_eck] + imul16(prestep, d.ldu)
	d.lv = d.poly.v[d.l_eck] + imul16(prestep, d.ldv)
	d.l_eck = index
}

func (d *Drawing) DDBresenhamLine(x0, y0, x1, y1 float32, value uint32) {
	if x0 < 1 {
		x0 = 1
	}
	if x0 > float32(d.width-2) {
		x0 = float32(d.width - 2)
	}
	if x1 < 1 {
		x1 = 1
	}
	if x1 > float32(d.width-2) {
		x1 = float32(d.width - 2)
	}
	if y0 < 1 {
		y0 = 1
	}
	if y0 > float32(d.height-2) {
		y0 = float32(d.height - 2)
	}
	if y1 < 1 {
		y1 = 1
	}
	if y1 > float32(d.height-2) {
		y1 = float32(d.height - 2)
	}

	if x0 == x1 && y0 == y1 {
		return
	}

	if math.Abs(float64(x1-x0)) > math.Abs(float64(y1-y0)) {
		if y0 < y1 {
			x0, x1 = x1, x0
			y0, y1 = y1, y0
		}
	} else {
		if x0 < x1 {
			x0, x1 = x1, x0
			y0, y1 = y1, y0
		}
	}

	x := int(x0)
	y := int(y0)
	dx := x1 - x0
	dy := y1 - y0

	incre := 1
	if dx < 0 {
		incre = -1
		dx = -dx
	}

	incrne := d.pitch
	if dy < 0 {
		incrne = -d.pitch
		dy = -dy
	}

	dxi := int(dx * 256.0)
	dyi := int(dy * 256.0)

	if dx < dy {
		offset := x + y*d.pitch
		length := int(dy) + 1
		error := int(float32(y0-float32(y))*float32(dxi) - float32(x0-float32(x))*float32(dyi))
		for i := 0; i < length; i++ {
			offset += incrne
			error += dxi
			if error > dyi {
				offset += incre
				error -= dyi
			}
			d.darkenPixel(offset, value)
		}
	} else {
		offset := x + y*d.pitch
		length := int(dx) + 1
		error := int(float32(x0-float32(x))*float32(dyi) - float32(y0-float32(y))*float32(dxi))
		for i := 0; i < length; i++ {
			offset += incre
			error += dyi
			if error > dxi {
				offset += incrne
				error -= dxi
			}
			d.darkenPixel(offset, value)
		}
	}
}

func (d *Drawing) darkenPixel(offset int, value uint32) {
	if offset < 0 || offset >= len(d.vscreen) {
		return
	}
	p := d.vscreen[offset]
	r := int((p>>16)&0xff) - int((value>>16)&0xff)
	g := int((p>>8)&0xff) - int((value>>8)&0xff)
	b := int(p&0xff) - int(value&0xff)
	if r < 0 {
		r = 0
	}
	if g < 0 {
		g = 0
	}
	if b < 0 {
		b = 0
	}
	d.vscreen[offset] = uint32((r << 16) | (g << 8) | b)
}

func ceil16r(x int) int {
	return (x + 0xffff) >> 16
}

func imul16(x, y int) int {
	return int((int64(x) * int64(y)) >> 16)
}

func imul14(x, y int) int {
	return int((int64(x) * int64(y)) >> 14)
}

func idiv16(x, y int) int {
	return int((int64(x) << 16) / int64(y))
}
