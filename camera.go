package main

import "math"

type Camera struct {
	m         Matrix
	perspX    float32
	perspY    float32
	roll      float32
	eyepoint  Vector
	target    Vector
	posKeys   []float32
	targetKeys []float32
	rollKeys  []float32
}

func (c *Camera) Init() {
	c.perspX = 1.0
	c.perspY = c.perspX * (float32(screenWidth) / float32(screenHeight))
	c.eyepoint = Vector{X: 0, Y: 0, Z: -2}
	c.roll = float32(math.Pi)
	c.target = Vector{X: 0, Y: 0, Z: 1}
	c.m.Identity()
}

func (c *Camera) Camera2Matrix() {
	var cVec, pivot Vector
	var focus float32
	var ax, ay, az float32
	var sinx, siny, sinz, cosx, cosy, cosz float32

	c.m.Identity()

	pivot = c.eyepoint.Negated()
	cVec = c.target.Sub(c.eyepoint)
	focus = cVec.Magnitude()

	ax = -float32(math.Atan2(float64(cVec.X), float64(cVec.Z)))
	ay = float32(math.Asin(float64(cVec.Y / focus)))
	az = -c.roll

	sinx = float32(math.Sin(float64(ax)))
	cosx = float32(math.Cos(float64(ax)))
	siny = float32(math.Sin(float64(ay)))
	cosy = float32(math.Cos(float64(ay)))
	sinz = float32(math.Sin(float64(az)))
	cosz = float32(math.Cos(float64(az)))

	c.m.xx = (sinx*siny*sinz + cosx*cosz) * c.perspX
	c.m.yx = (cosy * sinz) * c.perspX
	c.m.zx = (sinx*cosz - cosx*siny*sinz) * c.perspX

	c.m.xy = (sinx*siny*cosz - cosx*sinz) * c.perspY
	c.m.yy = (cosy * cosz) * c.perspY
	c.m.zy = (-cosx*siny*cosz - sinx*sinz) * c.perspY

	c.m.xz = -sinx * cosy
	c.m.yz = siny
	c.m.zz = cosx * cosy

	c.m.PreTranslation(pivot)
}

func (c *Camera) BuildCam(r *binReader) {
	c.eyepoint = Vector{X: r.readF32(), Y: r.readF32(), Z: r.readF32()}
	c.target = Vector{X: r.readF32(), Y: r.readF32(), Z: r.readF32()}
	c.roll = r.readF32()
}

func (c *Camera) BuildCamKeys(r *binReader) {
	l1 := int(r.readU16())
	l2 := int(r.readU16())
	l3 := int(r.readU16())
	if l1 > 0 {
		c.posKeys = r.readF32Slice(l1 * 4)
	}
	if l2 > 0 {
		c.targetKeys = r.readF32Slice(l2 * 4)
	}
	if l3 > 0 {
		c.rollKeys = r.readF32Slice(l3 * 2)
	}
}
