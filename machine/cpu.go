package machine

import (
	"fmt"
	"the-machine/machine/instruction"
	"the-machine/machine/register"
)

type cpu struct {
	registers *memory
	memory    *memory
}

func NewCpu() cpu {
	registers := NewMemory(int(register.Size()) * 2)
	return cpu{registers: registers, memory: NewMemory(255)}
}

func (cpu *cpu) LoadProgram(at address, program []byte) error {
	for idx, b := range program {
		if err := cpu.memory.SetByte(at+address(idx), b); err != nil {
			return fmt.Errorf("error loading program at %d+%d (0x%02x): %v", at, idx, b, err)
		}
	}
	return nil
}

func (cpu cpu) getRegister(r register.Register) (uint16, error) {
	if reg, err := cpu.registers.GetUint16(address(r.Address())); err != nil {
		return 0, fmt.Errorf("Unknown register: %d: %v", reg, err)
	} else {
		return reg, nil
	}
}

func (cpu *cpu) setRegister(r register.Register, v uint16) error {
	return cpu.registers.SetUint16(address(r.Address()), v)
}

func (cpu *cpu) nextInstruction() (byte, error) {
	ip, err := cpu.getRegister(register.Ip)
	if err != nil {
		return 0, fmt.Errorf("unable to access IP register: %v", err)
	}

	ipAddr := address(ip)
	instr, err := cpu.memory.GetByte(ipAddr)
	if err != nil {
		return instr, fmt.Errorf("unable to get next instruction: %v", err)
	}

	if err := cpu.setRegister(register.Ip, ip+1); err != nil {
		return instr, fmt.Errorf("unable to update IP register: %v", err)
	}

	return instr, nil
}

func (cpu *cpu) executeInstruction(instr byte) error {
	instruction, ok := Instructions[instruction.Instruction(instr)]
	if !ok {
		return fmt.Errorf("unknown instruction: 0x%02x", instr)
	}
	if instruction.execute == nil {
		return fmt.Errorf("missing executor for 0x%02x", instr)
	}
	if err := instruction.execute(cpu); err != nil {
		return fmt.Errorf("error executing 0x%02x: %v", instr, err)
	}
	return nil
}

func (cpu *cpu) tick() error {
	next, err := cpu.nextInstruction()
	if err != nil {
		return fmt.Errorf("unable to fetch next tick: %v", err)
	}

	if err := cpu.executeInstruction(next); err != nil {
		return fmt.Errorf("unable to execute tick: %v", err)
	}

	return nil
}

func (cpu cpu) debug() {
	ad, _ := cpu.getRegister(register.Ip)
	fmt.Printf("%02d:   ", ad)
	for i := 0; i < 8; i++ {
		b, _ := cpu.memory.GetByte(address(ad + uint16(i)))
		fmt.Printf("0x%02x ", b)
	}
	r1, _ := cpu.getRegister(register.R1)
	r2, _ := cpu.getRegister(register.R2)
	ac, _ := cpu.getRegister(register.Ac)
	fmt.Printf("  R1: %02d, R2: %02d, Ac: %02d", r1, r2, ac)
	fmt.Println()
}