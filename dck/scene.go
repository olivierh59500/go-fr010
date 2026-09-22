package fr010

import "github.com/olivierh59500/democonstructionkit/motion"

const texCount = 64

type Scene3D struct {
	cam     Camera
	objects []*Object3D
}

func NewScene3D() *Scene3D {
	s := &Scene3D{}
	s.cam.Init()
	return s
}

func (s *Scene3D) BuildScene(data []byte) {
	r := newBinReader(data)
	for {
		kind := int(r.readI16())
		if kind == -1 {
			return
		}
		switch kind {
		case 1:
			obj := NewObject3D()
			obj.BuildObj(r, len(s.objects)%texCount)
			s.objects = append(s.objects, obj)
		case 10:
			s.cam.BuildCam(r)
		case 11:
			s.cam.BuildCamKeys(r)
		}
	}
}

func (s *Scene3D) DrawScene(frame float32, ll *LineList, fl *FaceList) {
	if len(s.cam.posKeys) > 0 {
		applySpline(&s.cam.posTrack, s.cam.posKeys, frame, &s.cam.eyepoint)
	}
	if len(s.cam.targetKeys) > 0 {
		applySpline(&s.cam.targetTrack, s.cam.targetKeys, frame, &s.cam.target)
	}

	s.cam.Camera2Matrix()
	for _, obj := range s.objects {
		obj.Transform(s.cam.m)
		obj.Draw(ll, fl)
	}
}

func (s *Scene3D) BuildFirstScene(font *VectorFont) {
	s.objects = nil
	obj := NewObject3D()
	obj.BuildText(font, "Farbrausch", 30, 10, 0, -0.008, 1)
	s.objects = append(s.objects, obj)

	obj = NewObject3D()
	obj.BuildText(font, "and", 0, -5, 0, -0.008, 1)
	s.objects = append(s.objects, obj)

	obj = NewObject3D()
	obj.BuildText(font, "ScoopeX", 20, -15, 0, -0.008, 1)
	s.objects = append(s.objects, obj)
}

func (s *Scene3D) BuildSecondScene(font *VectorFont) {
	s.objects = nil
	obj := NewObject3D()
	obj.BuildText(font, "at", 30, 5, 0, -0.01, 1)
	s.objects = append(s.objects, obj)

	obj = NewObject3D()
	obj.BuildText(font, "ms2oo1", 20, -10, 0, -0.01, 1)
	s.objects = append(s.objects, obj)
}

func (s *Scene3D) BuildThirdScene(font *VectorFont) {
	s.objects = nil
	obj := NewObject3D()
	obj.BuildText(font, "art", 20, 0, 0, -0.02, 1)
	s.objects = append(s.objects, obj)

	obj = NewObject3D()
	obj.BuildCube(8.0)
	obj.m.Translation(Vector{X: -12, Y: -10, Z: 0})
	s.objects = append(s.objects, obj)
}

func (s *Scene3D) BuildGreetingScene(font *VectorFont, text string, x, y, z, size float32) {
	s.objects = nil
	obj := NewObject3D()
	obj.BuildText(font, text, x, y, z, size, 1)
	s.objects = append(s.objects, obj)
}

// applySpline compiles the authored camera track once and preserves its original
// float32 B-spline interpolation, including extrapolation outside key times.
func applySpline(track **motion.BSpline32, keys []float32, frame float32, out *Vector) {
	if *track == nil {
		compiled, err := motion.NewBSpline32(keys, 3)
		if err != nil {
			panic(err)
		}
		*track = compiled
	}
	var value [3]float32
	if (*track).Sample(value[:], frame) {
		out.X, out.Y, out.Z = value[0], value[1], value[2]
	}
}
