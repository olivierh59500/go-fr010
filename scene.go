package main

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
		applySpline(s.cam.posKeys, frame, &s.cam.eyepoint)
	}
	if len(s.cam.targetKeys) > 0 {
		applySpline(s.cam.targetKeys, frame, &s.cam.target)
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

func spline(keyArray []float32, curFrame float32, out []float32, numFloats int) {
	stride := numFloats + 1
	if len(keyArray) < stride*4 {
		return
	}

	keys := 0
	maxStart := (len(keyArray)/stride - 4) * stride
	for keys+stride < len(keyArray) && keyArray[keys+stride] < curFrame {
		keys += stride
	}
	if keys > maxStart {
		keys = maxStart
	}

	spT := (curFrame - keyArray[keys]) / (keyArray[keys+stride] - keyArray[keys])
	spT2 := spT * spT
	spT3 := spT2 * spT

	k1 := ((-1.0 / 6.0) * spT3) + (0.5 * spT2) - (0.5 * spT) + (1.0 / 6.0)
	k2 := (0.5 * spT3) - spT2 + (2.0 / 3.0)
	k3 := ((-0.5) * spT3) + (0.5 * spT2) + (0.5 * spT) + (1.0 / 6.0)
	k4 := (1.0 / 6.0) * spT3

	for i := 0; i < numFloats; i++ {
		out[i] = keyArray[keys+i+1]*k1 +
			keyArray[keys+i+1+stride]*k2 +
			keyArray[keys+i+1+stride*2]*k3 +
			keyArray[keys+i+1+stride*3]*k4
	}
}

func applySpline(keys []float32, frame float32, out *Vector) {
	var tmp [3]float32
	spline(keys, frame, tmp[:], 3)
	out.X = tmp[0]
	out.Y = tmp[1]
	out.Z = tmp[2]
}
