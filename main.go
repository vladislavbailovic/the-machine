package main

// https://github.dev/lowbyteproductions/16-Bit-Virtual-Machine/tree/master/episode-1

import (
	"fmt"
	"the-machine/machine"
	"the-machine/machine/device"
	"the-machine/machine/instruction"
	"the-machine/machine/register"
)

func packStatements(last instruction.Type, instr ...[]byte) []byte {
	res := make([]byte, 0, len(instr)+2)
	for _, b := range instr {
		res = append(res, b...)
	}
	res = append(res, last.Pack(0)...)
	return res
}

func packProgram(instr ...[]byte) []byte {
	return packStatements(instruction.HALT, instr...)
}

func run(vm machine.Machine) (int, error) {
	step := 0
	for step < 65534 {
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

func main() {
	vga := device.NewVideo()
	vm := machine.NewWithMemory(vga)
	vm.LoadProgram(0, packProgram(
		instruction.MOV_LIT_R1.Pack(4),                                               // R1 = 4
		instruction.SHL_REG_LIT.Pack(register.R1.AsUint16(), 4),                      // Ac = 64
		instruction.ADD_REG_LIT.Pack(register.Ac.AsUint16(), 1),                      // Ac = 65
		instruction.MOV_REG_REG.Pack(register.Ac.AsUint16(), register.R1.AsUint16()), // R1 = 65
		instruction.MOV_LIT_R2.Pack(8),                                               // R2 = 15
		instruction.SHL_REG_LIT.Pack(register.R2.AsUint16(), 8),                      // Ac = 2048
		instruction.SHL_REG_LIT.Pack(register.Ac.AsUint16(), 4),                      // Ac = 32768
		instruction.SUB_REG_LIT.Pack(register.Ac.AsUint16(), 1),                      // Ac = 32767
		instruction.MUL_REG_LIT.Pack(register.Ac.AsUint16(), 2),                      // Ac = 65534
		instruction.MOV_REG_REG.Pack(register.Ac.AsUint16(), register.R2.AsUint16()), // R2 = 65534
		instruction.MOV_LIT_R3.Pack(15),                                              // R3 = 15 (jump address)
		instruction.ADD_REG_LIT.Pack(register.R3.AsUint16(), 13),                     // Ac = 28 (jump address)
		instruction.MOV_REG_REG.Pack(register.Ac.AsUint16(), register.R3.AsUint16()), // R3 = 28 (jump address)
		instruction.MOV_REG_REG.Pack(register.R8.AsUint16(), register.Ac.AsUint16()), // Ac = 0
		instruction.MOV_REG_MEM.Pack(register.R1.AsUint16()),
		instruction.ADD_REG_LIT.Pack(register.Ac.AsUint16(), 1), // Ac++
		instruction.JLT.Pack(register.R2.AsUint16(), register.R3.AsUint16()),
	))
	run(vm)
	vm.Debug()
}

func main2() {
	vm := machine.NewMachine(255)
	vm.LoadProgram(0, []byte{
		instruction.MUL_REG_LIT.AsByte(), register.Ac.AsByte(), 0x02, 0x00,
		instruction.JLT.AsByte(), 0x07, 0x00, 0x03, 0x00,
		instruction.MOD_REG_LIT.AsByte(), register.Ac.AsByte(), 0x03, 0x00,
		instruction.HALT.AsByte(),
	})

	var response string
	for true {
		err := vm.Tick()
		vm.Debug()
		if err != nil {
			fmt.Println(err)
		}
		if vm.IsDone() {
			break
		}
		fmt.Scanln(&response)
	}
	fmt.Println("bye!")
}
