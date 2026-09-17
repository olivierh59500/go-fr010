package fr010

import (
	"golang.org/x/image/font"
	"golang.org/x/image/font/sfnt"
	"golang.org/x/image/math/fixed"
)

const (
	fontHeight         = 1000.0
	fallbackGlyphWidth = 600.0
	curveStep          = 200.0
	maxCurveSegments   = 32
	closeEpsilon       = 0.01
)

type LineSegment struct {
	A, B Vector
}

type Glyph struct {
	Width    float32
	Segments []LineSegment
}

type VectorFont struct {
	font   *sfnt.Font
	buf    sfnt.Buffer
	ppem   fixed.Int26_6
	scale  float32
	glyphs map[rune]Glyph
	Height float32
}

func NewVectorFont() (*VectorFont, error) {
	return NewVectorFontFromTTF(secrcodeData)
}

func NewVectorFontFromTTF(data []byte) (*VectorFont, error) {
	f, err := sfnt.Parse(data)
	if err != nil {
		return nil, err
	}

	units := float32(f.UnitsPerEm())
	if units <= 0 {
		units = fontHeight
	}

	return &VectorFont{
		font:   f,
		ppem:   fixed.Int26_6(f.UnitsPerEm()),
		scale:  fontHeight / units,
		glyphs: make(map[rune]Glyph),
		Height: fontHeight,
	}, nil
}

func (f *VectorFont) Glyph(r rune) Glyph {
	if g, ok := f.glyphs[r]; ok {
		return g
	}

	g, ok := f.buildGlyph(r)
	if !ok && r != '?' {
		g, ok = f.buildGlyph('?')
	}
	if !ok {
		g = Glyph{Width: fallbackGlyphWidth}
	}

	f.glyphs[r] = g
	return g
}

func (f *VectorFont) buildGlyph(r rune) (Glyph, bool) {
	idx, err := f.font.GlyphIndex(&f.buf, r)
	if err != nil || idx == 0 {
		return Glyph{}, false
	}

	segments, err := f.font.LoadGlyph(&f.buf, idx, f.ppem, nil)
	if err != nil {
		return Glyph{}, false
	}

	lines := f.flattenSegments(segments)
	width := glyphWidth(lines)

	if width < 1 {
		if adv, err := f.font.GlyphAdvance(&f.buf, idx, f.ppem, font.HintingNone); err == nil {
			width = float32(adv) * f.scale
		}
	}
	if width <= 0 {
		width = fallbackGlyphWidth
	}

	return Glyph{Width: width, Segments: lines}, true
}

func (f *VectorFont) DrawText(d *Drawing, text string, x, y, size, spacing float32) {
	if d == nil {
		return
	}
	textLine = true
	defer func() { textLine = false }()

	cursorX := x
	for _, r := range text {
		g := f.Glyph(r)
		for _, seg := range g.Segments {
			a := Vector{X: cursorX + seg.A.X*size, Y: y + seg.A.Y*size, Z: 0}
			b := Vector{X: cursorX + seg.B.X*size, Y: y + seg.B.Y*size, Z: 0}
			d.DrawLine(a, b)
		}
		cursorX += (g.Width + spacing) * size
	}
}

func (f *VectorFont) MeasureText(text string, spacing float32) float32 {
	width := float32(0)
	count := 0
	for _, r := range text {
		g := f.Glyph(r)
		width += g.Width
		count++
	}
	if count > 1 {
		width += spacing * float32(count-1)
	}
	return width
}

func (f *VectorFont) point(p fixed.Point26_6) Vector {
	return Vector{
		X: float32(p.X) * f.scale,
		Y: float32(p.Y) * f.scale,
		Z: 0,
	}
}

func (f *VectorFont) flattenSegments(segments sfnt.Segments) []LineSegment {
	var lines []LineSegment
	var start Vector
	var prev Vector
	hasContour := false

	for _, seg := range segments {
		switch seg.Op {
		case sfnt.SegmentOpMoveTo:
			if hasContour {
				lines = appendLineSegment(lines, prev, start)
			}
			start = f.point(seg.Args[0])
			prev = start
			hasContour = true
		case sfnt.SegmentOpLineTo:
			p := f.point(seg.Args[0])
			lines = appendLineSegment(lines, prev, p)
			prev = p
		case sfnt.SegmentOpQuadTo:
			c := f.point(seg.Args[0])
			p := f.point(seg.Args[1])
			lines = appendQuad(lines, prev, c, p)
			prev = p
		case sfnt.SegmentOpCubeTo:
			c1 := f.point(seg.Args[0])
			c2 := f.point(seg.Args[1])
			p := f.point(seg.Args[2])
			lines = appendCube(lines, prev, c1, c2, p)
			prev = p
		}
	}

	if hasContour {
		lines = appendLineSegment(lines, prev, start)
	}

	return lines
}

func appendLineSegment(lines []LineSegment, a, b Vector) []LineSegment {
	if pointsClose(a, b) {
		return lines
	}
	return append(lines, LineSegment{A: a, B: b})
}

func appendQuad(lines []LineSegment, p0, p1, p2 Vector) []LineSegment {
	length := p0.Sub(p1).Magnitude() + p1.Sub(p2).Magnitude()
	steps := curveSteps(length)

	prev := p0
	for i := 1; i <= steps; i++ {
		t := float32(i) / float32(steps)
		inv := 1 - t
		p := Vector{
			X: inv*inv*p0.X + 2*inv*t*p1.X + t*t*p2.X,
			Y: inv*inv*p0.Y + 2*inv*t*p1.Y + t*t*p2.Y,
			Z: 0,
		}
		lines = appendLineSegment(lines, prev, p)
		prev = p
	}

	return lines
}

func appendCube(lines []LineSegment, p0, p1, p2, p3 Vector) []LineSegment {
	length := p0.Sub(p1).Magnitude() + p1.Sub(p2).Magnitude() + p2.Sub(p3).Magnitude()
	steps := curveSteps(length)

	prev := p0
	for i := 1; i <= steps; i++ {
		t := float32(i) / float32(steps)
		inv := 1 - t
		inv2 := inv * inv
		t2 := t * t
		p := Vector{
			X: inv2*inv*p0.X + 3*inv2*t*p1.X + 3*inv*t2*p2.X + t2*t*p3.X,
			Y: inv2*inv*p0.Y + 3*inv2*t*p1.Y + 3*inv*t2*p2.Y + t2*t*p3.Y,
			Z: 0,
		}
		lines = appendLineSegment(lines, prev, p)
		prev = p
	}

	return lines
}

func curveSteps(length float32) int {
	steps := int(length/curveStep) + 1
	if steps < 1 {
		steps = 1
	}
	if steps > maxCurveSegments {
		steps = maxCurveSegments
	}
	return steps
}

func pointsClose(a, b Vector) bool {
	dx := a.X - b.X
	dy := a.Y - b.Y
	return dx*dx+dy*dy <= closeEpsilon*closeEpsilon
}

func glyphWidth(lines []LineSegment) float32 {
	if len(lines) == 0 {
		return 0
	}

	minX := lines[0].A.X
	maxX := lines[0].A.X

	for _, seg := range lines {
		if seg.A.X < minX {
			minX = seg.A.X
		}
		if seg.A.X > maxX {
			maxX = seg.A.X
		}
		if seg.B.X < minX {
			minX = seg.B.X
		}
		if seg.B.X > maxX {
			maxX = seg.B.X
		}
	}

	if maxX < minX {
		return 0
	}
	return maxX - minX
}
