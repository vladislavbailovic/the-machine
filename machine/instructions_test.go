package machine

import (
	"encoding/binary"
	"testing"
	"the-machine/machine/cpu"
	"the-machine/machine/instruction"
	"the-machine/machine/memory"
	"the-machine/machine/register"
)

func Test_Execute_Program(t *testing.T) {
	vm := Machine{cpu: cpu.NewCpu(), memory: memory.NewMemory(255)}
	b1 := make([]byte, 2)
	binary.LittleEndian.PutUint16(b1, uint16(1312))
	b2 := make([]byte, 2)
	binary.LittleEndian.PutUint16(b2, uint16(161))
	vm.LoadProgram(0, []byte{
		instruction.MOV_LIT_R1.AsByte(), b1[0], b1[1],
		instruction.MOV_LIT_R2.AsByte(), b2[0], b2[1],
		instruction.ADD_REG_REG.AsByte(), register.R1.AsByte(), register.R2.AsByte(),
	})

	step := 0
	for step < 20 {
		if err := vm.Tick(); err != nil {
			t.Fatalf("error at tick %d: %v", step, err)
		}
		step++
		if vm.IsDone() {
			break
		}
	}
	vm.Debug()
	if val, err := vm.cpu.GetRegister(register.Ac); err != nil || val != 1312+161 {
		t.Fatalf("expected %d in Ac, got %#02x (%d) instead, err: %v", 1312+161, val, val, err)
	}
	if 20 != step {
		t.Fatalf("expected nohalt program to run to external limit, break after %d", step)
	}
}

func Test_Execute_Loop(t *testing.T) {
	vm := Machine{cpu: cpu.NewCpu(), memory: memory.NewMemory(255)}
	vm.LoadProgram(0, []byte{
		instruction.MOV_LIT_AC.AsByte(), 0x01, 0x00,
		instruction.MOV_LIT_R1.AsByte(), 0x01, 0x00,
		instruction.ADD_REG_REG.AsByte(), register.Ac.AsByte(), register.R1.AsByte(),
		instruction.JNE.AsByte(), 0x03, 0x00, 0x06, 0x00,
		instruction.MOV_LIT_R2.AsByte(), 0xab, 0xac,
		instruction.HALT.AsByte(),
	})

	step := 0
	for step < 20 {
		vm.Debug()
		if err := vm.Tick(); err != nil {
			t.Fatalf("error at tick %d: %v", step, err)
		}
		step++
		if vm.IsDone() {
			break
		}
	}

	if val, err := vm.cpu.GetRegister(register.R2); err != nil || val != 0xacab {
		t.Fatalf("expected %#02x in R2, got %#02x instead, err: %v", 0xacab, val, err)
	}
	if val, err := vm.cpu.GetRegister(register.Ac); err != nil || val != 3 {
		t.Fatalf("expected 3 in Ac, got %#02x instead, err: %v", val, err)
	}
	if step != 8 {
		t.Fatalf("expected exactly %d steps but got %d", 8, step)
	}
}

func Test_CopyToMemory(t *testing.T) {
	vm := Machine{cpu: cpu.NewCpu(), memory: memory.NewMemory(255)}
	address := make([]byte, 2)
	binary.LittleEndian.PutUint16(address, uint16(161))
	vm.LoadProgram(0, []byte{
		instruction.MOV_LIT_R1.AsByte(), 0x12, 0x00,
		instruction.MOV_REG_MEM.AsByte(), register.R1.AsByte(), 0x13, 0x00,
		instruction.MOV_LIT_MEM.AsByte(), 0xab, 0xac, address[0], address[1],
		instruction.HALT.AsByte(),
	})

	step := 0
	for step < 20 {
		vm.Debug()
		if err := vm.Tick(); err != nil {
			t.Fatalf("error at tick %d: %v", step, err)
		}
		step++
		if vm.IsDone() {
			break
		}
	}

	if val, err := vm.memory.GetUint16(memory.Address(0x13)); err != nil || val != 0x12 {
		t.Fatalf("expected %#02x (%d) at address %#02x (%d), got %#02x (%d), err %v",
			0x12, 0x12, 0x13, 0x13, val, val, err)
	}

	if val, err := vm.memory.GetUint16(memory.Address(161)); err != nil || val != 0xacab {
		t.Fatalf("expected %#02x (%d) at address %#02x (%d), got %#02x (%d), err %v",
			0xacab, 0xacab, 161, 161, val, val, err)
	}

	if step != 4 {
		t.Fatalf("expected exactly %d steps but got %d", 4, step)
	}
}
