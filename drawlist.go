package fr010

type DrawFaceObj struct {
	x1, x2, x3 float64
	y1, y2, y3 float64
	z1, z2, z3 float64
	minx, maxx float64
	miny, maxy float64
	minZ, maxZ float64
	u1, v1     float64
	u2, v2     float64
	u3, v3     float64
	l          [3]*DrawLineObj
	lcount     int
	visible    bool
	col        *PastelColor
	texNum     int
}

type DrawLineObj struct {
	x1, x2  float64
	y1, y2  float64
	z1, z2  float64
	minZ    float64
	maxZ    float64
	visible bool
}

type LineList struct {
	dl           []*DrawLineObj
	count        int
	scratchCount int
	scratchLines []DrawLineObj
}

func NewLineList() *LineList {
	return &LineList{
		dl:           make([]*DrawLineObj, maxLine),
		scratchLines: make([]DrawLineObj, maxLine),
	}
}

func (ll *LineList) Clear() {
	ll.count = 0
	ll.scratchCount = 0
}

func (ll *LineList) Add(v1, v2 *Vector) *DrawLineObj {
	if ll.scratchCount >= len(ll.scratchLines) || ll.count >= len(ll.dl) {
		return nil
	}
	l := &ll.scratchLines[ll.scratchCount]
	ll.scratchCount++
	l.visible = true
	l.x1 = float64(v1.X)
	l.y1 = float64(v1.Y)
	l.z1 = float64(v1.Z)
	l.x2 = float64(v2.X)
	l.y2 = float64(v2.Y)
	l.z2 = float64(v2.Z)
	ll.dl[ll.count] = l
	ll.count++
	return l
}

func (ll *LineList) Draw(d *Drawing) {
	for i := 0; i < ll.count; i++ {
		l := ll.dl[i]
		if l.visible {
			d.DrawLine(
				Vector{X: float32(l.x1), Y: float32(l.y1), Z: float32(l.z1)},
				Vector{X: float32(l.x2), Y: float32(l.y2), Z: float32(l.z2)},
			)
		}
	}
}

type FaceList struct {
	dl           []*DrawFaceObj
	count        int
	scratchCount int
	scratchFaces []DrawFaceObj
}

func NewFaceList() *FaceList {
	return &FaceList{
		dl:           make([]*DrawFaceObj, maxFace),
		scratchFaces: make([]DrawFaceObj, maxFace),
	}
}

func (fl *FaceList) Clear() {
	fl.count = 0
	fl.scratchCount = 0
}

func (fl *FaceList) Add(verts []Vertex, face *Face) {
	if fl.scratchCount >= len(fl.scratchFaces) || fl.count >= len(fl.dl) {
		return
	}
	f := &fl.scratchFaces[fl.scratchCount]
	fl.scratchCount++
	f.visible = true
	v1 := verts[face.pIndex[0]].ws
	v2 := verts[face.pIndex[1]].ws
	v3 := verts[face.pIndex[2]].ws
	f.x1 = float64(v1.X)
	f.y1 = float64(v1.Y)
	f.z1 = float64(v1.Z)
	f.x2 = float64(v2.X)
	f.y2 = float64(v2.Y)
	f.z2 = float64(v2.Z)
	f.x3 = float64(v3.X)
	f.y3 = float64(v3.Y)
	f.z3 = float64(v3.Z)
	f.u1 = float64(face.tc[0].U)
	f.v1 = float64(face.tc[0].V)
	f.u2 = float64(face.tc[1].U)
	f.v2 = float64(face.tc[1].V)
	f.u3 = float64(face.tc[2].U)
	f.v3 = float64(face.tc[2].V)
	f.col = face.col
	f.texNum = (face.texNum + texnum) & (texCount - 1)
	f.lcount = 0
	fl.dl[fl.count] = f
	fl.count++
	face.theDrawFace = f
}

func (fl *FaceList) Draw(d *Drawing) {
	for i := 0; i < fl.count; i++ {
		f := fl.dl[i]
		if f.visible {
			d.drawFace(f)
		}
	}
}
