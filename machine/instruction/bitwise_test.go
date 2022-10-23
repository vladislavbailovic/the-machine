package instruction

import (
	"testing"
	"the-machine/machine/cpu"
	"the-machine/machine/memory"
	"the-machine/machine/register"
)

func Test_Shl_RegLit(t *testing.T) {
	cpu := cpu.NewCpu()
	mem := memory.NewMemory(255)
	packeds := [][]byte{
		MOV_LIT_R1.Pack(12),
		SHL_REG_LIT.Pack(uint16(register.R1.AsByte()), uint16(4)),
	}

	if cpu.GetRegister(register.Ac) != 0 {
		t.Fatalf("machine initial state error: expected empty Ac, got: %d", cpu.GetRegister(register.Ac))
	}

	for idx, packed := range packeds {
		instr, raw := unpackInstruction(packed)
		if err := instr.Executor.Execute(raw, cpu, mem); err != nil {
			t.Fatalf("%d: error executing instruction %v: %v", idx, instr, err)
		}
	}

	if cpu.GetRegister(register.R1) != 12 {
		t.Fatalf("error setting immediate value to register R1")
	}
	if cpu.GetRegister(register.Ac) != 192 {
		t.Fatalf("error setting result value in accumulator, expected 192 got %d",
			cpu.GetRegister(register.Ac))
	}
}

func Test_Shr_RegLit(t *testing.T) {
	cpu := cpu.NewCpu()
	mem := memory.NewMemory(255)
	packeds := [][]byte{
		MOV_LIT_R1.Pack(12),
		SHR_REG_LIT.Pack(uint16(register.R1.AsByte()), uint16(2)),
	}

	if cpu.GetRegister(register.Ac) != 0 {
		t.Fatalf("machine initial state error: expected empty Ac, got: %d", cpu.GetRegister(register.Ac))
	}

	for idx, packed := range packeds {
		instr, raw := unpackInstruction(packed)
		if err := instr.Executor.Execute(raw, cpu, mem); err != nil {
			t.Fatalf("%d: error executing instruction %v: %v", idx, instr, err)
		}
	}

	if cpu.GetRegister(register.R1) != 12 {
		t.Fatalf("error setting immediate value to register R1")
	}
	if cpu.GetRegister(register.Ac) != 3 {
		t.Fatalf("error setting result value in accumulator, expected 3 got %d",
			cpu.GetRegister(register.Ac))
	}
}

/*
func Test_Add_RegReg(t *testing.T) {
	cpu := cpu.NewCpu()
	mem := memory.NewMemory(255)
	packeds := [][]byte{
		MOV_LIT_R1.Pack(13),
		MOV_LIT_R2.Pack(12),
		ADD_REG_REG.Pack(uint16(register.R1.AsByte()), uint16(register.R2.AsByte())),
	}

	if cpu.GetRegister(register.Ac) != 0 {
		t.Fatalf("machine initial state error: expected empty Ac, got: %d", cpu.GetRegister(register.Ac))
	}

	for idx, packed := range packeds {
		instr, raw := unpackInstruction(packed)
		if err := instr.Executor.Execute(raw, cpu, mem); err != nil {
			t.Fatalf("%d: error executing instruction %v: %err", idx, instr, err)
		}
	}

	if cpu.GetRegister(register.R1) != 13 {
		t.Fatalf("error setting immediate value to register R1")
	}
	if cpu.GetRegister(register.R2) != 12 {
		t.Fatalf("error setting immediate value to register R2")
	}
	if cpu.GetRegister(register.Ac) != 25 {
		t.Fatalf("error setting result value in accumulator, expected 25 got %d",
			cpu.GetRegister(register.Ac))
	}
}
*/
