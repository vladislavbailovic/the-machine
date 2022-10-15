package machine

import (
	"testing"
	"the-machine/machine/register"
)

func Test_Register_InstructionPointer(t *testing.T) {
	cpu := NewCpu()
	if val, err := cpu.getRegister(register.Ip); err != nil || val != 0 {
		t.Fatalf("expected zero value in uninitialized register, got %d and error %v", val, err)
	}
	if err := cpu.setRegister(register.Ip, 161); err != nil {
		t.Fatalf("expected setting register to succeed: %v", err)
	}
	if val, err := cpu.getRegister(register.Ip); err != nil || val != 161 {
		t.Fatalf("expected specific value 161 in set register, got %d and error %v", val, err)
	}
}
