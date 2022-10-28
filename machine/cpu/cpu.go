package cpu

import (
	"fmt"
	"the-machine/machine/memory"
	"the-machine/machine/register"
)

const stackSize = 255

type Cpu struct {
	ip        uint16
	sp        uint16
	fp        uint16
	ac        uint16
	registers map[register.Register]uint16
	stack     *memory.Memory
	stackSize int
}

func NewCpu() *Cpu {
	registers := make(map[register.Register]uint16, 8)
	registers[register.R1] = 0
	registers[register.R2] = 0
	registers[register.R3] = 0
	registers[register.R4] = 0
	registers[register.R5] = 0
	registers[register.R6] = 0
	registers[register.R7] = 0
	registers[register.R8] = 0
	mem := memory.NewMemory(stackSize).(*memory.Memory)
	return &Cpu{registers: registers, stack: mem}
}

func (cpu Cpu) GetRegister(r register.Register) uint16 {
	switch r {
	case register.Ip:
		return cpu.ip
	case register.Sp:
		return cpu.sp
	case register.Fp:
		return cpu.fp
	case register.Ac:
		return cpu.ac
	default:
		if reg, ok := cpu.registers[r]; ok {
			return reg
		}
		return 0
	}
}

func (cpu *Cpu) SetRegister(r register.Register, v uint16) {
	switch r {
	case register.Ip:
		cpu.ip = v
	case register.Sp:
		cpu.sp = v
	case register.Fp:
		cpu.fp = v
	case register.Ac:
		cpu.ac = v
	default:
		cpu.registers[r] = v
	}
}

func (cpu *Cpu) Push(value uint16) error {
	address := cpu.GetRegister(register.Sp)
	address += 2
	if address >= stackSize {
		return fmt.Errorf("stack overflow, unable to push %d (%#02x) to %d (%#02x)", value, value, address, address)
	}

	if err := cpu.stack.SetUint16(memory.Address(address), value); err != nil {
		return fmt.Errorf("stack overflow: %w", err)
	}
	cpu.stackSize++
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
		return value, fmt.Errorf("stack underflow: %w", err)
	}
	cpu.stackSize--
	cpu.SetRegister(register.Sp, address-2)
	return value, nil
}

func (cpu *Cpu) StoreFrame() error {
	stackHead := uint16(cpu.stackSize)

	if err := cpu.Push(cpu.GetRegister(register.R1)); err != nil {
		return fmt.Errorf("error storing register R1: %w", err)
	}
	if err := cpu.Push(cpu.GetRegister(register.R2)); err != nil {
		return fmt.Errorf("error storing register R2: %w", err)
	}
	if err := cpu.Push(cpu.GetRegister(register.R3)); err != nil {
		return fmt.Errorf("error storing register R3: %w", err)
	}
	if err := cpu.Push(cpu.GetRegister(register.R4)); err != nil {
		return fmt.Errorf("error storing register R4: %w", err)
	}
	if err := cpu.Push(cpu.GetRegister(register.Ip)); err != nil {
		return fmt.Errorf("error storing register Ip: %w", err)
	}
	if err := cpu.Push(stackHead); err != nil {
		return fmt.Errorf("error storing stack head: %w", err)
	}

	cpu.stackSize = 0
	cpu.SetRegister(register.Fp, cpu.GetRegister(register.Sp))
	return nil
}

func (cpu *Cpu) RestoreFrame() error {
	framePointer := cpu.GetRegister(register.Fp)
	cpu.SetRegister(register.Sp, framePointer)

	stackHead, err := cpu.Pop()
	if err != nil {
		return fmt.Errorf("error restoring frame, no stack head: %w", err)
	}

	if value, err := cpu.Pop(); err != nil {
		return fmt.Errorf("error restoring instruction pointer: %w", err)
	} else {
		cpu.SetRegister(register.Ip, value)
	}

	if value, err := cpu.Pop(); err != nil {
		return fmt.Errorf("error restoring register 4: %w", err)
	} else {
		cpu.SetRegister(register.R4, value)
	}

	if value, err := cpu.Pop(); err != nil {
		return fmt.Errorf("error restoring register 3: %w", err)
	} else {
		cpu.SetRegister(register.R3, value)
	}

	if value, err := cpu.Pop(); err != nil {
		return fmt.Errorf("error restoring register 2: %w", err)
	} else {
		cpu.SetRegister(register.R2, value)
	}

	if value, err := cpu.Pop(); err != nil {
		return fmt.Errorf("error restoring register 1: %w", err)
	} else {
		cpu.SetRegister(register.R1, value)
	}

	// reset stack
	cpu.stackSize = int(stackHead)
	for i := 0; i < int(stackHead); i++ {
		cpu.Pop()
	}

	cpu.SetRegister(register.Fp, cpu.GetRegister(register.Sp))

	return nil
}

// TODO: Fugly, used just in debugging
func (x Cpu) GetStack() (int, memory.MemoryAccess) {
	return x.stackSize, x.stack
}
