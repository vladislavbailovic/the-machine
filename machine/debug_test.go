package machine

import "testing"

func Test_GetFormat_GetPos(t *testing.T) {
	var pos string
	dbg := Debugger{}

	pos, _ = dbg.getFormat(Binary, Byte)
	if pos != "%10d" {
		t.Fatalf("wrong position format for Binary, Byte: %v", pos)
	}
}
