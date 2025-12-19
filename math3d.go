package main

import "math"

type Quaternion struct {
	X, Y, Z, W float32
}

func (q *Quaternion) Identity() {
	q.X, q.Y, q.Z = 0, 0, 0
	q.W = 1
}

type Vector struct {
	X, Y, Z float32
}

func (v Vector) Add(o Vector) Vector {
	return Vector{X: v.X + o.X, Y: v.Y + o.Y, Z: v.Z + o.Z}
}

func (v *Vector) AddInPlace(o Vector) {
	v.X += o.X
	v.Y += o.Y
	v.Z += o.Z
}

func (v Vector) Sub(o Vector) Vector {
	return Vector{X: v.X - o.X, Y: v.Y - o.Y, Z: v.Z - o.Z}
}

func (v *Vector) SubInPlace(o Vector) {
	v.X -= o.X
	v.Y -= o.Y
	v.Z -= o.Z
}

func (v Vector) MulScalar(s float32) Vector {
	return Vector{X: v.X * s, Y: v.Y * s, Z: v.Z * s}
}

func (v *Vector) MulScalarInPlace(s float32) {
	v.X *= s
	v.Y *= s
	v.Z *= s
}

func (v Vector) DivScalar(s float32) Vector {
	inv := 1 / s
	return Vector{X: v.X * inv, Y: v.Y * inv, Z: v.Z * inv}
}

func (v Vector) Dot(o Vector) float32 {
	return v.X*o.X + v.Y*o.Y + v.Z*o.Z
}

func (v Vector) Cross(o Vector) Vector {
	return Vector{
		X: v.Y*o.Z - v.Z*o.Y,
		Y: v.Z*o.X - v.X*o.Z,
		Z: v.X*o.Y - v.Y*o.X,
	}
}

func (v Vector) Magnitude() float32 {
	return float32(math.Sqrt(float64(v.X*v.X + v.Y*v.Y + v.Z*v.Z)))
}

func (v Vector) MagnitudeSq() float32 {
	return v.X*v.X + v.Y*v.Y + v.Z*v.Z
}

func (v *Vector) Normalize() {
	len := v.Magnitude()
	if len == 0 {
		return
	}
	inv := 1 / len
	v.X *= inv
	v.Y *= inv
	v.Z *= inv
}

func (v Vector) Negated() Vector {
	return Vector{X: -v.X, Y: -v.Y, Z: -v.Z}
}

func (v Vector) Abs() Vector {
	ax := v.X
	if ax < 0 {
		ax = -ax
	}
	ay := v.Y
	if ay < 0 {
		ay = -ay
	}
	az := v.Z
	if az < 0 {
		az = -az
	}
	return Vector{X: ax, Y: ay, Z: az}
}

type Matrix struct {
	xx, xy, xz float32
	yx, yy, yz float32
	zx, zy, zz float32
	wx, wy, wz float32
}

func (m *Matrix) Identity() {
	m.xy, m.xz = 0, 0
	m.yx, m.yz = 0, 0
	m.zx, m.zy = 0, 0
	m.wx, m.wy, m.wz = 0, 0, 0
	m.xx, m.yy, m.zz = 1, 1, 1
}

func (m *Matrix) Translation(v Vector) {
	m.wx += v.X
	m.wy += v.Y
	m.wz += v.Z
}

func (m *Matrix) PreTranslation(v Vector) {
	m.wx = m.xx*v.X + m.yx*v.Y + m.zx*v.Z + m.wx
	m.wy = m.xy*v.X + m.yy*v.Y + m.zy*v.Z + m.wy
	m.wz = m.xz*v.X + m.yz*v.Y + m.zz*v.Z + m.wz
}

func (m *Matrix) Rotation(v Vector) {
	sx := float32(math.Sin(float64(v.X)))
	sy := float32(math.Sin(float64(v.Y)))
	sz := float32(math.Sin(float64(v.Z)))
	cx := float32(math.Cos(float64(v.X)))
	cy := float32(math.Cos(float64(v.Y)))
	cz := float32(math.Cos(float64(v.Z)))

	m.xx = cy * cz
	m.xy = cy * sz
	m.xz = -sy
	m.yx = sx*sy*cz - cx*sz
	m.yy = sx*sy*sz + cx*cz
	m.yz = sx * cy
	m.zx = cx*sy*cz + sx*sz
	m.zy = cx*sy*sz - sx*cz
	m.zz = cx * cy
}

func (m Matrix) Mul(b Matrix) Matrix {
	return Matrix{
		xx: m.xx*b.xx + m.yx*b.xy + m.zx*b.xz,
		yx: m.xx*b.yx + m.yx*b.yy + m.zx*b.yz,
		zx: m.xx*b.zx + m.yx*b.zy + m.zx*b.zz,
		wx: m.xx*b.wx + m.yx*b.wy + m.zx*b.wz + m.wx,
		xy: m.xy*b.xx + m.yy*b.xy + m.zy*b.xz,
		yy: m.xy*b.yx + m.yy*b.yy + m.zy*b.yz,
		zy: m.xy*b.zx + m.yy*b.zy + m.zy*b.zz,
		wy: m.xy*b.wx + m.yy*b.wy + m.zy*b.wz + m.wy,
		xz: m.xz*b.xx + m.yz*b.xy + m.zz*b.xz,
		yz: m.xz*b.yx + m.yz*b.yy + m.zz*b.yz,
		zz: m.xz*b.zx + m.yz*b.zy + m.zz*b.zz,
		wz: m.xz*b.wx + m.yz*b.wy + m.zz*b.wz + m.wz,
	}
}

func (m *Matrix) MulInPlace(b Matrix) {
	*m = m.Mul(b)
}

func (m Matrix) Apply(v Vector) Vector {
	return Vector{
		X: v.X*m.xx + v.Y*m.yx + v.Z*m.zx + m.wx,
		Y: v.X*m.xy + v.Y*m.yy + v.Z*m.zy + m.wy,
		Z: v.X*m.xz + v.Y*m.yz + v.Z*m.zz + m.wz,
	}
}

func (m Matrix) ApplyNoTranslation(v Vector) Vector {
	return Vector{
		X: v.X*m.xx + v.Y*m.yx + v.Z*m.zx,
		Y: v.X*m.xy + v.Y*m.yy + v.Z*m.zy,
		Z: v.X*m.xz + v.Y*m.yz + v.Z*m.zz,
	}
}
