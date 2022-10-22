package machine

import (
	"fmt"
	"testing"
	"the-machine/machine/cpu"
	"the-machine/machine/instruction"
	"the-machine/machine/memory"
	"the-machine/machine/register"
)

func run(vm Machine) (int, error) {
	step := 0
	for step < 127 {
		if err := vm.Tick(); err != nil {
			return step, fmt.Errorf("error at tick %d: %v", step, err)
		}
		step++
		if vm.IsDone() {
			break
		}
	}
	return step, nil
}

func packInstruction(kind instruction.Type, value uint16) []byte {
	inst1 := value | (uint16(kind.AsByte()) << 10)
	// fmt.Printf("val:%016b\nins:%016b\nbot:%016b\n", value, uint16(kind), inst1)
	// fmt.Printf("shl:%016b\nshr:%016b\n", (inst1 >> 8), (inst1 << 8))
	return []byte{
		byte(inst1),
		byte(inst1 >> 8),
	}
}

func packProgram(instr ...[]byte) []byte {
	res := make([]byte, 0, len(instr)+2)
	for _, b := range instr {
		res = append(res, b...)
	}
	halt := packInstruction(instruction.HALT, 0)
	res = append(res, halt...)
	return res
}

func Test_SingleUint16Instruction_All(t *testing.T) {
	values := []uint16{161, 13, 12, 255, 512, 1023}
	registers := map[register.Register]instruction.Type{
		register.R1: instruction.MOV_LIT_R1,
		register.R2: instruction.MOV_LIT_R2,
		register.R3: instruction.MOV_LIT_R3,
		register.R4: instruction.MOV_LIT_R4,
	}

	for idx, value := range values {
		regIdx := 0
		for reg, instr := range registers {
			// fmt.Printf("--- %d::%d: %d into %v ---\n", idx, regIdx, value, reg)
			vm := Machine{cpu: cpu.NewCpu(), memory: memory.NewMemory(255)}
			program := packInstruction(instr, value)
			vm.LoadProgram(0, packProgram(program))
			if step, err := run(vm); err != nil || step > 2 {
				t.Fatalf("error running machine or machine stuck: step %d, error: %v", step, err)
			}
			if vm.cpu.GetRegister(reg) != value {
				vm.Debug()
				t.Fatalf("%d: error setting immediate value %d to register %v: %d",
					idx, value, reg, vm.cpu.GetRegister(reg))
			}
			regIdx++
		}
	}

}

/*
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
	if val := vm.cpu.GetRegister(register.Ac); val != 1312+161 {
		t.Fatalf("expected %d in Ac, got %#02x (%d) instead", 1312+161, val, val)
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

	if val := vm.cpu.GetRegister(register.R2); val != 0xacab {
		t.Fatalf("expected %#02x in R2, got %#02x instead", 0xacab, val)
	}
	if val := vm.cpu.GetRegister(register.Ac); val != 3 {
		t.Fatalf("expected 3 in Ac, got %#02x instead", val)
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
*/
