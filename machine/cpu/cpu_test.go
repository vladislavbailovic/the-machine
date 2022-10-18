package cpu

import (
	"math"
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

func Test_StackOverflow(t *testing.T) {
	cpu := NewCpu()
	stack := int(math.Floor(stackSize / 2))
	for i := 0; i < stack-1; i++ {
		if err := cpu.Push(uint16(1312 + i)); err != nil {
			t.Fatalf("error pushing to stack at idx %d: %v", i, err)
		}
	}
	sp := cpu.GetRegister(register.Sp)
	if sp != stackSize-(2+stackSize%2) {
		t.Fatalf("expected stack pointer at %d, got %d", stackSize-3, sp)
	}

	start := stack - 1
	for i := start; i > 0; i-- {
		expected := uint16(1312 + i - 1)
		if val, err := cpu.Pop(); err != nil || val != expected {
			t.Fatalf("error popping from stack at idx %d (from %d), got value %d (expected %d) and err %v",
				i, start, val, expected, err)
		}
	}

	sp = cpu.GetRegister(register.Sp)
	if sp != 0 {
		t.Fatalf("expected stack to be empty, got %d", sp)
	}
}
