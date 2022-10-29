package device

import (
	"testing"
	"the-machine/machine/memory"
)

func Test_WriteToStdout(t *testing.T) {
	io := NewIoMap()
	if err := io.SetByte(memory.Address(Stdout), byte('H')); err != nil {
		t.Fatalf("unexpected error writing to stdout: %v", err)
	}
}
