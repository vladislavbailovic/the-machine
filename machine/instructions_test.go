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
		byte(instruction.MOV_LIT_R1), b1[0], b1[1],
		byte(instruction.MOV_LIT_R2), b2[0], b2[1],
		byte(instruction.ADD_REG_REG), register.R1.AsByte(), register.R2.AsByte(),
	})

	step := 0
	for step < 20 {
		if err := vm.Tick(); err != nil {
			t.Fatalf("error at tick %d: %v", step, err)
		}
		step++
	}
	vm.Debug()
}

func Test_Execute_Loop(t *testing.T) {
	vm := Machine{cpu: cpu.NewCpu(), memory: memory.NewMemory(255)}
	vm.LoadProgram(0, []byte{
		byte(instruction.MOV_LIT_AC), 0x01, 0x00,
		byte(instruction.MOV_LIT_R1), 0x01, 0x00,
		byte(instruction.ADD_REG_REG), register.Ac.AsByte(), register.R1.AsByte(),
		byte(instruction.JNE), 0x03, 0x00, 0x06, 0x00,
		byte(instruction.MOV_LIT_R2), 0xac, 0xab,
		byte(instruction.HALT),
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
}
