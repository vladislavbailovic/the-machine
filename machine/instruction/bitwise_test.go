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
		SHL_REG_LIT.Pack(register.R1.AsUint16(), 4),
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
		SHR_REG_LIT.Pack(register.R1.AsUint16(), 2),
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

func Test_And_RegLit(t *testing.T) {
	cpu := cpu.NewCpu()
	mem := memory.NewMemory(255)
	packeds := [][]byte{
		MOV_LIT_R1.Pack(13),                          // 0b1101
		AND_REG_LIT.Pack(register.R1.AsUint16(), 11), // 0b1011
	} // 0b1001 == 9

	if cpu.GetRegister(register.Ac) != 0 {
		t.Fatalf("machine initial state error: expected empty Ac, got: %d", cpu.GetRegister(register.Ac))
	}

	for idx, packed := range packeds {
		instr, raw := unpackInstruction(packed)
		if err := instr.Executor.Execute(raw, cpu, mem); err != nil {
			t.Fatalf("%d: error executing instruction %v: %v", idx, instr, err)
		}
	}

	if cpu.GetRegister(register.R1) != 13 {
		t.Fatalf("error setting immediate value to register R1")
	}
	if cpu.GetRegister(register.Ac) != 9 {
		t.Fatalf("error setting result value in accumulator, expected 9 got %d",
			cpu.GetRegister(register.Ac))
	}
}

func Test_And_RegReg(t *testing.T) {
	cpu := cpu.NewCpu()
	mem := memory.NewMemory(255)
	packeds := [][]byte{
		MOV_LIT_R1.Pack(13), // 0b1101
		MOV_LIT_R2.Pack(11), // 0b1011
		AND_REG_REG.Pack(register.R1.AsUint16(), register.R2.AsUint16()), // 0b1001 == 9
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

	if cpu.GetRegister(register.R1) != 13 {
		t.Fatalf("error setting immediate value to register R1")
	}
	if cpu.GetRegister(register.R2) != 11 {
		t.Fatalf("error setting immediate value to register R2")
	}
	if cpu.GetRegister(register.Ac) != 9 {
		t.Fatalf("error setting result value in accumulator, expected 9 got %d",
			cpu.GetRegister(register.Ac))
	}
}

func Test_Or_RegLit(t *testing.T) {
	cpu := cpu.NewCpu()
	mem := memory.NewMemory(255)
	packeds := [][]byte{
		MOV_LIT_R1.Pack(13),                         // 0b1101
		OR_REG_LIT.Pack(register.R1.AsUint16(), 11), // 0b1011
	} // 0b1111 == 15

	if cpu.GetRegister(register.Ac) != 0 {
		t.Fatalf("machine initial state error: expected empty Ac, got: %d", cpu.GetRegister(register.Ac))
	}

	for idx, packed := range packeds {
		instr, raw := unpackInstruction(packed)
		if err := instr.Executor.Execute(raw, cpu, mem); err != nil {
			t.Fatalf("%d: error executing instruction %v: %v", idx, instr, err)
		}
	}

	if cpu.GetRegister(register.R1) != 13 {
		t.Fatalf("error setting immediate value to register R1")
	}
	if cpu.GetRegister(register.Ac) != 15 {
		t.Fatalf("error setting result value in accumulator, expected 15 got %d",
			cpu.GetRegister(register.Ac))
	}
}

func Test_Or_RegReg(t *testing.T) {
	cpu := cpu.NewCpu()
	mem := memory.NewMemory(255)
	packeds := [][]byte{
		MOV_LIT_R1.Pack(13), // 0b1101
		MOV_LIT_R2.Pack(11), // 0b1011
		OR_REG_REG.Pack(register.R1.AsUint16(), register.R2.AsUint16()), // 0b1111 == 15
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

	if cpu.GetRegister(register.R1) != 13 {
		t.Fatalf("error setting immediate value to register R1")
	}
	if cpu.GetRegister(register.R2) != 11 {
		t.Fatalf("error setting immediate value to register R2")
	}
	if cpu.GetRegister(register.Ac) != 15 {
		t.Fatalf("error setting result value in accumulator, expected 15 got %d",
			cpu.GetRegister(register.Ac))
	}
}

func Test_Xor_RegLit(t *testing.T) {
	cpu := cpu.NewCpu()
	mem := memory.NewMemory(255)
	packeds := [][]byte{
		MOV_LIT_R1.Pack(13),                          // 0b1101
		XOR_REG_LIT.Pack(register.R1.AsUint16(), 11), // 0b1011
	} // 0b0110 == 6

	if cpu.GetRegister(register.Ac) != 0 {
		t.Fatalf("machine initial state error: expected empty Ac, got: %d", cpu.GetRegister(register.Ac))
	}

	for idx, packed := range packeds {
		instr, raw := unpackInstruction(packed)
		if err := instr.Executor.Execute(raw, cpu, mem); err != nil {
			t.Fatalf("%d: error executing instruction %v: %v", idx, instr, err)
		}
	}

	if cpu.GetRegister(register.R1) != 13 {
		t.Fatalf("error setting immediate value to register R1")
	}
	if cpu.GetRegister(register.Ac) != 6 {
		t.Fatalf("error setting result value in accumulator, expected 6 got %d",
			cpu.GetRegister(register.Ac))
	}
}

func Test_Xor_RegReg(t *testing.T) {
	cpu := cpu.NewCpu()
	mem := memory.NewMemory(255)
	packeds := [][]byte{
		MOV_LIT_R1.Pack(13), // 0b1101
		MOV_LIT_R2.Pack(11), // 0b1011
		XOR_REG_REG.Pack(register.R1.AsUint16(), register.R2.AsUint16()), // 0b0110 == 6
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

	if cpu.GetRegister(register.R1) != 13 {
		t.Fatalf("error setting immediate value to register R1")
	}
	if cpu.GetRegister(register.R2) != 11 {
		t.Fatalf("error setting immediate value to register R2")
	}
	if cpu.GetRegister(register.Ac) != 6 {
		t.Fatalf("error setting result value in accumulator, expected 6 got %d",
			cpu.GetRegister(register.Ac))
	}
}
