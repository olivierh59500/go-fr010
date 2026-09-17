// Package fr010 implements the FR-010 demo.
package fr010

import (
	"log"
	"time"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/audio"
)

type Game struct {
	pixels []uint32
	rgba   []byte
	frame  *ebiten.Image

	drawing  *Drawing
	lineList *LineList
	faceList *FaceList
	font     *VectorFont

	intro  *Scene3D
	scene1 *Scene3D
	scene2 *Scene3D
	scene3 *Scene3D
	scene4 *Scene3D
	scene5 *Scene3D
	scene6 *Scene3D

	audioContext  *audio.Context
	audioPlayer   *audio.Player
	ym            *YMPlayer
	audioReady    bool
	fallbackStart time.Time
	exitRequested bool
}

const (
	introDurationMS = 20000
	sampleRate      = 48000
	maxLayoutWidth  = 1280
)

func NewGame() (*Game, error) {
	textures := NewTextureManager(texCount)
	drawing := NewDrawing(screenWidth, screenHeight, textures)
	pixels, rgba := newPixelBuffer(screenWidth * screenHeight)
	drawing.SetBuffer(pixels, screenWidth)

	font, err := NewVectorFont()
	if err != nil {
		return nil, err
	}

	intro := NewScene3D()
	introLine1 := "A go/ebiten conversion"
	introLine2 := "by biliizir"
	introSize := float32(-0.006)
	introSpacing := float32(1)
	introX1 := -font.MeasureText(introLine1, introSpacing) * introSize * 0.5
	introX2 := -font.MeasureText(introLine2, introSpacing) * introSize * 0.5

	obj := NewObject3D()
	obj.BuildText(font, introLine1, introX1, 6, 0, introSize, introSpacing)
	intro.objects = append(intro.objects, obj)

	obj = NewObject3D()
	obj.BuildText(font, introLine2, introX2, -6, 0, introSize, introSpacing)
	intro.objects = append(intro.objects, obj)

	intro.cam.posKeys = scene1CamKeys

	scene1 := NewScene3D()
	scene1.BuildFirstScene(font)
	scene1.cam.posKeys = scene1CamKeys

	scene2 := NewScene3D()
	scene2.BuildSecondScene(font)
	scene2.cam.posKeys = scene2CamKeys

	scene3 := NewScene3D()
	scene3.BuildThirdScene(font)
	scene3.cam.posKeys = scene3CamKeys

	scene4 := NewScene3D()
	scene4.BuildScene(kasparovData)

	scene5 := NewScene3D()
	scene5.BuildScene(bugData)

	scene6 := NewScene3D()
	scene6.BuildScene(stadtData)

	return &Game{
		pixels:        pixels,
		rgba:          rgba,
		frame:         ebiten.NewImage(screenWidth, screenHeight),
		drawing:       drawing,
		lineList:      NewLineList(),
		faceList:      NewFaceList(),
		font:          font,
		intro:         intro,
		scene1:        scene1,
		scene2:        scene2,
		scene3:        scene3,
		scene4:        scene4,
		scene5:        scene5,
		scene6:        scene6,
		fallbackStart: time.Now(),
	}, nil
}

func (g *Game) Close() {
	if g.audioPlayer != nil {
		if err := g.audioPlayer.Close(); err != nil {
			log.Printf("close audio player: %v", err)
		}
		g.audioPlayer = nil
	}
	if g.ym != nil {
		if err := g.ym.Close(); err != nil {
			log.Printf("close YM player: %v", err)
		}
		g.ym = nil
	}
}

func (g *Game) Update() error {
	// mobile.SetGame constructs Game before Android has installed Ebitengine's
	// context. Opening audio on the first tick avoids blocking native startup.
	if !g.audioReady {
		g.audioReady = true
		g.fallbackStart = time.Now()
		g.initAudio()
	}

	if ebiten.IsKeyPressed(ebiten.KeyEscape) {
		g.exitRequested = true
	}
	if g.exitRequested {
		return ebiten.Termination
	}
	return nil
}

func (g *Game) Draw(screen *ebiten.Image) {
	frame := g.currentFrameMS()
	if frame < 0 {
		frame = 0
	}

	g.renderFrame(frame)

	g.frame.WritePixels(g.rgba)
	screen.Clear()
	op := &ebiten.DrawImageOptions{}
	op.GeoM.Translate(float64((screen.Bounds().Dx()-screenWidth)/2), 0)
	screen.DrawImage(g.frame, op)
}

func (g *Game) Layout(outsideWidth, outsideHeight int) (int, int) {
	return logicalWidth(outsideWidth, outsideHeight), screenHeight
}

func (g *Game) currentFrameMS() int {
	if g.audioPlayer != nil {
		return int(g.audioPlayer.Position().Milliseconds())
	}
	return int(time.Since(g.fallbackStart).Milliseconds())
}

