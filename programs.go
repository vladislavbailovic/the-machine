package main

import (
	"encoding/binary"
	"fmt"
	"net/http"
	"os"
	"strings"
	"the-machine/cmd"
	"the-machine/machine"
	"the-machine/machine/debug"
	"the-machine/machine/device"
	"the-machine/machine/instruction"
	"the-machine/machine/memory"
	"the-machine/machine/register"
)

type responseStatusWriter struct {
	resp http.ResponseWriter
}

func (x responseStatusWriter) Write(b []byte) (int, error) {
	status := int(binary.LittleEndian.Uint16(b))
	x.resp.WriteHeader(status)
	return 0, nil
}

type responseBodyWriter struct {
	resp http.ResponseWriter
}

func (x responseBodyWriter) Write(b []byte) (int, error) {
	c := byte(binary.LittleEndian.Uint16(b))
	x.resp.Write([]byte{c})
	return 0, nil
}

func main_Microservice() {
	vm := machine.NewMachine(2048)

	method := device.FileDescriptor(12)
	path := device.FileDescriptor(13)

	responseStatus := device.FileDescriptor(16)
	responseBody := device.FileDescriptor(61)

	interceptor := func(w http.ResponseWriter, r *http.Request) {
		io, err := vm.GetIO()
		if err != nil {
			panic(err)
		}

		mlike := device.NewFilelike(method, device.Read, strings.NewReader(r.Method))
		io.SetDescriptor(method, mlike)

		plike := device.NewFilelike(path, device.Read, strings.NewReader(r.URL.String()))
		io.SetDescriptor(path, plike)

		slike := device.NewFilelike(responseStatus, device.Write, responseStatusWriter{resp: w})
		io.SetDescriptor(responseStatus, slike)

		rlike := device.NewFilelike(responseBody, device.Write, responseBodyWriter{resp: w})
		io.SetDescriptor(responseBody, rlike)

		cmd.Run(vm)
		vm.Reset()
	}

	program := packProgram(
		instruction.MOV_LIT_R1.Pack(150),
		instruction.CALL.Pack(register.R1.AsUint16()),

		instruction.MOV_LIT_BNK.Pack(uint16(memory.DeviceIO)),

		instruction.MOV_LIT_AC.Pack(uint16(responseStatus)),
		instruction.MOV_LIT_R1.Pack(201),
		instruction.MOV_REG_MEM.Pack(register.R1.AsUint16()),

		instruction.MOV_LIT_AC.Pack(uint16(responseBody)),
		instruction.MOV_LIT_R1.Pack(uint16('O')),
		instruction.MOV_REG_MEM.Pack(register.R1.AsUint16()),
		instruction.MOV_LIT_R1.Pack(uint16('K')),
		instruction.MOV_REG_MEM.Pack(register.R1.AsUint16()),
	)
	verifyMethod := packProgram(
		instruction.PUSH_LIT.Pack(uint16('T')),
		instruction.PUSH_LIT.Pack(uint16('E')),
		instruction.PUSH_LIT.Pack(uint16('G')),

		instruction.MOV_LIT_R3.Pack(158),
		instruction.MOV_LIT_R5.Pack(210),
		instruction.MOV_LIT_BNK.Pack(uint16(memory.DeviceIO)),

		instruction.MOV_LIT_AC.Pack(uint16(method)),
		instruction.MOV_MEM_REG.Pack(register.Ac.AsUint16(), register.R1.AsUint16()),

		instruction.POP_REG.Pack(register.R2.AsUint16()),

		instruction.MOV_REG_REG.Pack(register.R2.AsUint16(), register.Ac.AsUint16()),
		instruction.JEQ.Pack(register.R4.AsUint16(), register.R5.AsUint16()),

		instruction.MOV_REG_REG.Pack(register.R1.AsUint16(), register.Ac.AsUint16()),
		instruction.JEQ.Pack(register.R2.AsUint16(), register.R3.AsUint16()),

		instruction.MOV_LIT_AC.Pack(uint16(responseStatus)),
		instruction.MOV_LIT_R1.Pack(503),
		instruction.MOV_REG_MEM.Pack(register.R1.AsUint16()),
		instruction.MOV_LIT_AC.Pack(uint16(responseBody)),
		instruction.MOV_LIT_R1.Pack(uint16('n')),
		instruction.MOV_REG_MEM.Pack(register.R1.AsUint16()),
		instruction.MOV_LIT_R1.Pack(uint16('o')),
		instruction.MOV_REG_MEM.Pack(register.R1.AsUint16()),
		instruction.MOV_LIT_R1.Pack(uint16('t')),
		instruction.MOV_REG_MEM.Pack(register.R1.AsUint16()),
		instruction.MOV_LIT_R1.Pack(uint16(' ')),
		instruction.MOV_REG_MEM.Pack(register.R1.AsUint16()),
		instruction.MOV_LIT_R1.Pack(uint16('o')),
		instruction.MOV_REG_MEM.Pack(register.R1.AsUint16()),
		instruction.MOV_LIT_R1.Pack(uint16('k')),
		instruction.MOV_REG_MEM.Pack(register.R1.AsUint16()),
		instruction.HALT.Pack(0),

		instruction.RET.Pack(0),
	)
	vm.LoadProgram(0, program)
	vm.LoadProgram(150, verifyMethod)

	http.HandleFunc("/", interceptor)
	http.ListenAndServe(":6660", nil)
}

