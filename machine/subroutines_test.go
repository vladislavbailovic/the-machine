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

func Test_CalReg(t *testing.T) {
	subroutine := packSubroutine(
		instruction.MOV_LIT_R1.Pack(3),
		instruction.MOV_LIT_R2.Pack(2),
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
		instruction.MOV_LIT_R3.Pack(161),
	)

	vm := Machine{cpu: cpu.NewCpu(), memory: memory.NewMemory(2048)}
	vm.LoadProgram(225, subroutine)
	vm.LoadProgram(0, main)

	if steps, err := run(vm); err != nil || steps != 17 {
		vm.Debug()
		t.Fatalf("machine stuck (%d) or error running program: %v", steps, err)
	}

	// Ensure memory actually being set in subroutine:

	sb1, err := vm.memory.GetByte(1050)
	if err != nil {
		vm.Debug()
		t.Fatalf("expected subroutine memory set #1: %v", err)
	}
	if sb1 != 13 {
		vm.Debug()
		t.Fatalf("expected subroutine to set memory at offset to be 13, got %d (%#02x/%016b)", sb1, sb1, sb1)
	}

	sb2, err := vm.memory.GetByte(1051)
	if err != nil {
		vm.Debug()
		t.Fatalf("expected subroutine memory set #1: %v", err)
	}
	if sb2 != 12 {
		vm.Debug()
		t.Fatalf("expected subroutine to set memory at offset to be 13, got %d (%#02x/%016b)", sb2, sb2, sb2)
	}

	// Ensure registers state being restored after call:

	if vm.cpu.GetRegister(register.R1) != 13 {
		vm.Debug()
		t.Fatalf("expected register R1 to be preserved on return")
	}
	if vm.cpu.GetRegister(register.R2) != 12 {
		vm.Debug()
		t.Fatalf("expected register R2 to be preserved on return")
	}

	// Ensure execution continuing after call:

	if vm.cpu.GetRegister(register.R3) != 161 {
		vm.Debug()
		t.Fatalf("expected execution to continue where it left off")
	}
}
