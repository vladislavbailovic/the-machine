package main

// https://github.dev/lowbyteproductions/16-Bit-Virtual-Machine/tree/master/episode-1

import (
	"fmt"
	"os"
	"the-machine/machine"
	"the-machine/machine/debug"
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

func packSubroutine(instr ...[]byte) []byte {
	return packStatements(instruction.RET, instr...)
}

func run(vm machine.Machine) (int, error) {
	step := 0
	for step < 0xffff {
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

func main3() {
	vga := device.NewVideo()
	vm := machine.NewWithMemory(vga, 1024)
	setLimit := packSubroutine(
		instruction.PUSH_LIT.Pack(1023),
		instruction.PUSH_LIT.Pack(17),
		instruction.MUL_STACK.Pack(),
		instruction.PUSH_LIT.Pack(459),
		instruction.ADD_STACK.Pack(),
		instruction.POP_REG.Pack(register.Ac.AsUint16()),
	)
	vm.LoadProgram(500, setLimit)
	vm.LoadProgram(0, packProgram(
		instruction.PUSH_LIT.Pack(65),
		instruction.POP_REG.Pack(register.R1.AsUint16()), // R1 = 65 (draw char)
		instruction.PUSH_LIT.Pack(500),
		instruction.POP_REG.Pack(register.R4.AsUint16()), // R4 = 500 (subroutine address)
		instruction.CALL.Pack(register.R4.AsUint16()),
		instruction.MOV_REG_REG.Pack(register.Ac.AsUint16(), register.R2.AsUint16()), // R2 = 17920 (limit)
		instruction.PUSH_LIT.Pack(9*2),
		instruction.POP_REG.Pack(register.R3.AsUint16()),                             // R3 = 18 (jump address-1) *2
		instruction.MOV_REG_REG.Pack(register.R8.AsUint16(), register.Ac.AsUint16()), // Ac = 0
		instruction.MOV_REG_MEM.Pack(register.R1.AsUint16()),                         // Draw
		instruction.ADD_REG_LIT.Pack(register.Ac.AsUint16(), 1),                      // Ac++
		instruction.JLT.Pack(register.R2.AsUint16(), register.R3.AsUint16()),         // If Ac < R2, jump to R3
	))
	run(vm)

	fmtr := debug.Formatter{
		Numbers:   debug.Binary,
		OutputAs:  debug.Byte,
		Rendering: debug.Vertical,
	}
	dbg := machine.NewDebugger(&vm, fmtr)
	fmt.Println()
	fmt.Println(dbg.Peek(0, 8, machine.RAM))
	fmt.Println(dbg.Disassemble(0, 12))

	fmtr.Numbers = debug.Decimal
	fmtr.Rendering = debug.Horizontal
	dbg.SetFormatter(fmtr)
	fmt.Println(dbg.AllRegisters())

	dbg.Dump()
}

func main() {
	buffer, err := os.ReadFile("out.bin")
	if err != nil {
		panic(err)
	}
	vga := device.NewVideo()
	vm := machine.NewWithMemory(vga, 1024)
	vm.LoadProgram(0, buffer)

	fmtr := debug.Formatter{
		Numbers:   debug.Binary,
		OutputAs:  debug.Byte,
		Rendering: debug.Vertical,
	}
	dbg := machine.NewDebugger(&vm, fmtr)
	dbg.Run()
	// fmt.Println()
	// fmt.Println(dbg.Peek(0, 8, machine.RAM))
	// fmt.Println(dbg.Disassemble(0, 4))

	// fmtr.Numbers = debug.Decimal
	// fmtr.Rendering = debug.Horizontal
	// dbg.SetFormatter(fmtr)
	// fmt.Println(dbg.AllRegisters())

	// fmt.Println("^ that was loaded o.0")
}

func outAll() {
	vga := device.NewVideo()
	vm := machine.NewWithMemory(vga, 1024)
	vm.LoadProgram(0, packProgram(
		instruction.MOV_LIT_R1.Pack(4),                                               // 01: R1 = 4
		instruction.SHL_REG_LIT.Pack(register.R1.AsUint16(), 4),                      // 02: Ac = 64
		instruction.ADD_REG_LIT.Pack(register.Ac.AsUint16(), 1),                      // 03: Ac = 65
		instruction.MOV_REG_REG.Pack(register.Ac.AsUint16(), register.R1.AsUint16()), // 04: R1 = 65 (draw char)
		instruction.MOV_LIT_R2.Pack(8),                                               // 05: R2 = 15
		instruction.SHL_REG_LIT.Pack(register.R2.AsUint16(), 8),                      // 06: Ac = 2048
		instruction.SHL_REG_LIT.Pack(register.Ac.AsUint16(), 4),                      // 07: Ac = 32768
		instruction.SUB_REG_LIT.Pack(register.Ac.AsUint16(), 1),                      // 08: Ac = 32767
		instruction.MUL_REG_LIT.Pack(register.Ac.AsUint16(), 2),                      // 09: Ac = 65534
		instruction.MOV_REG_REG.Pack(register.Ac.AsUint16(), register.R2.AsUint16()), // 10: R2 = 65534 (limit)
		instruction.MOV_LIT_R3.Pack(15),                                              // 11: R3 = 15
		instruction.ADD_REG_LIT.Pack(register.R3.AsUint16(), 13),                     // 12: Ac = 28
		instruction.MOV_REG_REG.Pack(register.Ac.AsUint16(), register.R3.AsUint16()), // 13: R3 = 28 (jump address-1)*2
		instruction.MOV_REG_REG.Pack(register.R8.AsUint16(), register.Ac.AsUint16()), // 14: Ac = 0
		instruction.MOV_REG_MEM.Pack(register.R1.AsUint16()),                         // 15: Draw
		instruction.ADD_REG_LIT.Pack(register.Ac.AsUint16(), 1),                      // 16: Ac++
		instruction.JLT.Pack(register.R2.AsUint16(), register.R3.AsUint16()),         // 17: If Ac < R2, jump to R3
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
