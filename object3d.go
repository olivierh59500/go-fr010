package fr010

type TexCoord struct {
	U, V int
}

type Face struct {
	pIndex      [3]int
	tc          [3]TexCoord
	visible     bool
	theDrawFace *DrawFaceObj
	col         *PastelColor
	texNum      int
}

type Line struct {
	pIndex [2]int
	fcount int
	f      [2]*Face
}

type Vertex struct {
	os Vector
	ws Vector
}

type Object3D struct {
	pVert  []Vertex
	pFaces []Face
	pLines []Line
	m      Matrix
}

func NewObject3D() *Object3D {
	o := &Object3D{}
	o.m.Identity()
	return o
}

func (o *Object3D) BuildCube(size float32) {
	o.pVert = make([]Vertex, 8)
	o.pFaces = make([]Face, 12)
	o.pLines = make([]Line, 12)

	for i := 0; i < 8; i++ {
		x := -size
		if (i^(i>>1))&1 != 0 {
			x = size
		}
		y := size
		if i&2 != 0 {
			y = -size
		}
		z := -size
		if i&4 != 0 {
			z = size
		}
		o.pVert[i] = Vertex{os: Vector{X: x, Y: y, Z: z}, ws: Vector{X: x, Y: y, Z: z}}
	}

	o.pFaces[0] = newFace(0, 1, 2, Vector{X: 0.3, Y: 0.5, Z: 0.0}, TexCoord{U: 255, V: 0}, TexCoord{U: 0, V: 0}, TexCoord{U: 0, V: 255}, 0)
	o.pFaces[1] = newFace(0, 2, 3, Vector{X: 0.3, Y: 0.5, Z: 0.0}, TexCoord{U: 255, V: 0}, TexCoord{U: 0, V: 255}, TexCoord{U: 255, V: 255}, 0)
	o.pFaces[2] = newFace(5, 4, 7, Vector{X: 0.7, Y: 0.5, Z: 0.0}, TexCoord{U: 255, V: 0}, TexCoord{U: 0, V: 0}, TexCoord{U: 0, V: 255}, 1)
	o.pFaces[3] = newFace(5, 7, 6, Vector{X: 0.7, Y: 0.5, Z: 0.0}, TexCoord{U: 255, V: 0}, TexCoord{U: 0, V: 255}, TexCoord{U: 255, V: 255}, 1)
	o.pFaces[4] = newFace(1, 5, 6, Vector{X: 0.0, Y: 0.5, Z: 0.3}, TexCoord{U: 255, V: 0}, TexCoord{U: 0, V: 0}, TexCoord{U: 0, V: 255}, 2)
	o.pFaces[5] = newFace(1, 6, 2, Vector{X: 0.0, Y: 0.5, Z: 0.3}, TexCoord{U: 255, V: 0}, TexCoord{U: 0, V: 255}, TexCoord{U: 255, V: 255}, 2)
	o.pFaces[6] = newFace(4, 0, 3, Vector{X: 0.2, Y: 0.0, Z: 0.0}, TexCoord{U: 255, V: 0}, TexCoord{U: 0, V: 0}, TexCoord{U: 0, V: 255}, 3)
	o.pFaces[7] = newFace(4, 3, 7, Vector{X: 0.2, Y: 0.0, Z: 0.0}, TexCoord{U: 255, V: 0}, TexCoord{U: 0, V: 255}, TexCoord{U: 255, V: 255}, 3)
	o.pFaces[8] = newFace(4, 5, 1, Vector{X: 0.2, Y: 0.0, Z: 0.2}, TexCoord{U: 255, V: 0}, TexCoord{U: 0, V: 0}, TexCoord{U: 0, V: 255}, 4)
	o.pFaces[9] = newFace(4, 1, 0, Vector{X: 0.2, Y: 0.0, Z: 0.2}, TexCoord{U: 255, V: 0}, TexCoord{U: 0, V: 255}, TexCoord{U: 255, V: 255}, 4)
	o.pFaces[10] = newFace(3, 2, 6, Vector{X: 0.2, Y: 0.5, Z: 0.0}, TexCoord{U: 255, V: 0}, TexCoord{U: 0, V: 0}, TexCoord{U: 0, V: 255}, 5)
	o.pFaces[11] = newFace(3, 6, 7, Vector{X: 0.2, Y: 0.5, Z: 0.0}, TexCoord{U: 255, V: 0}, TexCoord{U: 0, V: 255}, TexCoord{U: 255, V: 255}, 5)

	o.pLines[0] = newLine(0, 1, &o.pFaces[0], &o.pFaces[9])
	o.pLines[1] = newLine(1, 2, &o.pFaces[0], &o.pFaces[5])
	o.pLines[2] = newLine(2, 3, &o.pFaces[1], &o.pFaces[10])
	o.pLines[3] = newLine(3, 0, &o.pFaces[1], &o.pFaces[6])
	o.pLines[4] = newLine(5, 4, &o.pFaces[2], &o.pFaces[8])
	o.pLines[5] = newLine(4, 7, &o.pFaces[2], &o.pFaces[7])
	o.pLines[6] = newLine(7, 6, &o.pFaces[3], &o.pFaces[11])
	o.pLines[7] = newLine(6, 5, &o.pFaces[3], &o.pFaces[4])
	o.pLines[8] = newLine(1, 5, &o.pFaces[4], &o.pFaces[8])
	o.pLines[9] = newLine(6, 2, &o.pFaces[5], &o.pFaces[10])
	o.pLines[10] = newLine(3, 7, &o.pFaces[7], &o.pFaces[11])
	o.pLines[11] = newLine(0, 4, &o.pFaces[6], &o.pFaces[9])
	o.m.Identity()
}

