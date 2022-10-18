package cpu

import (
	"fmt"
	"the-machine/machine/memory"
	"the-machine/machine/register"
)

const stackSize = 255

type Cpu struct {
	registers map[register.Register]uint16
	stack     *memory.Memory
}

func NewCpu() *Cpu {
	registers := make(map[register.Register]uint16, 6)
	registers[register.Ip] = 0
	registers[register.Ac] = 0
	registers[register.Sp] = 0
	registers[register.R1] = 0
	registers[register.R2] = 0
	registers[register.R3] = 0
	registers[register.R4] = 0
	return &Cpu{registers: registers, stack: memory.NewMemory(stackSize)}
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

func (cpu *Cpu) Push(value uint16) error {
	address := cpu.GetRegister(register.Sp)
	address += 2
	if address >= stackSize {
		return fmt.Errorf("stack overflow, unable to push %d (%#02x) to %d (%#02x)", value, value, address, address)
	}

	if err := cpu.stack.SetUint16(memory.Address(address), value); err != nil {
		return fmt.Errorf("stack overflow: %v", err)
	}
	cpu.SetRegister(register.Sp, address)
	return nil
}

func (cpu *Cpu) Pop() (uint16, error) {
	address := cpu.GetRegister(register.Sp)
	if address < 2 {
		return 0, fmt.Errorf("stack underflow, unable to pop from %d (%#02x))", address, address)
	}

	value, err := cpu.stack.GetUint16(memory.Address(address))
	if err != nil {
		return value, fmt.Errorf("stack underflow: %v", err)
	}
	cpu.SetRegister(register.Sp, address-2)
	return value, nil
}
