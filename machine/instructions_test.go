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
	return []byte{
		byte(inst1),
		byte(inst1 >> 8),
	}
}

func packInstruction2(kind instruction.Type, value1 uint16, value2 uint16) []byte {
	v1 := value1 << 12 & 0b1111_000000000000
	v2 := value2 << 8 & 0b0000_1111_0000000
	// fmt.Printf("v1: %016b\n", v1)
	// fmt.Printf("v2: %016b\n", v2)
	value := v1 | v2
	return packInstruction(kind, value)
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

func Test_AddRegReg_One(t *testing.T) {
	vm := Machine{cpu: cpu.NewCpu(), memory: memory.NewMemory(255)}
	program := packProgram(
		packInstruction(instruction.MOV_LIT_R1, 13),
		packInstruction(instruction.MOV_LIT_R2, 12),
		packInstruction2(instruction.ADD_REG_REG, uint16(register.R1.AsByte()), uint16(register.R2.AsByte())),
	)
	vm.LoadProgram(0, program)

	if vm.cpu.GetRegister(register.Ac) != 0 {
		vm.Debug()
		t.Fatalf("machine initial state error: expected empty Ac, got: %d", vm.cpu.GetRegister(register.Ac))
	}

	if step, err := run(vm); err != nil || step > 4 {
		t.Fatalf("error running machine or machine stuck: step %d, error: %v", step, err)
	}

	if vm.cpu.GetRegister(register.R1) != 13 {
		vm.Debug()
		t.Fatalf("error setting immediate value to register R1")
	}
	if vm.cpu.GetRegister(register.R2) != 12 {
		vm.Debug()
		t.Fatalf("error setting immediate value to register R2")
	}
	if vm.cpu.GetRegister(register.Ac) != 25 {
		vm.Debug()
		t.Fatalf("error adding R1 and R2: expected 25, got: %d", vm.cpu.GetRegister(register.Ac))
	}
}

func Test_AddRegLit_One(t *testing.T) {
	vm := Machine{cpu: cpu.NewCpu(), memory: memory.NewMemory(255)}
	program := packProgram(
		packInstruction(instruction.MOV_LIT_R1, 13),
		packInstruction2(instruction.ADD_REG_LIT, uint16(register.R1.AsByte()), uint16(3)), // 2 bytes left for literal
	)
	vm.LoadProgram(0, program)

	if vm.cpu.GetRegister(register.Ac) != 0 {
		vm.Debug()
		t.Fatalf("machine initial state error: expected empty Ac, got: %d", vm.cpu.GetRegister(register.Ac))
	}

	if step, err := run(vm); err != nil || step > 3 {
		t.Fatalf("error running machine or machine stuck: step %d, error: %v", step, err)
	}

	if vm.cpu.GetRegister(register.R1) != 13 {
		vm.Debug()
		t.Fatalf("error setting immediate value to register R1")
	}
	if vm.cpu.GetRegister(register.Ac) != 16 {
		vm.Debug()
		t.Fatalf("error adding R1 with lit 3: expected 16, got: %d", vm.cpu.GetRegister(register.Ac))
	}
}

func Test_SubRegReg_One(t *testing.T) {
	vm := Machine{cpu: cpu.NewCpu(), memory: memory.NewMemory(255)}
	program := packProgram(
		packInstruction(instruction.MOV_LIT_R1, 13),
		packInstruction(instruction.MOV_LIT_R2, 12),
		packInstruction2(instruction.SUB_REG_REG, uint16(register.R1.AsByte()), uint16(register.R2.AsByte())),
	)
	vm.LoadProgram(0, program)

	if vm.cpu.GetRegister(register.Ac) != 0 {
		vm.Debug()
		t.Fatalf("machine initial state error: expected empty Ac, got: %d", vm.cpu.GetRegister(register.Ac))
	}

	if step, err := run(vm); err != nil || step > 4 {
		t.Fatalf("error running machine or machine stuck: step %d, error: %v", step, err)
	}

	if vm.cpu.GetRegister(register.R1) != 13 {
		vm.Debug()
		t.Fatalf("error setting immediate value to register R1")
	}
	if vm.cpu.GetRegister(register.R2) != 12 {
		vm.Debug()
		t.Fatalf("error setting immediate value to register R2")
	}
	if vm.cpu.GetRegister(register.Ac) != 1 {
		vm.Debug()
		t.Fatalf("error subtracting R1 and R2: expected 1, got: %d", vm.cpu.GetRegister(register.Ac))
	}
}

func Test_SubRegLit_One(t *testing.T) {
	vm := Machine{cpu: cpu.NewCpu(), memory: memory.NewMemory(255)}
	program := packProgram(
		packInstruction(instruction.MOV_LIT_R1, 13),
		packInstruction2(instruction.SUB_REG_LIT, uint16(register.R1.AsByte()), uint16(3)), // 2 bytes left for literal
	)
	vm.LoadProgram(0, program)

	if vm.cpu.GetRegister(register.Ac) != 0 {
		vm.Debug()
		t.Fatalf("machine initial state error: expected empty Ac, got: %d", vm.cpu.GetRegister(register.Ac))
	}

	if step, err := run(vm); err != nil || step > 3 {
		t.Fatalf("error running machine or machine stuck: step %d, error: %v", step, err)
	}

	if vm.cpu.GetRegister(register.R1) != 13 {
		vm.Debug()
		t.Fatalf("error setting immediate value to register R1")
	}
	if vm.cpu.GetRegister(register.Ac) != 10 {
		vm.Debug()
		t.Fatalf("error subtracting R1 and lit 3: expected 10, got: %d", vm.cpu.GetRegister(register.Ac))
	}
}

func Test_MulRegReg_One(t *testing.T) {
	vm := Machine{cpu: cpu.NewCpu(), memory: memory.NewMemory(255)}
	program := packProgram(
		packInstruction(instruction.MOV_LIT_R1, 13),
		packInstruction(instruction.MOV_LIT_R2, 12),
		packInstruction2(instruction.MUL_REG_REG, uint16(register.R1.AsByte()), uint16(register.R2.AsByte())),
	)
	vm.LoadProgram(0, program)

	if vm.cpu.GetRegister(register.Ac) != 0 {
		vm.Debug()
		t.Fatalf("machine initial state error: expected empty Ac, got: %d", vm.cpu.GetRegister(register.Ac))
	}

	if step, err := run(vm); err != nil || step > 4 {
		t.Fatalf("error running machine or machine stuck: step %d, error: %v", step, err)
	}

	if vm.cpu.GetRegister(register.R1) != 13 {
		vm.Debug()
		t.Fatalf("error setting immediate value to register R1")
	}
	if vm.cpu.GetRegister(register.R2) != 12 {
		vm.Debug()
		t.Fatalf("error setting immediate value to register R2")
	}
	if vm.cpu.GetRegister(register.Ac) != 156 {
		vm.Debug()
		t.Fatalf("error multiplying R1 and R2: expected 156, got: %d", vm.cpu.GetRegister(register.Ac))
	}
}

func Test_MulRegLit_One(t *testing.T) {
	vm := Machine{cpu: cpu.NewCpu(), memory: memory.NewMemory(255)}
	program := packProgram(
		packInstruction(instruction.MOV_LIT_R1, 13),
		packInstruction2(instruction.MUL_REG_LIT, uint16(register.R1.AsByte()), uint16(3)), // 2 bytes left for literal
	)
	vm.LoadProgram(0, program)

	if vm.cpu.GetRegister(register.Ac) != 0 {
		vm.Debug()
		t.Fatalf("machine initial state error: expected empty Ac, got: %d", vm.cpu.GetRegister(register.Ac))
	}

	if step, err := run(vm); err != nil || step > 3 {
		t.Fatalf("error running machine or machine stuck: step %d, error: %v", step, err)
	}

	if vm.cpu.GetRegister(register.R1) != 13 {
		vm.Debug()
		t.Fatalf("error setting immediate value to register R1")
	}
	if vm.cpu.GetRegister(register.Ac) != 39 {
		vm.Debug()
		t.Fatalf("error multiplying R1 and lit 3: expected 39, got: %d", vm.cpu.GetRegister(register.Ac))
	}
}

func Test_MulRegReg_Overflow(t *testing.T) {
	vm := Machine{cpu: cpu.NewCpu(), memory: memory.NewMemory(255)}
	var val uint16 = 257
	program := packProgram(
		packInstruction(instruction.MOV_LIT_R1, val),
		packInstruction(instruction.MOV_LIT_R2, val),
		packInstruction2(instruction.MUL_REG_REG, uint16(register.R1.AsByte()), uint16(register.R2.AsByte())),
	)
	vm.LoadProgram(0, program)

	if vm.cpu.GetRegister(register.Ac) != 0 {
		vm.Debug()
		t.Fatalf("machine initial state error: expected empty Ac, got: %d", vm.cpu.GetRegister(register.Ac))
	}

	if step, err := run(vm); err != nil || step > 4 {
		t.Fatalf("error running machine or machine stuck: step %d, error: %v", step, err)
	}

	if vm.cpu.GetRegister(register.R1) != val {
		vm.Debug()
		t.Fatalf("error setting immediate value to register R1")
	}
	if vm.cpu.GetRegister(register.R2) != val {
		vm.Debug()
		t.Fatalf("error setting immediate value to register R2")
	}

	expected := uint16((int(val)*int(val))-0b1111_1111_1111_1111) - 1
	if vm.cpu.GetRegister(register.Ac) != expected {
		vm.Debug()
		t.Fatalf("error multiplying with overflow R1 and R2: expected %d, got: %d",
			expected, vm.cpu.GetRegister(register.Ac))
	}
}

func Test_DivRegReg_Straight(t *testing.T) {
	vm := Machine{cpu: cpu.NewCpu(), memory: memory.NewMemory(255)}
	program := packProgram(
		packInstruction(instruction.MOV_LIT_R1, 120),
		packInstruction(instruction.MOV_LIT_R2, 12),
		packInstruction2(instruction.DIV_REG_REG, uint16(register.R1.AsByte()), uint16(register.R2.AsByte())),
	)
	vm.LoadProgram(0, program)

	if vm.cpu.GetRegister(register.Ac) != 0 {
		vm.Debug()
		t.Fatalf("machine initial state error: expected empty Ac, got: %d", vm.cpu.GetRegister(register.Ac))
	}

	if step, err := run(vm); err != nil || step > 4 {
		t.Fatalf("error running machine or machine stuck: step %d, error: %v", step, err)
	}

	if vm.cpu.GetRegister(register.R1) != 120 {
		vm.Debug()
		t.Fatalf("error setting immediate value to register R1")
	}
	if vm.cpu.GetRegister(register.R2) != 12 {
		vm.Debug()
		t.Fatalf("error setting immediate value to register R2")
	}
	if vm.cpu.GetRegister(register.Ac) != 10 {
		vm.Debug()
		t.Fatalf("error straight dividing R1 and R2: expected 10, got: %d", vm.cpu.GetRegister(register.Ac))
	}
}

func Test_DivRegLit_Straight(t *testing.T) {
	vm := Machine{cpu: cpu.NewCpu(), memory: memory.NewMemory(255)}
	program := packProgram(
		packInstruction(instruction.MOV_LIT_R1, 39),
		packInstruction2(instruction.DIV_REG_LIT, uint16(register.R1.AsByte()), uint16(3)), // 2 bytes left for literal
	)
	vm.LoadProgram(0, program)

	if vm.cpu.GetRegister(register.Ac) != 0 {
		vm.Debug()
		t.Fatalf("machine initial state error: expected empty Ac, got: %d", vm.cpu.GetRegister(register.Ac))
	}

	if step, err := run(vm); err != nil || step > 3 {
		t.Fatalf("error running machine or machine stuck: step %d, error: %v", step, err)
	}

	if vm.cpu.GetRegister(register.R1) != 39 {
		vm.Debug()
		t.Fatalf("error setting immediate value to register R1")
	}
	if vm.cpu.GetRegister(register.Ac) != 13 {
		vm.Debug()
		t.Fatalf("error dividing R1 and lit 3: expected 13, got: %d", vm.cpu.GetRegister(register.Ac))
	}
}

func Test_DivRegReg_Round(t *testing.T) {
	vm := Machine{cpu: cpu.NewCpu(), memory: memory.NewMemory(255)}
	program := packProgram(
		packInstruction(instruction.MOV_LIT_R1, 128),
		packInstruction(instruction.MOV_LIT_R2, 12),
		packInstruction2(instruction.DIV_REG_REG, uint16(register.R1.AsByte()), uint16(register.R2.AsByte())),
	)
	vm.LoadProgram(0, program)

	if vm.cpu.GetRegister(register.Ac) != 0 {
		vm.Debug()
		t.Fatalf("machine initial state error: expected empty Ac, got: %d", vm.cpu.GetRegister(register.Ac))
	}

	if step, err := run(vm); err != nil || step > 4 {
		t.Fatalf("error running machine or machine stuck: step %d, error: %v", step, err)
	}

	if vm.cpu.GetRegister(register.R1) != 128 {
		vm.Debug()
		t.Fatalf("error setting immediate value to register R1")
	}
	if vm.cpu.GetRegister(register.R2) != 12 {
		vm.Debug()
		t.Fatalf("error setting immediate value to register R2")
	}
	if vm.cpu.GetRegister(register.Ac) != 10 {
		vm.Debug()
		t.Fatalf("error dividing R1 and R2 with rounding down: expected 10, got: %d", vm.cpu.GetRegister(register.Ac))
	}
}

func Test_DivRegLit_Round(t *testing.T) {
	vm := Machine{cpu: cpu.NewCpu(), memory: memory.NewMemory(255)}
	program := packProgram(
		packInstruction(instruction.MOV_LIT_R1, 40),
		packInstruction2(instruction.DIV_REG_LIT, uint16(register.R1.AsByte()), uint16(3)), // 2 bytes left for literal
	)
	vm.LoadProgram(0, program)

	if vm.cpu.GetRegister(register.Ac) != 0 {
		vm.Debug()
		t.Fatalf("machine initial state error: expected empty Ac, got: %d", vm.cpu.GetRegister(register.Ac))
	}

	if step, err := run(vm); err != nil || step > 3 {
		t.Fatalf("error running machine or machine stuck: step %d, error: %v", step, err)
	}

	if vm.cpu.GetRegister(register.R1) != 40 {
		vm.Debug()
		t.Fatalf("error setting immediate value to register R1")
	}
	if vm.cpu.GetRegister(register.Ac) != 13 {
		vm.Debug()
		t.Fatalf("error dividing with rounding R1 and lit 3: expected 13, got: %d", vm.cpu.GetRegister(register.Ac))
	}
}

func Test_ModRegReg_One(t *testing.T) {
	vm := Machine{cpu: cpu.NewCpu(), memory: memory.NewMemory(255)}
	program := packProgram(
		packInstruction(instruction.MOV_LIT_R1, 13),
		packInstruction(instruction.MOV_LIT_R2, 12),
		packInstruction2(instruction.MOD_REG_REG, uint16(register.R1.AsByte()), uint16(register.R2.AsByte())),
	)
	vm.LoadProgram(0, program)

	if vm.cpu.GetRegister(register.Ac) != 0 {
		vm.Debug()
		t.Fatalf("machine initial state error: expected empty Ac, got: %d", vm.cpu.GetRegister(register.Ac))
	}

	if step, err := run(vm); err != nil || step > 4 {
		t.Fatalf("error running machine or machine stuck: step %d, error: %v", step, err)
	}

	if vm.cpu.GetRegister(register.R1) != 13 {
		vm.Debug()
		t.Fatalf("error setting immediate value to register R1")
	}
	if vm.cpu.GetRegister(register.R2) != 12 {
		vm.Debug()
		t.Fatalf("error setting immediate value to register R2")
	}
	if vm.cpu.GetRegister(register.Ac) != 1 {
		vm.Debug()
		t.Fatalf("error modding R1 and R2: expected 1, got: %d", vm.cpu.GetRegister(register.Ac))
	}
}

func Test_ModRegLit_One(t *testing.T) {
	vm := Machine{cpu: cpu.NewCpu(), memory: memory.NewMemory(255)}
	program := packProgram(
		packInstruction(instruction.MOV_LIT_R1, 40),
		packInstruction2(instruction.MOD_REG_LIT, uint16(register.R1.AsByte()), uint16(3)), // 2 bytes left for literal
	)
	vm.LoadProgram(0, program)

	if vm.cpu.GetRegister(register.Ac) != 0 {
		vm.Debug()
		t.Fatalf("machine initial state error: expected empty Ac, got: %d", vm.cpu.GetRegister(register.Ac))
	}

	if step, err := run(vm); err != nil || step > 3 {
		t.Fatalf("error running machine or machine stuck: step %d, error: %v", step, err)
	}

	if vm.cpu.GetRegister(register.R1) != 40 {
		vm.Debug()
		t.Fatalf("error setting immediate value to register R1")
	}
	if vm.cpu.GetRegister(register.Ac) != 1 {
		vm.Debug()
		t.Fatalf("error modding R1 and lit 3: expected 13, got: %d", vm.cpu.GetRegister(register.Ac))
	}
}

func Test_MovRegMem(t *testing.T) {
	vm := Machine{cpu: cpu.NewCpu(), memory: memory.NewMemory(2048)}
	program := packProgram(
		packInstruction(instruction.MOV_LIT_R1, 1023),
		packInstruction(instruction.MOV_LIT_R2, 289),
		packInstruction2(instruction.ADD_REG_REG, uint16(register.R1.AsByte()), uint16(register.R2.AsByte())),
		packInstruction(instruction.MOV_LIT_R3, 161),
		packInstruction(instruction.MOV_REG_MEM, uint16(register.R3.AsByte())),
	)
	vm.LoadProgram(0, program)

	if vm.cpu.GetRegister(register.Ac) != 0 {
		vm.Debug()
		t.Fatalf("machine initial state error: expected empty Ac, got: %d", vm.cpu.GetRegister(register.Ac))
	}

	if step, err := run(vm); err != nil || step > 6 {
		t.Fatalf("error running machine or machine stuck: step %d, error: %v", step, err)
	}

	if vm.cpu.GetRegister(register.R1) != 1023 {
		vm.Debug()
		t.Fatalf("error setting immediate value to register R1")
	}
	if vm.cpu.GetRegister(register.R2) != 289 {
		vm.Debug()
		t.Fatalf("error setting immediate value to register R2")
	}
	if vm.cpu.GetRegister(register.R3) != 161 {
		vm.Debug()
		t.Fatalf("error setting immediate value to register R3")
	}
	if vm.cpu.GetRegister(register.Ac) != 1312 {
		vm.Debug()
		t.Fatalf("error setting memory address in accumulator")
	}

	res, err := vm.memory.GetUint16(1312)
	if err != nil {
		vm.Debug()
		t.Fatalf("error setting memory at 1312: %v", err)
	}
	if res != 161 {
		vm.Debug()
		t.Fatalf("error setting memory at 1312: expected 161, got: %d", res)
	}
}

func Test_MovLitMem(t *testing.T) {
	vm := Machine{cpu: cpu.NewCpu(), memory: memory.NewMemory(2048)}
	program := packProgram(
		packInstruction(instruction.MOV_LIT_R1, 1023),
		packInstruction(instruction.MOV_LIT_R2, 289),
		packInstruction2(instruction.ADD_REG_REG, uint16(register.R1.AsByte()), uint16(register.R2.AsByte())),
		packInstruction(instruction.MOV_LIT_MEM, uint16(161)),
	)
	vm.LoadProgram(0, program)

	if vm.cpu.GetRegister(register.Ac) != 0 {
		vm.Debug()
		t.Fatalf("machine initial state error: expected empty Ac, got: %d", vm.cpu.GetRegister(register.Ac))
	}

	if step, err := run(vm); err != nil || step > 6 {
		t.Fatalf("error running machine or machine stuck: step %d, error: %v", step, err)
	}

	if vm.cpu.GetRegister(register.R1) != 1023 {
		vm.Debug()
		t.Fatalf("error setting immediate value to register R1")
	}
	if vm.cpu.GetRegister(register.R2) != 289 {
		vm.Debug()
		t.Fatalf("error setting immediate value to register R2")
	}
	if vm.cpu.GetRegister(register.Ac) != 1312 {
		vm.Debug()
		t.Fatalf("error setting memory address in accumulator")
	}

	res, err := vm.memory.GetUint16(1312)
	if err != nil {
		vm.Debug()
		t.Fatalf("error setting memory at 1312: %v", err)
	}
	if res != 161 {
		vm.Debug()
		t.Fatalf("error setting memory at 1312: expected 161, got: %d", res)
	}
}
