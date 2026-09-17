package fr010

import (
	"encoding/binary"
	"math"
	"testing"
)

func TestPixelBufferUsesRGBAByteOrder(t *testing.T) {
	pixels, rgba := newPixelBuffer(2)
	pixels[0] = packRGB(0x12, 0x34, 0x56)
	pixels[1] = opaqueWhite
	want := []byte{0x12, 0x34, 0x56, 0xff, 0xff, 0xff, 0xff, 0xff}
	for i := range want {
		if rgba[i] != want[i] {
			t.Fatalf("rgba[%d] = %#02x, want %#02x", i, rgba[i], want[i])
		}
	}
}

func TestDarkenColorSaturatesChannels(t *testing.T) {
	got := darkenColor(packRGB(0x20, 0x80, 0xff))
	want := packRGB(0x00, 0x20, 0x9f)
	if got != want {
		t.Fatalf("darkenColor = %#08x, want %#08x", got, want)
	}
}

func TestFadeTableAccuracy(t *testing.T) {
	const samples = 10000
	maxError := float64(0)
	for i := 0; i <= samples; i++ {
		x := float32(i) / samples
		want := 0.5 * (1 - math.Cos(float64(x)*math.Pi))
		error := math.Abs(float64(fsc(x)) - want)
		if error > maxError {
			maxError = error
		}
	}
	if maxError > 8e-4 {
		t.Fatalf("fade lookup max error = %g, want <= 8e-4", maxError)
	}
}

func TestYMPlayerReadProducesStereoWithoutAllocating(t *testing.T) {
	player, err := NewYMPlayer(ymData, sampleRate, true)
	if err != nil {
		t.Fatal(err)
	}
	t.Cleanup(func() {
		if err := player.Close(); err != nil {
			t.Errorf("Close: %v", err)
		}
	})

	buffer := make([]byte, ymBufferSize*4)
	read := func() {
		n, err := player.Read(buffer)
		if err != nil {
			t.Fatalf("Read: %v", err)
		}
		if n != len(buffer) {
			t.Fatalf("Read bytes = %d, want %d", n, len(buffer))
		}
	}

	read()
	nonZero := false
	for i := 0; i < len(buffer); i += 4 {
		left := binary.LittleEndian.Uint16(buffer[i : i+2])
		right := binary.LittleEndian.Uint16(buffer[i+2 : i+4])
		if left != right {
			t.Fatalf("frame %d is not mono duplicated to stereo: %d != %d", i/4, left, right)
		}
		nonZero = nonZero || left != 0
	}
	if !nonZero {
		t.Fatal("YM stream contains only silence")
	}

	if allocations := testing.AllocsPerRun(20, read); allocations != 0 {
		t.Fatalf("Read allocations = %v, want 0", allocations)
	}
}

func TestLogicalWidth(t *testing.T) {
	tests := []struct {
		name          string
		outsideWidth  int
		outsideHeight int
		want          int
	}{
		{name: "unknown", want: screenWidth},
		{name: "native", outsideWidth: 640, outsideHeight: 480, want: 640},
		{name: "portrait", outsideWidth: 1080, outsideHeight: 2424, want: 640},
		{name: "Pixel 10a landscape", outsideWidth: 2424, outsideHeight: 1080, want: 1078},
		{name: "ultrawide cap", outsideWidth: 4000, outsideHeight: 480, want: maxLayoutWidth},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			if got := logicalWidth(test.outsideWidth, test.outsideHeight); got != test.want {
				t.Fatalf("logicalWidth(%d, %d) = %d, want %d", test.outsideWidth, test.outsideHeight, got, test.want)
			}
		})
	}
}
