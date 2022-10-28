package machine

import (
	"fmt"
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
	renderer *debug.Renderer
	skin     *debug.Interface
}

func NewDebugger(vm *Machine, f debug.Formatter) *Debugger {
	return &Debugger{vm: vm, renderer: debug.NewRenderer(f), skin: debug.NewInterface()}
}

func (x Debugger) current() {
	x.renderer.Out("")
	x.currentRam()
	x.currentDisassembly()
	x.currentRegisters()
	x.renderer.Out("")
}

func (x Debugger) currentRam() {
	x.renderer.Out("[ Memory ]")
	x.renderer.Out(x.Peek(0, 8, RAM))
}

func (x Debugger) currentRom() {
	memPos := x.vm.cpu.GetRegister(register.Ip)
	x.renderer.Out("[ Program ]")
	x.renderer.Out(x.Peek(memory.Address(memPos), 8, ROM))
}

func (x Debugger) currentStack() {
	x.renderer.Out("[ Stack ]")
	_, stack := x.vm.cpu.GetStack()
	x.renderer.Out(x.renderer.Stack(x.vm.cpu.GetRegister(register.Sp), stack))
}

func (x Debugger) currentDisassembly() {
	memPos := x.vm.cpu.GetRegister(register.Ip)
	x.renderer.Out("[ Disassembly ]")
	x.renderer.Out(x.Disassemble(memory.Address(memPos), 4))
}

func (x Debugger) currentRegisters() {
	x.renderer.Out("[ Registers ]")
	old := x.renderer.GetFormatter()
	f := debug.Formatter{
		Numbers:   debug.Decimal,
		OutputAs:  debug.Uint,
		Rendering: debug.Horizontal,
	}
	x.renderer.SetFormatter(f)
	x.renderer.Out(x.AllRegisters())
	x.renderer.SetFormatter(old)
}

func (x Debugger) Run() {
	doTick := false
	ticks := 0
	for true {
		if doTick {
			err := x.vm.Tick()
			if err != nil {
				x.renderer.Out(fmt.Sprintf("ERROR: runtime error: %v", err))
			} else {
				ticks++
			}
			if x.vm.IsDone() {
				break
			}
		}
		doTick = true
		x.skin.Prompt(ticks)
		cmd, err := x.skin.GetCommand()
		if err != nil {
			x.renderer.Out(fmt.Sprintf("ERROR: debugger error: %v", err))
			continue
		}
		switch cmd.Action {
		case debug.Tick:
			x.renderer.Out("")
			continue
		case debug.Next:
			x.current()
			continue
		case debug.Inspect:
			x.current()
			doTick = false
			continue
		case debug.PeekRam:
			x.currentRam()
			doTick = false
			continue
		case debug.PeekRom:
			x.currentRom()
			doTick = false
			continue
		case debug.Stack:
			x.currentStack()
			doTick = false
			continue
		case debug.Disassemble:
			x.currentDisassembly()
			doTick = false
			continue
		case debug.Registers:
			x.currentRegisters()
			doTick = false
			continue
		case debug.Dump:
			x.Dump()
			continue
		case debug.Quit:
			break
		}
	}
	x.renderer.Out("bye!")
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
		x.renderer.Out(fmt.Sprintf("ERROR: unknown source type: %v", srcType))
		return ""
	}

	return x.renderer.Memory(source, startAt, outputLen)
}

func (x Debugger) Disassemble(startAt memory.Address, outputLen int) string {
	source := x.vm.rom
	return x.renderer.Disassembly(source, startAt, outputLen)
}

func (x Debugger) Dump() error {
	dumper := debug.NewDumper()
	return dumper.Dump(x.vm.rom)
}
