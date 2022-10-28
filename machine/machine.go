package machine

import (
	"fmt"
	"the-machine/machine/cpu"
	"the-machine/machine/debug"
	"the-machine/machine/instruction"
	"the-machine/machine/internal"
	"the-machine/machine/memory"
	"the-machine/machine/register"
)

type Status uint8

const (
	Ready   Status = 0
	Loaded  Status = iota
	Running Status = iota
	Done    Status = iota
	Error   Status = iota
)

type Cycle uint8

const (
	Idle    Cycle = 0
	Fetch   Cycle = iota
	Decode  Cycle = iota
	Execute Cycle = iota
)

type Machine struct {
	cpu    *cpu.Cpu
	rom    memory.MemoryAccess
	ram    memory.MemoryAccess
	status Status
	cycle  Cycle
}

func NewMachine(memsize int) Machine {
	return Machine{
		cpu:    cpu.NewCpu(),
		rom:    memory.NewMemory(memsize),
		ram:    memory.NewMemory(memsize),
		status: Ready,
		cycle:  Idle,
	}
}

func (vm *Machine) Reset() {
	vm.cpu.Reset()
	vm.status = Ready
	vm.cycle = Idle
}

func NewWithMemory(mem memory.MemoryAccess, ramSize int) Machine {
	return Machine{
		cpu:    cpu.NewCpu(),
		rom:    memory.NewMemory(ramSize),
		ram:    mem,
		status: Ready,
		cycle:  Idle,
	}
}

func (vm *Machine) LoadProgram(at memory.Address, program []byte) error {
	for idx, b := range program {
		if err := vm.rom.SetByte(at+memory.Address(idx), b); err != nil {
			return internal.Error(fmt.Sprintf("error loading program at %d+%d (%#02x)", at, idx, b), err, internal.ErrorLoading)
		}
	}
	vm.status = Loaded
	return nil
}

func (vm *Machine) fetch() (uint16, error) {
	vm.cycle = Fetch
	ip := vm.cpu.GetRegister(register.Ip)

	ipAddr := memory.Address(ip)
	instr, err := vm.rom.GetUint16(ipAddr)
	if err != nil {
		vm.status = Error
		return instr, internal.Error("unable to get next instruction", err, internal.ErrorRuntime)
	}

	vm.cpu.SetRegister(register.Ip, ip+2)

	return instr, nil
}

func (vm *Machine) decode(instr uint16) (instruction.Instruction, error) {
	vm.cycle = Decode

	kind, raw := instruction.Decode(instr)
	if kind == instruction.HALT {
		vm.status = Done
		return instruction.Descriptors[instruction.NOP], nil
	}

	decoded, ok := instruction.Descriptors[kind]
	if !ok {
		vm.status = Error
		return instruction.Descriptors[instruction.NOP], internal.Error(fmt.Sprintf("unknown instruction: %#02x", instr), nil, internal.ErrorRuntime)
	}

	// fmt.Printf("cmd: %v (%d)\npass:\n%016b\n%016b\n", decoded.Description, kind, instr, raw)

	decoded.Raw = raw
	return decoded, nil
}

func (vm *Machine) execute(instr instruction.Instruction) error {
	vm.cycle = Execute
	if err := instr.Execute(vm.cpu, vm.ram); err != nil {
		vm.status = Error
		return internal.Error(fmt.Sprintf("error executing %#02x", instr), err, internal.ErrorRuntime)
	}
	return nil
}

func (vm *Machine) Tick() error {
	if vm.IsDone() {
		return nil
	}
	vm.status = Running
	vm.cycle = Idle

	next, err := vm.fetch()
	if err != nil {
		return internal.Error("unable to fetch next tick", err, internal.ErrorRuntime)
	}

	decoded, err := vm.decode(next)
	if err != nil {
		return internal.Error(fmt.Sprintf("unable to decode instruction: %#02x", next), err, internal.ErrorRuntime)
	}

	if err := vm.execute(decoded); err != nil {
		return internal.Error("unable to execute tick", err, internal.ErrorRuntime)
	}

	vm.cycle = Idle

	return nil
}

func (vm Machine) IsDone() bool {
	return vm.status == Done || vm.status == Error
}

func (vm *Machine) DebugError(err error) {
	fmtr := debug.Formatter{
		Numbers:   debug.Binary,
		OutputAs:  debug.Byte,
		Rendering: debug.Vertical,
	}
	dbg := NewDebugger(vm, fmtr)

	if err != nil {
		dbg.OutError(err)
	}

	dbg.Current()
	dbg.Run()
}
func (vm *Machine) Debug() {
	fmtr := debug.Formatter{
		Numbers:   debug.Binary,
		OutputAs:  debug.Byte,
		Rendering: debug.Vertical,
	}
	dbg := NewDebugger(vm, fmtr)

	dbg.Current()
	dbg.Run()
}
