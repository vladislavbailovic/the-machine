package machine

import (
	"testing"
	"the-machine/machine/cpu"
	"the-machine/machine/instruction"
	"the-machine/machine/memory"
	"the-machine/machine/register"
)

func packSubroutine(instr ...[]byte) []byte {
	return packStatements(instruction.RET, instr...)
}

func Test_CalLit(t *testing.T) {
	subroutine := packSubroutine(
		instruction.MOV_REG_REG.Pack(register.R8.AsUint16(), register.Ac.AsUint16()), // reset Ac
		instruction.ADD_REG_LIT.Pack(register.Ac.AsUint16(), 14),
		instruction.MUL_REG_LIT.Pack(register.Ac.AsUint16(), 15),
		instruction.MUL_REG_LIT.Pack(register.Ac.AsUint16(), 5), // Ac = 1050
		instruction.MOV_LIT_MEM.Pack(13),
		instruction.ADD_REG_LIT.Pack(register.Ac.AsUint16(), 1), // Ac = 1051
		instruction.MOV_LIT_MEM.Pack(12),
	)
	main := packProgram(
		instruction.MOV_LIT_R1.Pack(13),
		instruction.MOV_LIT_R2.Pack(12),
		instruction.ADD_REG_LIT.Pack(register.Ac.AsUint16(), 15),
		instruction.MUL_REG_LIT.Pack(register.Ac.AsUint16(), 15), // Ac = 225
		instruction.CALL.Pack(register.Ac.AsUint16()),
	)

	vm := Machine{cpu: cpu.NewCpu(), memory: memory.NewMemory(2048)}
	vm.LoadProgram(225, subroutine)
	vm.LoadProgram(0, main)

	if _, err := run(vm); err != nil {
		vm.Debug()
		t.Fatalf("error running program: %v", err)
	}
	vm.Debug()

	sb1, err := vm.memory.GetByte(1050)
	if err != nil {
		t.Fatalf("expected subroutine memory set #1: %v", err)
	}
	if sb1 != 13 {
		t.Fatalf("expected subroutine to set memory at offset to be 13, got %d (%#02x/%016b)", sb1, sb1, sb1)
	}

	sb2, err := vm.memory.GetByte(1051)
	if err != nil {
		t.Fatalf("expected subroutine memory set #1: %v", err)
	}
	if sb2 != 12 {
		t.Fatalf("expected subroutine to set memory at offset to be 13, got %d (%#02x/%016b)", sb2, sb2, sb2)
	}
}
