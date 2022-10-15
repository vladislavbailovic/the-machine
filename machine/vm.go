package machine

import (
	"fmt"
	"the-machine/machine/cpu"
	"the-machine/machine/instruction"
	"the-machine/machine/memory"
	"the-machine/machine/register"
)

type Vm struct {
	cpu    *cpu.Cpu
	memory *memory.Memory
}

func (vm *Vm) LoadProgram(at memory.Address, program []byte) error {
	for idx, b := range program {
		if err := vm.memory.SetByte(at+memory.Address(idx), b); err != nil {
			return fmt.Errorf("error loading program at %d+%d (0x%02x): %v", at, idx, b, err)
		}
	}
	return nil
}

func (vm *Vm) nextInstruction() (byte, error) {
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

func (vm *Vm) executeInstruction(instr byte) error {
	instructionType := instruction.Type(instr)
	if instructionType == instruction.END || instructionType == instruction.HALT {
		return nil
	}

	instruction, ok := Instructions[instructionType]
	if !ok {
		return fmt.Errorf("unknown instruction: 0x%02x", instr)
	}
	if err := instruction.Execute(vm.cpu, vm.memory); err != nil {
		return fmt.Errorf("error executing 0x%02x: %v", instr, err)
	}
	return nil
}

func (vm *Vm) tick() error {
	next, err := vm.nextInstruction()
	if err != nil {
		return fmt.Errorf("unable to fetch next tick: %v", err)
	}

	if next == byte(instruction.END) {
		// We are done here
		return nil
	}

	if err := vm.executeInstruction(next); err != nil {
		return fmt.Errorf("unable to execute tick: %v", err)
	}

	return nil
}

func (vm *Vm) debug() {
	ad, _ := vm.cpu.GetRegister(register.Ip)
	fmt.Printf("%02d:   ", ad)
	for i := 0; i < 8; i++ {
		b, _ := vm.memory.GetByte(memory.Address(ad + uint16(i)))
		fmt.Printf("0x%02x ", b)
	}
	r1, _ := vm.cpu.GetRegister(register.R1)
	r2, _ := vm.cpu.GetRegister(register.R2)
	ac, _ := vm.cpu.GetRegister(register.Ac)
	fmt.Printf("  R1: %02d, R2: %02d, Ac: %02d", r1, r2, ac)
	fmt.Println()
}
