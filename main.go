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
		byte(instruction.MOV_LIT_AC), 0x02, 0x00,
		byte(instruction.MUL_REG_LIT), register.Ac.AsByte(), 0x02, 0x00,
		byte(instruction.JLT), 0x07, 0x00, 0x03, 0x00,
		byte(instruction.MOD_REG_LIT), register.Ac.AsByte(), 0x03, 0x00,
		byte(instruction.HALT),
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