func (g *Game) initAudio() {
	context := audio.NewContext(sampleRate)
	ym, err := NewYMPlayer(ymData, sampleRate, true)
	if err != nil {
		log.Printf("initialize YM stream: %v", err)
		return
	}

	player, err := context.NewPlayer(ym)
	if err != nil {
		_ = ym.Close()
		log.Printf("initialize Ebitengine audio: %v", err)
		return
	}

	g.audioContext = context
	g.audioPlayer = player
	g.ym = ym
	player.Play()
}

func logicalWidth(outsideWidth, outsideHeight int) int {
	if outsideWidth <= 0 || outsideHeight <= 0 {
		return screenWidth
	}

	width := (outsideWidth*screenHeight + outsideHeight - 1) / outsideHeight
	if width < screenWidth {
		return screenWidth
	}
	if width > maxLayoutWidth {
		return maxLayoutWidth
	}
	return width
}

func (g *Game) renderFrame(frame int) {
	demoFrame := frame
	introActive := frame < introDurationMS
	if !introActive {
		demoFrame = frame - introDurationMS
	}

	texnum = (demoFrame / 150) & (texCount - 1)
	g.drawing.t = float32(demoFrame) * 0.001

	for i := range g.pixels {
		g.pixels[i] = opaqueWhite
	}

	g.faceList.Clear()
	g.lineList.Clear()

	mz := false
	doDraw := true

	if introActive {
		g.intro.DrawScene(float32(demoFrame), g.lineList, g.faceList)
		g.faceList.DoClip(mz)
		g.lineList.DoClip(g.faceList)
		g.faceList.Draw(g.drawing)
		g.lineList.Draw(g.drawing)
		return
	}

	switch {
	case demoFrame < 20000:
		g.scene1.DrawScene(float32(demoFrame), g.lineList, g.faceList)
	case demoFrame < 40000:
		g.scene2.DrawScene(float32(demoFrame-20000), g.lineList, g.faceList)
	case demoFrame < 65000:
		if len(g.scene3.objects) > 1 {
			g.scene3.objects[1].m.Rotation(Vector{X: float32(demoFrame) / 2000.0, Y: float32(demoFrame) / 1900.0, Z: float32(demoFrame) / 1700.0})
		}
		g.scene3.DrawScene(float32(demoFrame-40000), g.lineList, g.faceList)
	case demoFrame < 68000:
		// pause
	case demoFrame <= 89500:
		g.scene4.DrawScene(float32(demoFrame-69000)/50.0, g.lineList, g.faceList)
	case demoFrame < 136000:
		doDraw = false
		g.scene5.DrawScene(float32(demoFrame-90000)/50.0, g.lineList, g.faceList)
		mz = true
		g.faceList.DoClip(mz)
		g.lineList.DoClip(g.faceList)
		g.faceList.Draw(g.drawing)
		g.lineList.Draw(g.drawing)
		g.drawGreetings(demoFrame)
	case demoFrame <= 157000:
		g.scene6.DrawScene(float32(demoFrame-140000)/20.0, g.lineList, g.faceList)
	case demoFrame < 167000:
		doDraw = false
		g.drawCredits()
	default:
		g.exitRequested = true
		return
	}

	if doDraw {
		g.faceList.DoClip(mz)
		g.lineList.DoClip(g.faceList)
		g.faceList.Draw(g.drawing)
		g.lineList.Draw(g.drawing)
	}
}

func (g *Game) drawGreetings(frame int) {
	if frame < 92000 {
		return
	}

	size := float32(0.05)
	spacing := float32(200)
	if frame < 94000 {
		g.font.DrawText(g.drawing, "PDM", 350, 250, 0.06, spacing)
		return
	}
	if frame < 98000 {
		g.font.DrawText(g.drawing, "Jinx", 100, 200, size, spacing)
		return
	}
	if frame < 102000 {
		g.font.DrawText(g.drawing, "All members of DMA", 35, 200, size, spacing)
		return
	}
	if frame < 106000 {
		g.font.DrawText(g.drawing, "The Union", 50, 250, size, spacing)
		return
	}
	if frame < 110000 {
		g.font.DrawText(g.drawing, "#TEAMG1", 50, 150, size, spacing)
		return
	}
	if frame < 114000 {
		g.font.DrawText(g.drawing, "Leonard (Oxg)", 25, 100, size, spacing)
		return
	}
	if frame < 118000 {
		g.font.DrawText(g.drawing, "RMS", 75, 200, size, spacing)
		return
	}
	if frame < 122000 {
		g.font.DrawText(g.drawing, "All I", 75, 150, size, spacing)
		g.font.DrawText(g.drawing, "forgot", 275, 250, size, spacing)
	}
}

func (g *Game) drawCredits() {
	size := float32(0.05)
	spacing := float32(200)
	g.font.DrawText(g.drawing, "original code:", 25, 75, size, spacing)
	g.font.DrawText(g.drawing, "entropy", 100, 150, size, spacing)
	g.font.DrawText(g.drawing, "peci", 100, 250, size, spacing)
	g.font.DrawText(g.drawing, "3d-design:", 25, 375, size, spacing)
	g.font.DrawText(g.drawing, "ruul", 100, 450, size, spacing)
}