type rawFace struct {
	a, b, c int
	flag    byte
	u1, v1  byte
	u2, v2  byte
	u3, v3  byte
}

func (o *Object3D) BuildObj(r *binReader, texNum int) {
	verts := int(r.readI16())
	faces := int(r.readI16())

	o.pVert = make([]Vertex, verts)
	for i := 0; i < verts; i++ {
		x := r.readF32()
		y := r.readF32()
		z := r.readF32()
		o.pVert[i] = Vertex{os: Vector{X: x, Y: y, Z: z}, ws: Vector{X: x, Y: y, Z: z}}
	}

	rawFaces := make([]rawFace, faces)
	for i := 0; i < faces; i++ {
		rawFaces[i] = rawFace{
			a:    int(r.readU8()),
			b:    int(r.readU8()),
			c:    int(r.readU8()),
			flag: r.readU8(),
			u1:   r.readU8(),
			v1:   r.readU8(),
			u2:   r.readU8(),
			v2:   r.readU8(),
			u3:   r.readU8(),
			v3:   r.readU8(),
		}
	}

	cr := 1.0 - float32(r.readU8())/255.0
	cg := 1.0 - float32(r.readU8())/255.0
	cb := 1.0 - float32(r.readU8())/255.0
	_ = r.readU8()
	colorVec := Vector{X: cr, Y: cg, Z: cb}

	o.pFaces = make([]Face, faces)
	for i := 0; i < faces; i++ {
		rf := rawFaces[i]
		o.pFaces[i] = newFace(rf.a, rf.b, rf.c, colorVec,
			TexCoord{U: int(rf.u1), V: int(rf.v1)},
			TexCoord{U: int(rf.u2), V: int(rf.v2)},
			TexCoord{U: int(rf.u3), V: int(rf.v3)},
			texNum,
		)
	}

	lines := make([]Line, 0, faces)
	for i := 0; i < faces; i++ {
		rf := rawFaces[i]
		if rf.flag&1 != 0 {
			la, lb := rf.a, rf.b
			if !edgeShared(rawFaces, i+1, la, lb) {
				f1, f2 := findFaceRefs(rawFaces, o.pFaces, i, la, lb)
				lines = append(lines, newLine(la, lb, f1, f2))
			}
		}
		if rf.flag&2 != 0 {
			la, lb := rf.b, rf.c
			if !edgeShared(rawFaces, i+1, la, lb) {
				f1, f2 := findFaceRefs(rawFaces, o.pFaces, i, la, lb)
				lines = append(lines, newLine(la, lb, f1, f2))
			}
		}
		if rf.flag&4 != 0 {
			la, lb := rf.c, rf.a
			if !edgeShared(rawFaces, i+1, la, lb) {
				f1, f2 := findFaceRefs(rawFaces, o.pFaces, i, la, lb)
				lines = append(lines, newLine(la, lb, f1, f2))
			}
		}
	}

	o.pLines = lines
	o.m.Identity()
}

func (o *Object3D) Transform(mat Matrix) {
	n := mat
	n.MulInPlace(o.m)
	for i := range o.pVert {
		o.pVert[i].ws = n.Apply(o.pVert[i].os)
	}
}

