package cpu

import (
	"fmt"
	"the-machine/machine/memory"
	"the-machine/machine/register"
)

type Cpu struct {
	registers *memory.Memory
	memory    *memory.Memory
}

func NewCpu() *Cpu {
	registers := memory.NewMemory(register.Size())
	return &Cpu{registers: registers, memory: memory.NewMemory(255)}
}

func (cpu *Cpu) LoadProgram(at memory.Address, program []byte) error {
	for idx, b := range program {
		if err := cpu.memory.SetByte(at+memory.Address(idx), b); err != nil {
			return fmt.Errorf("error loading program at %d+%d (0x%02x): %v", at, idx, b, err)
		}
	}
	return nil
}

func (cpu Cpu) GetRegister(r register.Register) (uint16, error) {
	if reg, err := cpu.registers.GetUint16(r.AsAddress()); err != nil {
		return 0, fmt.Errorf("Unknown register: %d: %v", reg, err)
	} else {
		return reg, nil
	}
}

func (cpu *Cpu) SetRegister(r register.Register, v uint16) error {
	return cpu.registers.SetUint16(r.AsAddress(), v)
}
