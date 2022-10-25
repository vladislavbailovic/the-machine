package machine

import (
	"fmt"
	"io"
	"the-machine/machine/debug"
	"the-machine/machine/memory"
	"the-machine/machine/register"
)

type MemoryType uint8

const (
	RAM MemoryType = 0
	ROM MemoryType = iota
)

type Debugger struct {
	vm       *Machine
	stream   io.Writer
	renderer *debug.Renderer
}

func NewDebugger(vm *Machine, f debug.Formatter) *Debugger {
	return &Debugger{vm: vm, renderer: debug.NewRenderer(f)}
}

func (x *Debugger) SetFormatter(f debug.Formatter) {
	x.renderer.SetFormatter(f)
}

func (x Debugger) CoreRegisters() string {
	return x.renderer.Registers(x.vm.cpu, []register.Register{
		register.Ip,
		register.Ac,
		register.Sp,
		register.Fp,
	})
}

func (x Debugger) GeneralRegisters() string {
	return x.renderer.Registers(x.vm.cpu, []register.Register{
		register.R1,
		register.R2,
		register.R3,
		register.R4,
		register.R5,
		register.R6,
		register.R7,
		register.R8,
	})
}

func (x Debugger) AllRegisters() string {
	return x.renderer.Registers(x.vm.cpu, []register.Register{
		register.Ip,
		register.Ac,
		register.Sp,
		register.Fp,
		register.R1,
		register.R2,
		register.R3,
		register.R4,
		register.R5,
		register.R6,
		register.R7,
		register.R8,
	})
}

func (x Debugger) Peek(startAt memory.Address, outputLen int, srcType MemoryType) string {
	var source memory.MemoryAccess
	switch srcType {
	case RAM:
		source = x.vm.ram
	case ROM:
		source = x.vm.rom
	default:
		x.out(fmt.Sprintf("ERROR: unknown source type: %v", srcType))
		return ""
	}

	return x.renderer.Memory(source, startAt, outputLen)
}

func (x Debugger) Disassemble(startAt memory.Address, outputLen int) string {
	source := x.vm.rom
	return x.renderer.Disassembly(source, x.vm.decode, startAt, outputLen)
}

func (x Debugger) Dump() error {
	dumper := debug.NewDumper()
	return dumper.Dump(x.vm.rom)
}

// TODO: figure this out
func (x Debugger) out(msg string) {
	fmt.Println(msg)
}