func (o *Object3D) Draw(ll *LineList, fl *FaceList) {
	for i := range o.pFaces {
		f := &o.pFaces[i]
		z1 := float64(o.pVert[f.pIndex[0]].ws.Z)
		z2 := float64(o.pVert[f.pIndex[1]].ws.Z)
		z3 := float64(o.pVert[f.pIndex[2]].ws.Z)
		if z1 < nearZ {
			z1 = nearZ
		}
		if z2 < nearZ {
			z2 = nearZ
		}
		if z3 < nearZ {
			z3 = nearZ
		}
		z1 = 1.0 / z1
		z2 = 1.0 / z2
		z3 = 1.0 / z3

		x1 := float64(o.pVert[f.pIndex[0]].ws.X) * z1
		y1 := float64(o.pVert[f.pIndex[0]].ws.Y) * z1
		x2 := float64(o.pVert[f.pIndex[1]].ws.X) * z2
		y2 := float64(o.pVert[f.pIndex[1]].ws.Y) * z2
		x3 := float64(o.pVert[f.pIndex[2]].ws.X) * z3
		y3 := float64(o.pVert[f.pIndex[2]].ws.Y) * z3

		f.visible = ((x1-x2)*(y1-y3) - (x1-x3)*(y1-y2)) < 0
		if f.visible {
			fl.Add(o.pVert, f)
		} else {
			f.theDrawFace = nil
		}
	}

	for i := range o.pLines {
		l := &o.pLines[i]
		visible := true
		if l.fcount > 0 {
			visible = l.f[0].visible
		}
		if l.fcount > 1 {
			visible = visible || l.f[1].visible
		}

		if visible {
			dl := ll.Add(&o.pVert[l.pIndex[0]].ws, &o.pVert[l.pIndex[1]].ws)
			if dl == nil {
				continue
			}
			if l.fcount > 0 && l.f[0].theDrawFace != nil {
				df := l.f[0].theDrawFace
				if df.lcount < len(df.l) {
					df.l[df.lcount] = dl
					df.lcount++
				}
			}
			if l.fcount > 1 && l.f[1].theDrawFace != nil {
				df := l.f[1].theDrawFace
				if df.lcount < len(df.l) {
					df.l[df.lcount] = dl
					df.lcount++
				}
			}
		}
	}
}

func newFace(a, b, c int, col Vector, t1, t2, t3 TexCoord, tex int) Face {
	return Face{
		pIndex: [3]int{a, b, c},
		tc:     [3]TexCoord{t1, t2, t3},
		col:    NewPastelColor(col),
		texNum: tex,
	}
}

func newLine(a, b int, f1, f2 *Face) Line {
	l := Line{pIndex: [2]int{a, b}}
	if f1 != nil {
		l.f[l.fcount] = f1
		l.fcount++
	}
	if f2 != nil {
		l.f[l.fcount] = f2
		l.fcount++
	}
	return l
}

func edgeShared(faces []rawFace, start int, a, b int) bool {
	for j := start; j < len(faces); j++ {
		rf := faces[j]
		if rf.flag&1 != 0 && edgeMatch(a, b, rf.a, rf.b) {
			return true
		}
		if rf.flag&2 != 0 && edgeMatch(a, b, rf.b, rf.c) {
			return true
		}
		if rf.flag&4 != 0 && edgeMatch(a, b, rf.c, rf.a) {
			return true
		}
	}
	return false
}

func findFaceRefs(faces []rawFace, built []Face, end int, a, b int) (*Face, *Face) {
	var f1, f2 *Face
	for j := 0; j <= end; j++ {
		rf := faces[j]
		if rf.flag&1 != 0 && edgeMatch(a, b, rf.a, rf.b) {
			f1, f2 = addFaceRef(f1, f2, &built[j])
		}
		if rf.flag&2 != 0 && edgeMatch(a, b, rf.b, rf.c) {
			f1, f2 = addFaceRef(f1, f2, &built[j])
		}
		if rf.flag&4 != 0 && edgeMatch(a, b, rf.c, rf.a) {
			f1, f2 = addFaceRef(f1, f2, &built[j])
		}
	}
	return f1, f2
}

func addFaceRef(f1, f2 *Face, f *Face) (*Face, *Face) {
	if f1 == nil {
		return f, f2
	}
	if f2 == nil && f1 != f {
		return f1, f
	}
	return f1, f2
}

func edgeMatch(a, b, x, y int) bool {
	return (a == x && b == y) || (a == y && b == x)
}
