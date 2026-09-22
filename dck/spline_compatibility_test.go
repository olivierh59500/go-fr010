package fr010

import (
	"github.com/olivierh59500/democonstructionkit/motion"
	"math"
	"testing"
)

// legacySpline preserves the previous segment selection, float32 weights and
// authored endpoint convention independently of DCK's binary-search track.
func legacySpline(keyArray []float32, curFrame float32, out []float32, numFloats int) {
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

func TestSharedBSplinePreservesAllCameraTracks(t *testing.T) {
	tracks := [][]float32{scene1CamKeys, scene2CamKeys, scene3CamKeys}
	for _, data := range [][]byte{kasparovData, bugData, stadtData} {
		scene := NewScene3D()
		scene.BuildScene(data)
		if len(scene.cam.posKeys) > 0 {
			tracks = append(tracks, scene.cam.posKeys)
		}
		if len(scene.cam.targetKeys) > 0 {
			tracks = append(tracks, scene.cam.targetKeys)
		}
	}
	for index, keys := range tracks {
		compiled, err := motion.NewBSpline32(keys, 3)
		if err != nil {
			t.Fatalf("track %d: %v", index, err)
		}
		times := []float32{-1000, -1, 0, .25, 1, 100, 2000, 22000, 30000}
		for key := 0; key < len(keys); key += 4 {
			times = append(times, keys[key]-.25, keys[key], keys[key]+.25)
		}
		var got, want [3]float32
		for _, at := range times {
			legacySpline(keys, at, want[:], 3)
			if !compiled.Sample(got[:], at) {
				t.Fatal("sample rejected")
			}
			for component := range got {
				if math.Float32bits(got[component]) != math.Float32bits(want[component]) {
					t.Fatalf("track %d at %g component %d: %g vs %g", index, at, component, got[component], want[component])
				}
			}
		}
		var vector Vector
		var cached *motion.BSpline32
		applySpline(&cached, keys, 12.5, &vector)
		if allocations := testing.AllocsPerRun(100, func() { applySpline(&cached, keys, 12.5, &vector) }); allocations != 0 {
			t.Fatalf("camera sample allocations = %g", allocations)
		}
	}
}
