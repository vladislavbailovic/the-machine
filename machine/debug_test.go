package machine

import "testing"

func Test_GetFormat_GetPos(t *testing.T) {
	var pos string
	dbg := Formatter{Numbers: Binary, OutputAs: Byte}

	pos, _ = dbg.getFormat()
	if pos != "%10d" {
		t.Fatalf("wrong position format for Binary, Byte: %v", pos)
	}
}
