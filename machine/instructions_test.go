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

func packProgram(instr ...[]byte) []byte {
	res := make([]byte, 0, len(instr)+2)
	for _, b := range instr {
		res = append(res, b...)
	}
	halt := instruction.HALT.Pack(0)
	res = append(res, halt...)
	return res
}

func Test_Pack2Regs(t *testing.T) {
	values := []uint16{161, 13, 12, 255, 512, 1023}
	registers := []register.Register{
		register.R2,
		register.R3,
		register.R4,
		register.R5,
	}
	for vid, value := range values {
		for rid, destination := range registers {
			vm := Machine{cpu: cpu.NewCpu(), memory: memory.NewMemory(255)}
			program := packProgram(
				instruction.MOV_LIT_R7.Pack(value),
				instruction.MOV_REG_REG.Pack(uint16(register.R7.AsByte()), uint16(destination.AsByte())),
			)
			vm.LoadProgram(0, program)
			run(vm)
			if vm.cpu.GetRegister(register.R7) != value {
				vm.Debug()
				t.Fatalf("%d::%d: Invalid value in source register", vid, rid)
			}
			if vm.cpu.GetRegister(destination) != value {
				vm.Debug()
				t.Fatalf("%d::%d: Invalid value in destination register %v: %d",
					vid, rid, destination, vm.cpu.GetRegister(destination))
			}
		}
	}
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
			program := instr.Pack(value)
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
		instruction.MOV_LIT_R1.Pack(13),
		instruction.MOV_LIT_R2.Pack(12),
		instruction.ADD_REG_REG.Pack(uint16(register.R1.AsByte()), uint16(register.R2.AsByte())),
	)
	vm.LoadProgram(0, program)

	if vm.cpu.GetRegister(register.Ac) != 0 {
		vm.Debug()
		t.Fatalf("machine initial state error: expected empty Ac, got: %d", vm.cpu.GetRegister(register.Ac))
	}

	if step, err := run(vm); err != nil || step > 4 {
		vm.Debug()
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
		instruction.MOV_LIT_R1.Pack(13),
		instruction.ADD_REG_LIT.Pack(uint16(register.R1.AsByte()), uint16(3)), // 2 bytes left for literal
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
		instruction.MOV_LIT_R1.Pack(13),
		instruction.MOV_LIT_R2.Pack(12),
		instruction.SUB_REG_REG.Pack(uint16(register.R1.AsByte()), uint16(register.R2.AsByte())),
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
		instruction.MOV_LIT_R1.Pack(13),
		instruction.SUB_REG_LIT.Pack(uint16(register.R1.AsByte()), uint16(3)), // 2 bytes left for literal
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
		instruction.MOV_LIT_R1.Pack(13),
		instruction.MOV_LIT_R2.Pack(12),
		instruction.MUL_REG_REG.Pack(uint16(register.R1.AsByte()), uint16(register.R2.AsByte())),
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
		instruction.MOV_LIT_R1.Pack(13),
		instruction.MUL_REG_LIT.Pack(uint16(register.R1.AsByte()), uint16(3)), // 2 bytes left for literal
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
		instruction.MOV_LIT_R1.Pack(val),
		instruction.MOV_LIT_R2.Pack(val),
		instruction.MUL_REG_REG.Pack(uint16(register.R1.AsByte()), uint16(register.R2.AsByte())),
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
		instruction.MOV_LIT_R1.Pack(120),
		instruction.MOV_LIT_R2.Pack(12),
		instruction.DIV_REG_REG.Pack(uint16(register.R1.AsByte()), uint16(register.R2.AsByte())),
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
		instruction.MOV_LIT_R1.Pack(39),
		instruction.DIV_REG_LIT.Pack(uint16(register.R1.AsByte()), uint16(3)), // 2 bytes left for literal
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
		instruction.MOV_LIT_R1.Pack(128),
		instruction.MOV_LIT_R2.Pack(12),
		instruction.DIV_REG_REG.Pack(uint16(register.R1.AsByte()), uint16(register.R2.AsByte())),
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
		instruction.MOV_LIT_R1.Pack(40),
		instruction.DIV_REG_LIT.Pack(uint16(register.R1.AsByte()), uint16(3)), // 2 bytes left for literal
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
		instruction.MOV_LIT_R1.Pack(13),
		instruction.MOV_LIT_R2.Pack(12),
		instruction.MOD_REG_REG.Pack(uint16(register.R1.AsByte()), uint16(register.R2.AsByte())),
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
		instruction.MOV_LIT_R1.Pack(40),
		instruction.MOD_REG_LIT.Pack(uint16(register.R1.AsByte()), uint16(3)), // 2 bytes left for literal
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
		instruction.MOV_LIT_R1.Pack(1023),
		instruction.MOV_LIT_R2.Pack(289),
		instruction.ADD_REG_REG.Pack(uint16(register.R1.AsByte()), uint16(register.R2.AsByte())),
		instruction.MOV_LIT_R3.Pack(161),
		instruction.MOV_REG_MEM.Pack(uint16(register.R3.AsByte())),
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
		instruction.MOV_LIT_R1.Pack(1023),
		instruction.MOV_LIT_R2.Pack(289),
		instruction.ADD_REG_REG.Pack(uint16(register.R1.AsByte()), uint16(register.R2.AsByte())),
		instruction.MOV_LIT_MEM.Pack(uint16(161)),
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

func Test_MovRegReg_GeneralPurpose(t *testing.T) {
	vm := Machine{cpu: cpu.NewCpu(), memory: memory.NewMemory(255)}
	program := packProgram(
		instruction.MOV_LIT_R1.Pack(161),
		instruction.MOV_REG_REG.Pack(uint16(register.R1.AsByte()), uint16(register.R2.AsByte())),
	)
	vm.LoadProgram(0, program)

	if step, err := run(vm); err != nil || step > 3 {
		t.Fatalf("error running machine or machine stuck: step %d, error: %v", step, err)
	}

	if vm.cpu.GetRegister(register.R1) != 161 {
		vm.Debug()
		t.Fatalf("error setting immediate value to register R1")
	}
	if vm.cpu.GetRegister(register.R2) != 161 {
		vm.Debug()
		t.Fatalf("error copying R1 value to R2: %d", vm.cpu.GetRegister(register.R2))
	}
}

func Test_MovRegReg_Ac2General(t *testing.T) {
	vm := Machine{cpu: cpu.NewCpu(), memory: memory.NewMemory(255)}
	program := packProgram(
		instruction.MOV_LIT_R1.Pack(160),
		instruction.ADD_REG_LIT.Pack(uint16(register.R1.AsByte()), uint16(1)),
		instruction.MOV_REG_REG.Pack(uint16(register.Ac.AsByte()), uint16(register.R2.AsByte())),
	)
	vm.LoadProgram(0, program)

	if vm.cpu.GetRegister(register.Ac) != 0 {
		vm.Debug()
		t.Fatalf("machine initial state error: expected empty Ac, got: %d", vm.cpu.GetRegister(register.Ac))
	}

	if step, err := run(vm); err != nil || step > 4 {
		t.Fatalf("error running machine or machine stuck: step %d, error: %v", step, err)
	}

	if vm.cpu.GetRegister(register.R1) != 160 {
		vm.Debug()
		t.Fatalf("error setting immediate value to register R1: %d", vm.cpu.GetRegister(register.R1))
	}
	if vm.cpu.GetRegister(register.Ac) != 161 {
		vm.Debug()
		t.Fatalf("error in Accumulator: %d", vm.cpu.GetRegister(register.Ac))
	}
	if vm.cpu.GetRegister(register.R2) != 161 {
		vm.Debug()
		t.Fatalf("error copying Accumulator to R2: %d", vm.cpu.GetRegister(register.R2))
	}
}

func Test_Jne(t *testing.T) {
	vm := Machine{cpu: cpu.NewCpu(), memory: memory.NewMemory(255)}
	program := packProgram(
		instruction.MOV_LIT_R2.Pack(13),
		instruction.MOV_LIT_R3.Pack(4), // Multiple of 2 because uint16 addresses
		instruction.ADD_REG_LIT.Pack(uint16(register.R1.AsByte()), uint16(1)),
		instruction.MOV_REG_REG.Pack(uint16(register.Ac.AsByte()), uint16(register.R1.AsByte())),
		instruction.JNE.Pack(uint16(register.R2.AsByte()), uint16(register.R3.AsByte())),
	)
	vm.LoadProgram(0, program)

	if vm.cpu.GetRegister(register.Ac) != 0 {
		vm.Debug()
		t.Fatalf("machine initial state error: expected empty Ac, got: %d", vm.cpu.GetRegister(register.Ac))
	}

	if step, err := run(vm); err != nil || step > (13*3)+3 {
		vm.Debug()
		t.Fatalf("error running machine or machine stuck: step %d, error: %v", step, err)
	}

	if vm.cpu.GetRegister(register.R1) != 13 {
		vm.Debug()
		t.Fatalf("error increasing by immediate value in register R1: %d", vm.cpu.GetRegister(register.R1))
	}
	if vm.cpu.GetRegister(register.R2) != 13 {
		vm.Debug()
		t.Fatalf("error setting immediate value for check in register R2: %d", vm.cpu.GetRegister(register.R2))
	}
	if vm.cpu.GetRegister(register.R3) != 4 {
		vm.Debug()
		t.Fatalf("error setting immediate value for address jump in register R3: %d", vm.cpu.GetRegister(register.R3))
	}
	if vm.cpu.GetRegister(register.Ac) != 13 {
		vm.Debug()
		t.Fatalf("error in Accumulator: %d", vm.cpu.GetRegister(register.Ac))
	}
}

func Test_Jeq(t *testing.T) {
	vm := Machine{cpu: cpu.NewCpu(), memory: memory.NewMemory(255)}
	program := packProgram(
		instruction.MOV_LIT_R1.Pack(11),
		instruction.MOV_LIT_R2.Pack(13),
		instruction.MOV_LIT_R3.Pack(6), // Multiple of 2 because uint16 addresses
		instruction.ADD_REG_LIT.Pack(uint16(register.R1.AsByte()), uint16(1)),
		instruction.MOV_REG_REG.Pack(uint16(register.Ac.AsByte()), uint16(register.R1.AsByte())),
		instruction.JNE.Pack(uint16(register.R2.AsByte()), uint16(register.R3.AsByte())),
	)
	vm.LoadProgram(0, program)

	if vm.cpu.GetRegister(register.Ac) != 0 {
		vm.Debug()
		t.Fatalf("machine initial state error: expected empty Ac, got: %d", vm.cpu.GetRegister(register.Ac))
	}

	if step, err := run(vm); err != nil || step != 10 {
		vm.Debug()
		t.Fatalf("error running machine or machine stuck: step %d, error: %v", step, err)
	}

	if vm.cpu.GetRegister(register.R1) != 13 {
		vm.Debug()
		t.Fatalf("error increasing by immediate value in register R1: %d", vm.cpu.GetRegister(register.R1))
	}
	if vm.cpu.GetRegister(register.R2) != 13 {
		vm.Debug()
		t.Fatalf("error setting immediate value for check in register R2: %d", vm.cpu.GetRegister(register.R2))
	}
	if vm.cpu.GetRegister(register.R3) != 6 {
		vm.Debug()
		t.Fatalf("error setting immediate value for address jump in register R3: %d", vm.cpu.GetRegister(register.R3))
	}
	if vm.cpu.GetRegister(register.Ac) != 13 {
		vm.Debug()
		t.Fatalf("error in Accumulator: %d", vm.cpu.GetRegister(register.Ac))
	}
}

func Test_Jgt(t *testing.T) {
	vm := Machine{cpu: cpu.NewCpu(), memory: memory.NewMemory(255)}
	program := packProgram(
		instruction.MOV_LIT_R1.Pack(25),
		instruction.MOV_LIT_R2.Pack(13),
		instruction.MOV_LIT_R3.Pack(6), // Multiple of 2 because uint16 addresses
		instruction.SUB_REG_LIT.Pack(uint16(register.R1.AsByte()), uint16(1)),
		instruction.MOV_REG_REG.Pack(uint16(register.Ac.AsByte()), uint16(register.R1.AsByte())),
		instruction.JGT.Pack(uint16(register.R2.AsByte()), uint16(register.R3.AsByte())),
	)
	vm.LoadProgram(0, program)

	if vm.cpu.GetRegister(register.Ac) != 0 {
		vm.Debug()
		t.Fatalf("machine initial state error: expected empty Ac, got: %d", vm.cpu.GetRegister(register.Ac))
	}

	if step, err := run(vm); err != nil || step > (13*3)+3 {
		vm.Debug()
		t.Fatalf("error running machine or machine stuck: step %d, error: %v", step, err)
	}

	if vm.cpu.GetRegister(register.R1) != 13 {
		vm.Debug()
		t.Fatalf("error increasing by immediate value in register R1: %d", vm.cpu.GetRegister(register.R1))
	}
	if vm.cpu.GetRegister(register.R2) != 13 {
		vm.Debug()
		t.Fatalf("error setting immediate value for check in register R2: %d", vm.cpu.GetRegister(register.R2))
	}
	if vm.cpu.GetRegister(register.R3) != 6 {
		vm.Debug()
		t.Fatalf("error setting immediate value for address jump in register R3: %d", vm.cpu.GetRegister(register.R3))
	}
	if vm.cpu.GetRegister(register.Ac) != 13 {
		vm.Debug()
		t.Fatalf("error in Accumulator: %d", vm.cpu.GetRegister(register.Ac))
	}
}

func Test_Jge(t *testing.T) {
	vm := Machine{cpu: cpu.NewCpu(), memory: memory.NewMemory(255)}
	program := packProgram(
		instruction.MOV_LIT_R1.Pack(25),
		instruction.MOV_LIT_R2.Pack(13),
		instruction.MOV_LIT_R3.Pack(6), // Multiple of 2 because uint16 addresses
		instruction.SUB_REG_LIT.Pack(uint16(register.R1.AsByte()), uint16(1)),
		instruction.MOV_REG_REG.Pack(uint16(register.Ac.AsByte()), uint16(register.R1.AsByte())),
		instruction.JGE.Pack(uint16(register.R2.AsByte()), uint16(register.R3.AsByte())),
	)
	vm.LoadProgram(0, program)

	if vm.cpu.GetRegister(register.Ac) != 0 {
		vm.Debug()
		t.Fatalf("machine initial state error: expected empty Ac, got: %d", vm.cpu.GetRegister(register.Ac))
	}

	if step, err := run(vm); err != nil || step > (13*3)+5 {
		vm.Debug()
		t.Fatalf("error running machine or machine stuck: step %d, error: %v", step, err)
	}

	if vm.cpu.GetRegister(register.R1) != 12 {
		vm.Debug()
		t.Fatalf("error increasing by immediate value in register R1: %d", vm.cpu.GetRegister(register.R1))
	}
	if vm.cpu.GetRegister(register.R2) != 13 {
		vm.Debug()
		t.Fatalf("error setting immediate value for check in register R2: %d", vm.cpu.GetRegister(register.R2))
	}
	if vm.cpu.GetRegister(register.R3) != 6 {
		vm.Debug()
		t.Fatalf("error setting immediate value for address jump in register R3: %d", vm.cpu.GetRegister(register.R3))
	}
	if vm.cpu.GetRegister(register.Ac) != 12 {
		vm.Debug()
		t.Fatalf("error in Accumulator: %d", vm.cpu.GetRegister(register.Ac))
	}
}

func Test_Jlt(t *testing.T) {
	vm := Machine{cpu: cpu.NewCpu(), memory: memory.NewMemory(255)}
	program := packProgram(
		instruction.MOV_LIT_R1.Pack(13),
		instruction.MOV_LIT_R2.Pack(25),
		instruction.MOV_LIT_R3.Pack(6), // Multiple of 2 because uint16 addresses
		instruction.ADD_REG_LIT.Pack(uint16(register.R1.AsByte()), uint16(1)),
		instruction.MOV_REG_REG.Pack(uint16(register.Ac.AsByte()), uint16(register.R1.AsByte())),
		instruction.JLT.Pack(uint16(register.R2.AsByte()), uint16(register.R3.AsByte())),
	)
	vm.LoadProgram(0, program)

	if vm.cpu.GetRegister(register.Ac) != 0 {
		vm.Debug()
		t.Fatalf("machine initial state error: expected empty Ac, got: %d", vm.cpu.GetRegister(register.Ac))
	}

	if step, err := run(vm); err != nil || step > (13*3)+3 {
		vm.Debug()
		t.Fatalf("error running machine or machine stuck: step %d, error: %v", step, err)
	}

	if vm.cpu.GetRegister(register.R1) != 25 {
		vm.Debug()
		t.Fatalf("error increasing by immediate value in register R1: %d", vm.cpu.GetRegister(register.R1))
	}
	if vm.cpu.GetRegister(register.R2) != 25 {
		vm.Debug()
		t.Fatalf("error setting immediate value for check in register R2: %d", vm.cpu.GetRegister(register.R2))
	}
	if vm.cpu.GetRegister(register.R3) != 6 {
		vm.Debug()
		t.Fatalf("error setting immediate value for address jump in register R3: %d", vm.cpu.GetRegister(register.R3))
	}
	if vm.cpu.GetRegister(register.Ac) != 25 {
		vm.Debug()
		t.Fatalf("error in Accumulator: %d", vm.cpu.GetRegister(register.Ac))
	}
}

func Test_Jle(t *testing.T) {
	vm := Machine{cpu: cpu.NewCpu(), memory: memory.NewMemory(255)}
	program := packProgram(
		instruction.MOV_LIT_R1.Pack(13),
		instruction.MOV_LIT_R2.Pack(26),
		instruction.MOV_LIT_R3.Pack(6), // Multiple of 2 because uint16 addresses
		instruction.ADD_REG_LIT.Pack(uint16(register.R1.AsByte()), uint16(1)),
		instruction.MOV_REG_REG.Pack(uint16(register.Ac.AsByte()), uint16(register.R1.AsByte())),
		instruction.JLE.Pack(uint16(register.R2.AsByte()), uint16(register.R3.AsByte())),
	)
	vm.LoadProgram(0, program)

	if vm.cpu.GetRegister(register.Ac) != 0 {
		vm.Debug()
		t.Fatalf("machine initial state error: expected empty Ac, got: %d", vm.cpu.GetRegister(register.Ac))
	}

	if step, err := run(vm); err != nil || step > (13*3)+7 {
		vm.Debug()
		t.Fatalf("error running machine or machine stuck: step %d, error: %v", step, err)
	}

	if vm.cpu.GetRegister(register.R1) != 27 {
		vm.Debug()
		t.Fatalf("error increasing by immediate value in register R1: %d", vm.cpu.GetRegister(register.R1))
	}
	if vm.cpu.GetRegister(register.R2) != 26 {
		vm.Debug()
		t.Fatalf("error setting immediate value for check in register R2: %d", vm.cpu.GetRegister(register.R2))
	}
	if vm.cpu.GetRegister(register.R3) != 6 {
		vm.Debug()
		t.Fatalf("error setting immediate value for address jump in register R3: %d", vm.cpu.GetRegister(register.R3))
	}
	if vm.cpu.GetRegister(register.Ac) != 27 {
		vm.Debug()
		t.Fatalf("error in Accumulator: %d", vm.cpu.GetRegister(register.Ac))
	}
}
