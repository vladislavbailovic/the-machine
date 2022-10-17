package cpu

import (
	"testing"
	"the-machine/machine/register"
)

func Test_Register_InstructionPointer(t *testing.T) {
	cpu := NewCpu()
	if val := cpu.GetRegister(register.Ip); val != 0 {
		t.Fatalf("expected zero value in uninitialized register, got %d", val)
	}
	cpu.SetRegister(register.Ip, 161)
	if val := cpu.GetRegister(register.Ip); val != 161 {
		t.Fatalf("expected specific value 161 in set register, got %d", val)
	}
}
