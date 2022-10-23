package machine

import (
	"fmt"
	"strings"
	"the-machine/machine/cpu"
	"the-machine/machine/instruction"
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
			return fmt.Errorf("error loading program at %d+%d (%#02x): %v", at, idx, b, err)
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
		return instr, fmt.Errorf("unable to get next instruction: %v", err)
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
		return instruction.Descriptors[instruction.NOP], fmt.Errorf("unknown instruction: %#02x", instr)
	}

	// fmt.Printf("cmd: %v (%d)\npass:\n%016b\n%016b\n", decoded.Description, kind, instr, raw)

	decoded.Raw = raw
	return decoded, nil
}

func (vm *Machine) execute(instr instruction.Instruction) error {
	vm.cycle = Execute
	if err := instr.Execute(vm.cpu, vm.ram); err != nil {
		vm.status = Error
		return fmt.Errorf("error executing %#02x: %v", instr, err)
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
		return fmt.Errorf("unable to fetch next tick: %v", err)
	}

	decoded, err := vm.decode(next)
	if err != nil {
		return fmt.Errorf("unable to decode instruction: %#02x", next)
	}

	if err := vm.execute(decoded); err != nil {
		return fmt.Errorf("unable to execute tick: %v", err)
	}

	vm.cycle = Idle

	return nil
}

func (vm Machine) IsDone() bool {
	return vm.status == Done || vm.status == Error
}

func (vm *Machine) Debug() {
	bpos := make([]string, 8)
	bval := make([]string, 8)
	positions := make([]string, 8)
	values := make([]string, 8)
	apos := make([]string, 8)
	aval := make([]string, 8)

	ad := vm.cpu.GetRegister(register.Ip)

	if ad > 8 {
		for i := 0; i < 8; i++ {
			pos := (ad + uint16(i)) - 8
			bpos[i] = fmt.Sprintf("%04d", pos)
			b, _ := vm.ram.GetByte(memory.Address(pos))
			bval[i] = fmt.Sprintf("%#02x", b)
		}
	}

	for i := 0; i < 8; i++ {
		pos := ad + uint16(i)
		positions[i] = fmt.Sprintf("%04d", pos)
		b, _ := vm.ram.GetByte(memory.Address(pos))
		values[i] = fmt.Sprintf("%#02x", b)
	}

	for i := 0; i < 8; i++ {
		pos := ad + uint16(i+8)
		apos[i] = fmt.Sprintf("%04d", pos)
		b, _ := vm.ram.GetByte(memory.Address(pos))
		aval[i] = fmt.Sprintf("%#02x", b)
	}

	positions = append(positions, " | ")
	values = append(values, " | ")

	putReg := func(name string, r register.Register) {
		r1 := vm.cpu.GetRegister(r)
		value := fmt.Sprintf("%4d", r1)
		format := fmt.Sprintf("%%%ds", len(value))

		positions = append(positions, fmt.Sprintf(format, name))
		values = append(values, value)
	}

	putReg("R1", register.R1)
	putReg("R2", register.R2)
	putReg("R3", register.R3)
	putReg("R4", register.R4)
	putReg("Ip", register.Ip)
	putReg("Ac", register.Ac)

	posStr := strings.Join(positions, " ")
	after := strings.Join(apos, " ")

	fmt.Println()
	if ad > 8 {
		before := strings.Join(bpos, " ")
		fmt.Println(before)
		fmt.Println(strings.Repeat("-", len(before)))
		fmt.Println(strings.Join(bval, " "))
		fmt.Println()
	}
	fmt.Println(posStr)
	fmt.Println(strings.Repeat("-", len(posStr)))
	fmt.Println(strings.Join(values, " "))
	fmt.Println()
	fmt.Println(after)
	fmt.Println(strings.Repeat("-", len(after)))
	fmt.Println(strings.Join(aval, " "))
	fmt.Println()
}
