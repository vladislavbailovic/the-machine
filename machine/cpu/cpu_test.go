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
	if cpu.stackSize != stack-1 {
		t.Fatalf("expected stack size at %d, got %d", stack-1, cpu.stackSize)
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
	if cpu.stackSize != 0 {
		t.Fatalf("expected stack size to be zero, got %d", cpu.stackSize)
	}
}

func Test_StoreFrame(t *testing.T) {
	cpu := NewCpu()

	if cpu.GetRegister(register.Fp) != 0 {
		t.Fatalf("expected initial frame pointer at zero")
	}

	cpu.Push(161)
	cpu.Push(1312)
	oldStackSize := uint16(cpu.stackSize)
	oldStackPointer := cpu.GetRegister(register.Sp)

	if err := cpu.StoreFrame(); err != nil {
		t.Fatalf("error storing frame: %v", err)
	}

	framePointer := cpu.GetRegister(register.Fp)
	if framePointer == 0 {
		t.Fatalf("expected frame pointer to update")
	}
	if framePointer == oldStackPointer {
		t.Fatalf("expected frame pointer not to point to old stack head(%d), got %d", oldStackPointer, framePointer)
	}
	if framePointer != cpu.GetRegister(register.Sp) {
		t.Fatalf("expected frame pointer to point to new stack head, got %d", framePointer)
	}

	if cpu.stackSize != 0 {
		t.Fatalf("expected stack size to reset on frame store")
	}

	restoredStackSize, err := cpu.Pop()
	if err != nil {
		t.Fatalf("error popping stack size: %v", err)
	}
	if restoredStackSize != oldStackSize {
		t.Fatalf("expected stack head to be old stack size (%d), got %d", oldStackSize, restoredStackSize)
	}
}

func Test_RestoreFrame(t *testing.T) {
	cpu := NewCpu()

	cpu.Push(161)
	cpu.SetRegister(register.Ip, 161)
	cpu.SetRegister(register.R1, 1)
	cpu.SetRegister(register.R2, 3)
	cpu.SetRegister(register.R3, 1)
	cpu.SetRegister(register.R4, 2)
	oldStackSize := cpu.stackSize
	oldStackPointer := cpu.GetRegister(register.Sp)

	if err := cpu.StoreFrame(); err != nil {
		t.Fatalf("error storing frame: %v", err)
	}
	cpu.Push(1312)
	cpu.SetRegister(register.Ip, 1312)
	cpu.SetRegister(register.R1, 8)
	cpu.SetRegister(register.R2, 8)
	cpu.SetRegister(register.R3, 8)
	cpu.SetRegister(register.R4, 8)

	if err := cpu.RestoreFrame(); err != nil {
		t.Fatalf("error restoring frame: %v", err)
	}

	if cpu.stackSize == oldStackSize {
		t.Fatalf("stack size reset: expected empty, got %d", cpu.stackSize)
	}
	newStackPointer := cpu.GetRegister(register.Sp)
	if oldStackPointer == newStackPointer {
		t.Fatalf("stack pointer mismatch: expected 0 (pre-store: %d), got %d", oldStackPointer, newStackPointer)
	}

	// if cpu.GetRegister(register.Fp) != oldStackPointer {
	// 	t.Fatalf("expected frame pointer to reset to %d, got %d", oldStackSize, cpu.GetRegister(register.Fp))
	// }

	if cpu.stackSize != 0 {
		t.Fatalf("stuff left on stack: %d", cpu.stackSize)
	}

	if _, err := cpu.Pop(); err == nil {
		t.Fatalf("stack should be empty!")
	}

	if 161 != cpu.GetRegister(register.Ip) {
		t.Fatalf("error restoring Ip register: %d", cpu.GetRegister(register.Ip))
	}

	if 1 != cpu.GetRegister(register.R1) {
		t.Fatalf("error restoring R1 register: %d", cpu.GetRegister(register.R1))
	}

	if 3 != cpu.GetRegister(register.R2) {
		t.Fatalf("error restoring R2 register: %d", cpu.GetRegister(register.R2))
	}

	if 1 != cpu.GetRegister(register.R3) {
		t.Fatalf("error restoring R3 register: %d", cpu.GetRegister(register.R3))
	}

	if 2 != cpu.GetRegister(register.R4) {
		t.Fatalf("error restoring R4 register: %d", cpu.GetRegister(register.R4))
	}
}

func Test_NestedStoreRestore(t *testing.T) {
	cpu := NewCpu()

	regs := []register.Register{
		register.Ip,
		register.R1,
		register.R2,
		register.R3,
		register.R4,
	}

	// Initial state
	for _, reg := range regs {
		cpu.SetRegister(reg, 161)
	}
	if err := cpu.StoreFrame(); err != nil {
		t.Fatalf("error storing initial state: %v", err)
	}

	previousFp := uint16(0)
	previousSp := uint16(0)

	// Nested stores
	for i := 1; i < 10; i++ {
		storeVals(cpu, uint16(i))
		if err := cpu.StoreFrame(); err != nil {
			t.Fatalf("%d: error storing state: %v", i, err)
		}

		if cpu.stackSize != 0 {
			t.Fatalf("stack size not reset on store: %d", cpu.stackSize)
		}

		if cpu.GetRegister(register.Fp) == previousFp {
			t.Fatalf("frame pointer not updated, expected it not to be %d", previousFp)
		}
		previousFp = cpu.GetRegister(register.Fp)

		if cpu.GetRegister(register.Sp) == previousSp {
			t.Fatalf("stack pointer not updated, expected it not to be %d", previousSp)
		}
		previousSp = cpu.GetRegister(register.Sp)

		for _, reg := range regs {
			regVal := cpu.GetRegister(reg)
			if regVal != uint16(i) {
				t.Fatalf("register not updated: %v", reg)
			}
		}
	}

	// Restores
	for i := 1; i < 10; i++ {
		if err := cpu.RestoreFrame(); err != nil {
			t.Fatalf("%d: error restoring state: %v", i, err)
		}

		if cpu.stackSize != 0 {
			t.Fatalf("%d: stack size not reset on restore: %d", i, cpu.stackSize)
		}

		if cpu.GetRegister(register.Fp) == previousFp {
			t.Fatalf("%d: frame pointer not updated, expected it not to be %d", i, previousFp)
		}
		previousFp = cpu.GetRegister(register.Fp)

		if cpu.GetRegister(register.Sp) == previousSp {
			t.Fatalf("%d: stack pointer not updated, expected it not to be %d", i, previousSp)
		}
		previousSp = cpu.GetRegister(register.Sp)

		for _, reg := range regs {
			regVal := cpu.GetRegister(reg)
			if regVal != uint16(10-i) {
				t.Fatalf("%d: register not updated: %v == %d", i, reg, regVal)
			}
		}
	}

	if err := cpu.RestoreFrame(); err != nil {
		t.Fatalf("error restoring initial state: %v", err)
	}
	for j, reg := range regs {
		if cpu.GetRegister(reg) != 161 {
			t.Fatalf("%d: register not updated: %v (%d)", j, reg, cpu.GetRegister(reg))
		}
	}
}

func storeVals(cpu *Cpu, val uint16) {
	cpu.Push(1312 + val)
	cpu.Push(161 + val)
	cpu.SetRegister(register.Ip, val)
	cpu.SetRegister(register.R1, val)
	cpu.SetRegister(register.R2, val)
	cpu.SetRegister(register.R3, val)
	cpu.SetRegister(register.R4, val)
}
