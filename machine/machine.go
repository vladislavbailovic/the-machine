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

type Machine struct {
	cpu    *cpu.Cpu
	memory *memory.Memory
	status Status
}

func NewMachine(memsize int) Machine {
	return Machine{
		cpu:    cpu.NewCpu(),
		memory: memory.NewMemory(memsize),
		status: Ready,
	}
}

func (vm *Machine) LoadProgram(at memory.Address, program []byte) error {
	for idx, b := range program {
		if err := vm.memory.SetByte(at+memory.Address(idx), b); err != nil {
			return fmt.Errorf("error loading program at %d+%d (%#02x): %v", at, idx, b, err)
		}
	}
	vm.status = Loaded
	return nil
}

func (vm *Machine) nextInstruction() (byte, error) {
	ip, err := vm.cpu.GetRegister(register.Ip)
	if err != nil {
		return 0, fmt.Errorf("unable to access IP register: %v", err)
	}

	ipAddr := memory.Address(ip)
	instr, err := vm.memory.GetByte(ipAddr)
	if err != nil {
		return instr, fmt.Errorf("unable to get next instruction: %v", err)
	}

	if err := vm.cpu.SetRegister(register.Ip, ip+1); err != nil {
		return instr, fmt.Errorf("unable to update IP register: %v", err)
	}

	return instr, nil
}

func (vm *Machine) executeInstruction(instr byte) error {
	instructionType := instruction.Type(instr)
	if instructionType == instruction.END || instructionType == instruction.HALT {
		vm.status = Done
		return nil
	}

	vm.status = Running
	instruction, ok := Instructions[instructionType]
	if !ok {
		vm.status = Error
		return fmt.Errorf("unknown instruction: %#02x", instr)
	}
	if err := instruction.Execute(vm.cpu, vm.memory); err != nil {
		vm.status = Error
		return fmt.Errorf("error executing %#02x: %v", instr, err)
	}
	return nil
}

func (vm *Machine) Tick() error {
	if vm.IsDone() {
		return nil
	}

	next, err := vm.nextInstruction()
	if err != nil {
		return fmt.Errorf("unable to fetch next tick: %v", err)
	}

	if next == instruction.END.AsByte() {
		// We are done here
		return nil
	}

	if err := vm.executeInstruction(next); err != nil {
		return fmt.Errorf("unable to execute tick: %v", err)
	}

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

	ad, _ := vm.cpu.GetRegister(register.Ip)

	if ad > 8 {
		for i := 0; i < 8; i++ {
			pos := (ad + uint16(i)) - 8
			bpos[i] = fmt.Sprintf("%04d", pos)
			b, _ := vm.memory.GetByte(memory.Address(pos))
			bval[i] = fmt.Sprintf("%#02x", b)
		}
	}

	for i := 0; i < 8; i++ {
		pos := ad + uint16(i)
		positions[i] = fmt.Sprintf("%04d", pos)
		b, _ := vm.memory.GetByte(memory.Address(pos))
		values[i] = fmt.Sprintf("%#02x", b)
	}

	for i := 0; i < 8; i++ {
		pos := ad + uint16(i+8)
		apos[i] = fmt.Sprintf("%04d", pos)
		b, _ := vm.memory.GetByte(memory.Address(pos))
		aval[i] = fmt.Sprintf("%#02x", b)
	}

	positions = append(positions, " | ")
	values = append(values, " | ")

	putReg := func(name string, r register.Register) {
		r1, _ := vm.cpu.GetRegister(r)
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
