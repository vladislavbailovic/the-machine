package main

// https://github.dev/lowbyteproductions/16-Bit-Virtual-Machine/tree/master/episode-1

import (
	"fmt"
	"the-machine/machine"
	"the-machine/machine/instruction"
	"the-machine/machine/register"
)

func main() {
	vm := machine.NewMachine(255)
	vm.LoadProgram(0, []byte{
		byte(instruction.MOV_LIT_AC), 0x01, 0x00,
		byte(instruction.MOV_LIT_R1), 0x01, 0x00,
		byte(instruction.ADD_REG_REG), byte(register.Ac), byte(register.R1),
		byte(instruction.JNE), 0x03, 0x00, 0x06, 0x00,
		byte(instruction.MOV_LIT_R2), 0xac, 0xab,
		byte(instruction.HALT),
	})

	var response string
	for true {
		vm.Tick()
		vm.Debug()
		if vm.IsDone() {
			break
		}
		fmt.Scanln(&response)
	}
	fmt.Println("bye!")
}
