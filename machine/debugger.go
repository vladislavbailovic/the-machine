package machine

import (
	"fmt"
	"the-machine/machine/debug"
	"the-machine/machine/internal"
	"the-machine/machine/memory"
	"the-machine/machine/register"
)

type Debugger struct {
	vm       *Machine
	renderer *debug.Renderer
	skin     *debug.Interface
}

func NewDebugger(vm *Machine, f debug.Formatter) *Debugger {
	return &Debugger{vm: vm, renderer: debug.NewRenderer(f), skin: debug.NewInterface()}
}

func (x Debugger) Current() {
	x.renderer.Out("")
	x.currentRam()
	x.currentDisassembly()
	x.currentRegisters()
	x.renderer.Out("")
}

func (x Debugger) currentRam() {
	x.ramAt(0, 8)
}

func (x Debugger) ramAt(memPos memory.Address, length int) {
	x.renderer.Out("[ Memory ]")
	x.renderer.Out(x.Peek(memPos, length, memory.RAM))
}

func (x Debugger) currentRom() {
	memPos := x.vm.cpu.GetRegister(register.Ip)
	x.romAt(memory.Address(memPos), 8)
}
func (x Debugger) romAt(memPos memory.Address, length int) {
	x.renderer.Out("[ Program ]")
	x.renderer.Out(x.Peek(memPos, length, memory.ROM))
}

func (x Debugger) currentStack() {
	x.renderer.Out("[ Stack ]")
	_, stack := x.vm.cpu.GetStack()
	x.renderer.Out(x.renderer.Stack(x.vm.cpu.GetRegister(register.Sp), stack))
}

func (x Debugger) currentDisassembly() {
	memPos := x.vm.cpu.GetRegister(register.Ip)
	x.disassemblyAt(memory.Address(memPos), 4)
}
func (x Debugger) disassemblyAt(memPos memory.Address, length int) {
	x.renderer.Out("[ Disassembly ]")
	x.renderer.Out(x.Disassemble(memory.Address(memPos), length))
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
				x.renderer.OutError("runtime error", err)
			} else {
				ticks++
				if x.vm.IsDone() {
					break
				}
			}
		}
		doTick = true
		x.skin.Prompt(ticks, x.vm.cpu.GetRegister(register.Ip))
		cmd, err := x.skin.GetCommand()
		if err != nil {
			x.renderer.OutError("debugger error", err)
			continue
		}
		switch cmd.GetAction() {
		case debug.Tick:
			x.renderer.Out("")
			continue
		case debug.Next:
			x.Current()
			continue
		case debug.Inspect:
			x.Current()
			doTick = false
			continue
		case debug.PeekRam:
			if peek, ok := cmd.(debug.PeekCommand); ok {
				x.ramAt(peek.At, peek.Length)
			} else {
				x.currentRam()
			}
			doTick = false
			continue
		case debug.PeekRom:
			if peek, ok := cmd.(debug.PeekCommand); ok {
				x.romAt(peek.At, peek.Length)
			} else {
				x.currentRom()
			}
			doTick = false
			continue
		case debug.Stack:
			x.currentStack()
			doTick = false
			continue
		case debug.Disassemble:
			if peek, ok := cmd.(debug.PeekCommand); ok {
				x.disassemblyAt(peek.At, peek.Length)
			} else {
				x.currentDisassembly()
			}
			doTick = false
			continue
		case debug.Registers:
			x.currentRegisters()
			doTick = false
			continue
		case debug.Dump:
			if err := x.Dump(); err != nil {
				x.renderer.OutError("debugger error", err)
			} else {
				x.renderer.Out("Successfully dumped memory to file")
			}
			doTick = false
			continue
		case debug.Load:
			if err := x.Load(); err != nil {
				x.renderer.OutError("debugger error", err)
			} else {
				x.renderer.Out("Successfully loaded buffer")
			}
			doTick = false
			continue
		case debug.Reset:
			x.vm.Reset()
			doTick = false
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

func (x Debugger) getMemory(srcType memory.MemoryType) memory.MemoryAccess {
	var source memory.MemoryAccess
	switch srcType {
	case memory.RAM:
		if ram, err := x.vm.getMemory(memory.RAM); err != nil {
			x.renderer.OutError("debugger error", internal.Error("unable to access RAM", err, internal.ErrorDebugger))
		} else {
			source = ram
		}
	case memory.ROM:
		if rom, err := x.vm.getMemory(memory.ROM); err != nil {
			x.renderer.OutError("debugger error", internal.Error("unable to access ROM", err, internal.ErrorDebugger))
		} else {
			source = rom
		}
	default:
		// x.renderer.OutError(fmt.Sprintf("ERROR: unknown source type: %v", srcType))
		x.renderer.OutError("debugger error", internal.Error(fmt.Sprintf("unknown source type: %v", srcType), nil, internal.ErrorDebugger))
	}
	return source
}

func (x Debugger) Peek(startAt memory.Address, outputLen int, srcType memory.MemoryType) string {
	source := x.getMemory(srcType)
	return x.renderer.Memory(source, startAt, outputLen)
}

func (x Debugger) Disassemble(startAt memory.Address, outputLen int) string {
	source := x.getMemory(memory.ROM)
	return x.renderer.Disassembly(source, startAt, outputLen)
}

func (x Debugger) Dump() error {
	dumper := debug.NewAsciiDumper(debug.Decimal)
	source := x.getMemory(memory.ROM)
	return dumper.Dump(source)
}

func (x Debugger) Load() error {
	dumper := debug.NewAsciiDumper(debug.Decimal)
	if rom, err := dumper.Load(); err != nil {
		return err
	} else {
		x.vm.LoadProgram(0, rom)
		x.vm.Reset()
	}
	return nil
}

func (x Debugger) OutError(err error) {
	x.renderer.OutError("error", err)
}
