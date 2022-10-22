package instruction

import (
	"testing"
	"the-machine/machine/cpu"
	"the-machine/machine/register"
)

func Test_PushReg(t *testing.T) {
	cpu := cpu.NewCpu()
	cpu.SetRegister(register.R1, 161)

	if head, err := cpu.Pop(); err == nil {
		t.Fatalf("expected no stack head, got %v", head)
	}

	packed := PUSH_REG.Pack(uint16(register.R1.AsByte()))

	if _, err := runPackedInstructionWithCpu(packed, cpu); err != nil {
		t.Fatalf("error pushing register to stack: %v", err)
	}

	head, err := cpu.Pop()
	if err != nil {
		t.Fatalf("expected stack head, got %v", err)
	}

	if head != 161 {
		t.Fatalf("expected stack head to be set to 161, got %v", head)
	}
}

func Test_PushLit(t *testing.T) {
	cpu := cpu.NewCpu()

	if head, err := cpu.Pop(); err == nil {
		t.Fatalf("expected no stack head, got %v", head)
	}

	packed := PUSH_LIT.Pack(uint16(161))

	if _, err := runPackedInstructionWithCpu(packed, cpu); err != nil {
		t.Fatalf("error pushing register to stack: %v", err)
	}

	head, err := cpu.Pop()
	if err != nil {
		t.Fatalf("expected stack head, got %v", err)
	}

	if head != 161 {
		t.Fatalf("expected stack head to be set to 161, got %v", head)
	}
}

func Test_PopReg(t *testing.T) {
	cpu := cpu.NewCpu()

	if head, err := cpu.Pop(); err == nil {
		t.Fatalf("expected no stack head, got %v", head)
	}

	if err := cpu.Push(161); err != nil {
		t.Fatalf("unable to set stack head, got %v", err)
	}

	packed := POP_REG.Pack(uint16(register.R1.AsByte()))

	if _, err := runPackedInstructionWithCpu(packed, cpu); err != nil {
		t.Fatalf("error pushing register to stack: %v", err)
	}

	head, err := cpu.Pop()
	if err == nil {
		t.Fatalf("expected instruction-removed stack head, got %v", head)
	}

	if cpu.GetRegister(register.R1) != 161 {
		t.Fatalf("expected register to be set to stack head, got %v", cpu.GetRegister(register.R1))
	}
}
