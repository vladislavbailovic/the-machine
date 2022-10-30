package main

import (
	"os"
	"the-machine/cmd"
	"the-machine/machine"
	"the-machine/machine/instruction"
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

func packSubroutine(instr ...[]byte) []byte {
	return packStatements(instruction.RET, instr...)
}

func main() {
	if len(os.Args) > 1 {
		fname := os.Args[1]
		// TODO: validate fname
		cmd.RunFile(fname)
	} else {
		main_InteractiveDebugger()
	}
}

func main_InteractiveDebugger() {
	vm := machine.NewMachine(0xffff)
	vm.Debug()
}
