package instruction

import (
	"testing"
	"the-machine/machine/cpu"
	"the-machine/machine/memory"
	"the-machine/machine/register"
)

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

func Test_Add_RegLit(t *testing.T) {
	cpu := cpu.NewCpu()
	mem := memory.NewMemory(255)
	packeds := [][]byte{
		MOV_LIT_R1.Pack(12),
		ADD_REG_LIT.Pack(uint16(register.R1.AsByte()), uint16(1)),
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

	if cpu.GetRegister(register.R1) != 12 {
		t.Fatalf("error setting immediate value to register R1")
	}
	if cpu.GetRegister(register.Ac) != 13 {
		t.Fatalf("error setting result value in accumulator, expected 25 got %d",
			cpu.GetRegister(register.Ac))
	}
}

func Test_Sub_RegReg(t *testing.T) {
	cpu := cpu.NewCpu()
	mem := memory.NewMemory(255)
	packeds := [][]byte{
		MOV_LIT_R1.Pack(13),
		MOV_LIT_R2.Pack(12),
		SUB_REG_REG.Pack(uint16(register.R1.AsByte()), uint16(register.R2.AsByte())),
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
	if cpu.GetRegister(register.Ac) != 1 {
		t.Fatalf("error setting result value in accumulator, expected 1 got %d",
			cpu.GetRegister(register.Ac))
	}
}

func Test_Sub_RegLit(t *testing.T) {
	cpu := cpu.NewCpu()
	mem := memory.NewMemory(255)
	packeds := [][]byte{
		MOV_LIT_R1.Pack(13),
		SUB_REG_LIT.Pack(uint16(register.R1.AsByte()), uint16(1)),
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
	if cpu.GetRegister(register.Ac) != 12 {
		t.Fatalf("error setting result value in accumulator, expected 12 got %d",
			cpu.GetRegister(register.Ac))
	}
}

func Test_Mul_RegReg(t *testing.T) {
	cpu := cpu.NewCpu()
	mem := memory.NewMemory(255)
	packeds := [][]byte{
		MOV_LIT_R1.Pack(13),
		MOV_LIT_R2.Pack(12),
		MUL_REG_REG.Pack(uint16(register.R1.AsByte()), uint16(register.R2.AsByte())),
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
	if cpu.GetRegister(register.Ac) != 156 {
		t.Fatalf("error setting result value in accumulator, expected 156 got %d",
			cpu.GetRegister(register.Ac))
	}
}

func Test_Mul_RegLit(t *testing.T) {
	cpu := cpu.NewCpu()
	mem := memory.NewMemory(255)
	packeds := [][]byte{
		MOV_LIT_R1.Pack(13),
		MUL_REG_LIT.Pack(uint16(register.R1.AsByte()), uint16(3)),
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
	if cpu.GetRegister(register.Ac) != 39 {
		t.Fatalf("error setting result value in accumulator, expected 39 got %d",
			cpu.GetRegister(register.Ac))
	}
}

func Test_Mul_RegReg_Overflow(t *testing.T) {
	cpu := cpu.NewCpu()
	mem := memory.NewMemory(255)
	var val uint16 = 257
	packeds := [][]byte{
		MOV_LIT_R1.Pack(val),
		MOV_LIT_R2.Pack(val),
		MUL_REG_REG.Pack(uint16(register.R1.AsByte()), uint16(register.R2.AsByte())),
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

	if cpu.GetRegister(register.R1) != val {
		t.Fatalf("error setting immediate value to register R1")
	}
	if cpu.GetRegister(register.R2) != val {
		t.Fatalf("error setting immediate value to register R2")
	}
	expected := uint16((int(val)*int(val))-0xffff) - 1
	if cpu.GetRegister(register.Ac) != expected {
		t.Fatalf("error setting result value in accumulator, expected %d got %d",
			expected, cpu.GetRegister(register.Ac))
	}
}

func Test_Div_RegReg(t *testing.T) {
	cpu := cpu.NewCpu()
	mem := memory.NewMemory(255)
	packeds := [][]byte{
		MOV_LIT_R1.Pack(130),
		MOV_LIT_R2.Pack(13),
		DIV_REG_REG.Pack(uint16(register.R1.AsByte()), uint16(register.R2.AsByte())),
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

	if cpu.GetRegister(register.R1) != 130 {
		t.Fatalf("error setting immediate value to register R1")
	}
	if cpu.GetRegister(register.R2) != 13 {
		t.Fatalf("error setting immediate value to register R2")
	}
	if cpu.GetRegister(register.Ac) != 10 {
		t.Fatalf("error setting result value in accumulator, expected 10 got %d",
			cpu.GetRegister(register.Ac))
	}
}

func Test_Div_RegLit(t *testing.T) {
	cpu := cpu.NewCpu()
	mem := memory.NewMemory(255)
	packeds := [][]byte{
		MOV_LIT_R1.Pack(39),
		DIV_REG_LIT.Pack(uint16(register.R1.AsByte()), uint16(3)),
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

	if cpu.GetRegister(register.R1) != 39 {
		t.Fatalf("error setting immediate value to register R1")
	}
	if cpu.GetRegister(register.Ac) != 13 {
		t.Fatalf("error setting result value in accumulator, expected 13 got %d",
			cpu.GetRegister(register.Ac))
	}
}

func Test_Mod_RegReg(t *testing.T) {
	cpu := cpu.NewCpu()
	mem := memory.NewMemory(255)
	packeds := [][]byte{
		MOV_LIT_R1.Pack(13),
		MOV_LIT_R2.Pack(12),
		MOD_REG_REG.Pack(uint16(register.R1.AsByte()), uint16(register.R2.AsByte())),
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
	if cpu.GetRegister(register.Ac) != 1 {
		t.Fatalf("error setting result value in accumulator, expected 1 got %d",
			cpu.GetRegister(register.Ac))
	}
}

func Test_Mod_RegLit(t *testing.T) {
	cpu := cpu.NewCpu()
	mem := memory.NewMemory(255)
	packeds := [][]byte{
		MOV_LIT_R1.Pack(40),
		MOD_REG_LIT.Pack(uint16(register.R1.AsByte()), uint16(3)),
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

	if cpu.GetRegister(register.R1) != 40 {
		t.Fatalf("error setting immediate value to register R1")
	}
	if cpu.GetRegister(register.Ac) != 1 {
		t.Fatalf("error setting result value in accumulator, expected 1 got %d",
			cpu.GetRegister(register.Ac))
	}
}