func main_RemapStdio_CopyToStdout() {
	vm := machine.NewMachine(2048)
	io, err := vm.GetIO()
	if err != nil {
		panic(err)
	}

	fd := device.FileDescriptor(13)
	filelike := device.NewFilelike(fd, device.Read, strings.NewReader("hai hello"))

	io.SetDescriptor(fd, filelike)

	program := packProgram(
		instruction.MOV_LIT_BNK.Pack(uint16(memory.DeviceIO)),

		instruction.MOV_LIT_AC.Pack(uint16(fd)),
		instruction.MOV_MEM_REG.Pack(register.Ac.AsUint16(), register.R1.AsUint16()),

		instruction.MOV_LIT_AC.Pack(uint16(device.Stdout)),
		instruction.MOV_REG_MEM.Pack(register.R1.AsUint16()),

		instruction.MOV_REG_REG.Pack(register.R1.AsUint16(), register.Ac.AsUint16()),
		instruction.JNE.Pack(register.R2.AsUint16(), register.R3.AsUint16()),
	)
	vm.LoadProgram(0, program)
	cmd.Run(vm)
}

func main_IoStdout_Machine() {
	vm := machine.NewMachine(1024)

	buffer := packProgram(
		instruction.MOV_LIT_BNK.Pack(uint16(memory.DeviceIO)),
		instruction.MOV_LIT_AC.Pack(uint16(device.Stdout)),

		instruction.MOV_LIT_R1.Pack(uint16('o')),
		instruction.MOV_REG_MEM.Pack(register.R1.AsUint16()),

		instruction.MOV_LIT_R1.Pack(uint16('h')),
		instruction.MOV_REG_MEM.Pack(register.R1.AsUint16()),

		instruction.MOV_LIT_R1.Pack(uint16('a')),
		instruction.MOV_REG_MEM.Pack(register.R1.AsUint16()),

		instruction.MOV_LIT_R1.Pack(uint16('i')),
		instruction.MOV_REG_MEM.Pack(register.R1.AsUint16()),

		instruction.MOV_LIT_R1.Pack(uint16(' ')),
		instruction.MOV_REG_MEM.Pack(register.R1.AsUint16()),

		instruction.MOV_LIT_R1.Pack(uint16('t')),
		instruction.MOV_REG_MEM.Pack(register.R1.AsUint16()),

		instruction.MOV_LIT_R1.Pack(uint16('h')),
		instruction.MOV_REG_MEM.Pack(register.R1.AsUint16()),

		instruction.MOV_LIT_R1.Pack(uint16('a')),
		instruction.MOV_REG_MEM.Pack(register.R1.AsUint16()),

		instruction.MOV_LIT_R1.Pack(uint16('r')),
		instruction.MOV_REG_MEM.Pack(register.R1.AsUint16()),

		instruction.MOV_LIT_R1.Pack(uint16('\n')),
		instruction.MOV_REG_MEM.Pack(register.R1.AsUint16()),
	)

	vm.LoadProgram(0, buffer)
	if _, err := cmd.Run(vm); err != nil {
		vm.DebugError(err)
	}
	// vm.Debug()
}

func main_IoStdin() {
	mem := device.NewIoMap()
	if b, err := mem.GetByte(memory.Address(device.Stdin)); err != nil {
		fmt.Printf("ERROR: %v", err)
	} else {
		fmt.Printf("SUCCESS: got %v", b)
	}
}

func main_IoStdout() {
	mem := device.NewIoMap()
	mem.SetByte(memory.Address(device.Stdout), byte('H'))
	mem.SetByte(memory.Address(device.Stdout), byte('e'))
	mem.SetByte(memory.Address(device.Stdout), byte('l'))
	mem.SetByte(memory.Address(device.Stdout), byte('l'))
	mem.SetByte(memory.Address(device.Stdout), byte('0'))
}

func main_InteractiveDebugger_WithProgram() {
	vm := machine.NewMachine(0xffff)
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
		instruction.MOV_REG_REG.Pack(register.Ac.AsUint16(), register.R2.AsUint16()), // R2 = 17850 (limit)
		instruction.PUSH_LIT.Pack(9*2),
		instruction.POP_REG.Pack(register.R3.AsUint16()),                             // R3 = 18 (jump address-1) *2
		instruction.MOV_REG_REG.Pack(register.R8.AsUint16(), register.Ac.AsUint16()), // Ac = 0
		instruction.MOV_REG_MEM.Pack(register.R1.AsUint16()),                         // Fake-Draw
		instruction.ADD_REG_LIT.Pack(register.Ac.AsUint16(), 1),                      // Ac++
		instruction.JLT.Pack(register.R2.AsUint16(), register.R3.AsUint16()),         // If Ac < R2, jump to R3
	))

	// cmd.Run(vm)
	vm.Debug()
}

func loadFromBuffer_Vga() {
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
	cmd.Run(vm)
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
