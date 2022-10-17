package cpu

import (
	"the-machine/machine/memory"
	"the-machine/machine/register"
)

type Cpu struct {
	registers map[register.Register]uint16
	memory    *memory.Memory
}

func NewCpu() *Cpu {
	registers := make(map[register.Register]uint16, register.Size())
	return &Cpu{registers: registers, memory: memory.NewMemory(255)}
}

func (cpu Cpu) GetRegister(r register.Register) uint16 {
	if reg, ok := cpu.registers[r]; ok {
		return reg
	}
	return 0
}

func (cpu *Cpu) SetRegister(r register.Register, v uint16) {
	cpu.registers[r] = v
}
