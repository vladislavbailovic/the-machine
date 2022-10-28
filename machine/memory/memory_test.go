package memory

import "testing"

func Test_Memory(t *testing.T) {
	m := NewMemory(255)
	var x uint16
	var err error

	if x, err = m.GetUint16(13); err != nil || x != 0 {
		t.Fatalf("expected zero at unintialized memory offset, got %d and error %v", x, err)
	}

	if err = m.SetUint16(13, 12); err != nil {
		t.Fatalf("expected success setting memory at 13, got: %v", err)
	}

	if x, err = m.GetUint16(13); err != nil || x != 12 {
		t.Fatalf("expected specific value 12 at set memory offset, got %d and error %v", x, err)
	}

	m.SetUint16(161, 1312)
	if x, err = m.GetUint16(161); err != nil || x != 1312 {
		t.Fatalf("expected specific value 1312 at set memory offset, got %d and error %v", x, err)
	}
}
