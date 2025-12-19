package main

func (o *Object3D) BuildText(font *VectorFont, text string, x, y, z, size, spacing float32) {
	if font == nil {
		return
	}

	segCount := 0
	for _, r := range text {
		g := font.Glyph(r)
		segCount += len(g.Segments)
	}

	o.pVert = make([]Vertex, segCount*2)
	o.pLines = make([]Line, segCount)
	o.pFaces = nil

	v := 0
	lineIdx := 0
	cursorX := x
	for _, r := range text {
		g := font.Glyph(r)
		for _, seg := range g.Segments {
			a := Vector{X: cursorX + seg.A.X*size, Y: y + seg.A.Y*size, Z: z}
			b := Vector{X: cursorX + seg.B.X*size, Y: y + seg.B.Y*size, Z: z}
			o.pVert[v] = Vertex{os: a, ws: a}
			o.pVert[v+1] = Vertex{os: b, ws: b}
			o.pLines[lineIdx] = newLine(v, v+1, nil, nil)
			v += 2
			lineIdx++
		}
		cursorX += (g.Width + spacing) * size
	}

	o.pLines = o.pLines[:lineIdx]
	o.m.Identity()
}
