package instruction

import (
	"testing"
	"the-machine/machine/cpu"
)

func Test_AddStack_Error(t *testing.T) {
	cpu := cpu.NewCpu()

	if head, err := cpu.Pop(); err == nil {
		t.Fatalf("expected no stack head, got %v", head)
	}

	if err := cpu.Push(161); err != nil {
		t.Fatalf("unable to set stack head, got %v", err)
	}

	packed := ADD_STACK.Pack()

	if _, err := runPackedInstructionWithCpu(packed, cpu); err == nil {
		t.Fatalf("expected stack underflow")
	}
}

func Test_AddStack(t *testing.T) {
	cpu := cpu.NewCpu()

	if head, err := cpu.Pop(); err == nil {
		t.Fatalf("expected no stack head, got %v", head)
	}

	if err := cpu.Push(13); err != nil {
		t.Fatalf("unable to set stack head, got %v", err)
	}

	if err := cpu.Push(12); err != nil {
		t.Fatalf("unable to set stack head, got %v", err)
	}

	packed := ADD_STACK.Pack()

	if _, err := runPackedInstructionWithCpu(packed, cpu); err != nil {
		t.Fatalf("error running instruction: %v", err)
	}

	result, err := cpu.Pop()
	if err != nil {
		t.Fatalf("expected stack op result as stack head, got: %v", err)
	}
	if result != 25 {
		t.Fatalf("expected stack op result to be 25, got %d", result)
	}
}

func Test_SubStack(t *testing.T) {
	cpu := cpu.NewCpu()

	if head, err := cpu.Pop(); err == nil {
		t.Fatalf("expected no stack head, got %v", head)
	}

	if err := cpu.Push(12); err != nil {
		t.Fatalf("unable to set stack head, got %v", err)
	}

	if err := cpu.Push(13); err != nil {
		t.Fatalf("unable to set stack head, got %v", err)
	}

	packed := SUB_STACK.Pack()

	if _, err := runPackedInstructionWithCpu(packed, cpu); err != nil {
		t.Fatalf("error running instruction: %v", err)
	}

	result, err := cpu.Pop()
	if err != nil {
		t.Fatalf("expected stack op result as stack head, got: %v", err)
	}
	if result != 1 {
		t.Fatalf("expected stack op result to be 1, got %d", result)
	}
}

func Test_MulStack(t *testing.T) {
	cpu := cpu.NewCpu()

	if head, err := cpu.Pop(); err == nil {
		t.Fatalf("expected no stack head, got %v", head)
	}

	if err := cpu.Push(12); err != nil {
		t.Fatalf("unable to set stack head, got %v", err)
	}

	if err := cpu.Push(13); err != nil {
		t.Fatalf("unable to set stack head, got %v", err)
	}

	packed := MUL_STACK.Pack()

	if _, err := runPackedInstructionWithCpu(packed, cpu); err != nil {
		t.Fatalf("error running instruction: %v", err)
	}

	result, err := cpu.Pop()
	if err != nil {
		t.Fatalf("expected stack op result as stack head, got: %v", err)
	}
	if result != 156 {
		t.Fatalf("expected stack op result to be 156, got %d", result)
	}
}

func Test_DivStack(t *testing.T) {
	cpu := cpu.NewCpu()

	if head, err := cpu.Pop(); err == nil {
		t.Fatalf("expected no stack head, got %v", head)
	}

	if err := cpu.Push(12); err != nil {
		t.Fatalf("unable to set stack head, got %v", err)
	}

	if err := cpu.Push(120); err != nil {
		t.Fatalf("unable to set stack head, got %v", err)
	}

	packed := DIV_STACK.Pack()

	if _, err := runPackedInstructionWithCpu(packed, cpu); err != nil {
		t.Fatalf("error running instruction: %v", err)
	}

	result, err := cpu.Pop()
	if err != nil {
		t.Fatalf("expected stack op result as stack head, got: %v", err)
	}
	if result != 10 {
		t.Fatalf("expected stack op result to be 10, got %d", result)
	}
}
