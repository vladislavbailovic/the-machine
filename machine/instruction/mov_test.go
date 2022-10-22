package instruction

import (
	"testing"
	"the-machine/machine/cpu"
	"the-machine/machine/memory"
	"the-machine/machine/register"
)

func unpackInstruction(packed []byte) (Instruction, uint16) {
	mem := memory.NewMemory(2)
	mem.SetByte(0, packed[0])
	mem.SetByte(1, packed[1])
	raw, _ := mem.GetUint16(0)
	kind, decoded := Decode(raw)
	instruction := Descriptors[kind]

	return instruction, decoded
}

func runPackedInstructionWithCpu(packed []byte, cpu *cpu.Cpu) (*memory.Memory, error) {
	mem := memory.NewMemory(2)
	mem.SetByte(0, packed[0])
	mem.SetByte(1, packed[1])
	raw, _ := mem.GetUint16(0)
	kind, decoded := Decode(raw)
	instruction := Descriptors[kind]
	err := instruction.Executor.Execute(decoded, cpu, mem)

	return mem, err
}

func runPackedInstruction(packed []byte) (*cpu.Cpu, *memory.Memory, error) {
	cpu := cpu.NewCpu()
	mem, err := runPackedInstructionWithCpu(packed, cpu)
	return cpu, mem, err
}

func Test_MovLitReg_All(t *testing.T) {
	values := []uint16{161, 13, 12, 255, 512, 1023}
	registers := map[register.Register]Type{
		register.R1: MOV_LIT_R1,
		register.R2: MOV_LIT_R2,
		register.R3: MOV_LIT_R3,
		register.R4: MOV_LIT_R4,
		register.R7: MOV_LIT_R7,
	}

	for idx, value := range values {
		regIdx := 0
		for reg, instr := range registers {
			cpu, _, err := runPackedInstruction(instr.Pack(value))

			if err != nil {
				t.Fatalf("%d: error executing instruction %v with value %d: %v",
					idx, instr, value, err)
			}

			if cpu.GetRegister(reg) != value {
				t.Fatalf("%d: error setting immediate value %d to register %v: %d",
					idx, value, reg, cpu.GetRegister(reg))
			}
			regIdx++
		}
	}
}

func Test_MovRegMem(t *testing.T) {
	cpu := cpu.NewCpu()
	mem := memory.NewMemory(2048)
	packeds := [][]byte{
		MOV_LIT_R1.Pack(1023),
		MOV_LIT_R2.Pack(289),
		ADD_REG_REG.Pack(uint16(register.R1.AsByte()), uint16(register.R2.AsByte())),
		MOV_LIT_R3.Pack(161),
		MOV_REG_MEM.Pack(uint16(register.R3.AsByte())),
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

	if cpu.GetRegister(register.R1) != 1023 {
		t.Fatalf("error setting immediate value to register R1")
	}
	if cpu.GetRegister(register.R2) != 289 {
		t.Fatalf("error setting immediate value to register R2")
	}
	if cpu.GetRegister(register.R3) != 161 {
		t.Fatalf("error setting immediate value to register R3")
	}
	if cpu.GetRegister(register.Ac) != 1312 {
		t.Fatalf("error setting memory address in accumulator")
	}

	res, err := mem.GetUint16(1312)
	if err != nil {
		t.Fatalf("error setting memory at 1312: %v", err)
	}
	if res != 161 {
		t.Fatalf("error setting memory at 1312: expected 161, got: %d", res)
	}
}

func Test_MovLitMem(t *testing.T) {
	cpu := cpu.NewCpu()
	mem := memory.NewMemory(2048)
	packeds := [][]byte{
		MOV_LIT_R1.Pack(1023),
		MOV_LIT_R2.Pack(289),
		ADD_REG_REG.Pack(uint16(register.R1.AsByte()), uint16(register.R2.AsByte())),
		MOV_LIT_MEM.Pack(uint16(161)),
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

	if cpu.GetRegister(register.R1) != 1023 {
		t.Fatalf("error setting immediate value to register R1")
	}
	if cpu.GetRegister(register.R2) != 289 {
		t.Fatalf("error setting immediate value to register R2")
	}
	if cpu.GetRegister(register.Ac) != 1312 {
		t.Fatalf("error setting memory address in accumulator")
	}

	res, err := mem.GetUint16(1312)
	if err != nil {
		t.Fatalf("error setting memory at 1312: %v", err)
	}
	if res != 161 {
		t.Fatalf("error setting memory at 1312: expected 161, got: %d", res)
	}
}

func Test_MovRegReg_GeneralPurpose(t *testing.T) {
	cpu := cpu.NewCpu()
	mem := memory.NewMemory(255)
	packeds := [][]byte{
		MOV_LIT_R1.Pack(161),
		MOV_REG_REG.Pack(uint16(register.R1.AsByte()), uint16(register.R2.AsByte())),
	}

	for idx, packed := range packeds {
		instr, raw := unpackInstruction(packed)
		if err := instr.Executor.Execute(raw, cpu, mem); err != nil {
			t.Fatalf("%d: error executing instruction %v: %err", idx, instr, err)
		}
	}

	if cpu.GetRegister(register.R1) != 161 {
		t.Fatalf("error setting immediate value to register R1")
	}
	if cpu.GetRegister(register.R2) != 161 {
		t.Fatalf("error copying R1 value to R2: %d", cpu.GetRegister(register.R2))
	}
}

func Test_MovRegReg_Ac2General(t *testing.T) {
	cpu := cpu.NewCpu()
	mem := memory.NewMemory(255)
	packeds := [][]byte{
		MOV_LIT_R1.Pack(160),
		ADD_REG_LIT.Pack(uint16(register.R1.AsByte()), uint16(1)),
		MOV_REG_REG.Pack(uint16(register.Ac.AsByte()), uint16(register.R2.AsByte())),
	}

	for idx, packed := range packeds {
		instr, raw := unpackInstruction(packed)
		if err := instr.Executor.Execute(raw, cpu, mem); err != nil {
			t.Fatalf("%d: error executing instruction %v: %err", idx, instr, err)
		}
	}

	if cpu.GetRegister(register.R1) != 160 {
		t.Fatalf("error setting immediate value to register R1")
	}
	if cpu.GetRegister(register.Ac) != 161 {
		t.Fatalf("error in Accumulator: %d", cpu.GetRegister(register.Ac))
	}
	if cpu.GetRegister(register.R2) != 161 {
		t.Fatalf("error copying Accumulator to R2: %d", cpu.GetRegister(register.R2))
	}
}
