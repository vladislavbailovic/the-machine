package device

import (
	"bytes"
	"testing"
	"the-machine/machine/memory"
)

func Test_AddressToCoords(t *testing.T) {
	vga := NewVideo().(Video)

	// (255*13)+12
	if coords, err := vga.addressToCoords(memory.Address(3327)); err != nil {
		t.Fatalf("expected success in resolving coordinates, got: %v", err)
	} else {
		if coords[0] != 12 {
			t.Fatalf("expected 12 X coordinate, got %d", coords[0])
		}
		if coords[1] != 13 {
			t.Fatalf("expected 13 Y coordinate, got %d", coords[1])
		}
	}

	if coords, err := vga.addressToCoords(memory.Address(1312)); err != nil {
		t.Fatalf("expected success in resolving coordinates, got: %v", err)
	} else {
		if coords[0] != 37 {
			t.Fatalf("expected 37 X coordinate, got %d", coords[0])
		}
		if coords[1] != 5 {
			t.Fatalf("expected 5 Y coordinate, got %d", coords[1])
		}
	}

	if coords, err := vga.addressToCoords(memory.Address(161)); err != nil {
		t.Fatalf("expected success in resolving coordinates, got: %v", err)
	} else {
		if coords[0] != 161 {
			t.Fatalf("expected 161 X coordinate, got %d", coords[0])
		}
		if coords[1] != 0 {
			t.Fatalf("expected 0 Y coordinate, got %d", coords[1])
		}
	}
}

func Test_SetByte_DrawsChar(t *testing.T) {
	var output bytes.Buffer
	vga := Video{stream: &output}

	if err := vga.SetByte(1312, 65); err != nil {
		t.Fatalf("error rendering byte: %v", err)
	}
	if output.String() != "\u001B[5;37HA" {
		t.Fatalf("unexpected output rendered: %s", output.String())
	}
}
