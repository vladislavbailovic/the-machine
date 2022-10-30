package main

import (
	"os"
	"strings"
	"the-machine/cmd"
	"the-machine/machine"
	"the-machine/machine/device"
	"the-machine/machine/instruction"
	"the-machine/machine/memory"
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

func packSubroutine(instr ...[]byte) []byte {
	return packStatements(instruction.RET, instr...)
}

func main() {
	vm := machine.NewMachine(2048)
	// interceptor := func(w http.ResponseWriter, r *http.Request) {
	// 	vm.
	// }
	io, err := vm.GetIO()
	if err != nil {
		panic(err)
	}

	fd := device.FileDescriptor(13)
	filelike := device.NewFilelike(fd, device.Read, strings.NewReader("hai hello"))

	io.SetDescriptor(fd, filelike)

	program := packProgram(
		instruction.MOV_LIT_BNK.Pack(uint16(memory.DeviceIO)),

		instruction.MOV_LIT_AC.Pack(13),
		instruction.MOV_MEM_REG.Pack(register.Ac.AsUint16(), register.R1.AsUint16()),

		instruction.MOV_LIT_AC.Pack(uint16(device.Stdout)),
		instruction.MOV_REG_MEM.Pack(register.R1.AsUint16()),

		instruction.MOV_REG_REG.Pack(register.R1.AsUint16(), register.Ac.AsUint16()),
		instruction.JNE.Pack(register.R2.AsUint16(), register.R3.AsUint16()),
	)
	vm.LoadProgram(0, program)
	cmd.Run(vm)
	// vm.Debug()
}

func main_Runner() {
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
