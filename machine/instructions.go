package machine

import (
	"fmt"
	"the-machine/machine/instruction"
	"the-machine/machine/register"
)

type instr struct {
	description string
	execute     func(*cpu) error
}

var Instructions = map[instruction.Instruction]instr{
	instruction.NOP: {
		description: "No-op",
		execute: func(cpu *cpu) error {
			return nil
		},
	},
	instruction.MOV_LIT_R1: {
		description: "Move literal to register R1",
		execute: func(cpu *cpu) error {
			fmt.Println("MOV_LIT_R1")
			cpu.debug()
			pos, err := cpu.getRegister(register.Ip)
			if err != nil {
				return fmt.Errorf("MOV_LIT_R1: error fetching Ip register: %v", err)
			}
			val, err := cpu.memory.GetUint16(address(pos))
			fmt.Printf("\tgot value: %d from %d\n", val, pos)
			if err != nil {
				return fmt.Errorf("MOV_LIT_R1: error getting value to store: %v", err)
			}
			pos += 2
			cpu.setRegister(register.Ip, pos)
			if err != cpu.setRegister(register.R1, val) {
				return fmt.Errorf("MOV_LIT_R1: error storing value %d into register: %v", val, err)
			}
			fmt.Printf("\tvalue %d is now in register %d\n", val, register.R1)
			return nil
		},
	},
	instruction.MOV_LIT_R2: {
		description: "Move literal to register R2",
		execute: func(cpu *cpu) error {
			fmt.Println("MOV_LIT_R2")
			cpu.debug()
			pos, err := cpu.getRegister(register.Ip)
			if err != nil {
				return fmt.Errorf("MOV_LIT_R2: error fetching Ip register: %v", err)
			}
			val, err := cpu.memory.GetUint16(address(pos))
			fmt.Printf("\tgot value: %d from %d\n", val, pos)
			if err != nil {
				return fmt.Errorf("MOV_LIT_R2: error getting value to store: %v", err)
			}
			pos += 2
			cpu.setRegister(register.Ip, pos)
			if err != cpu.setRegister(register.R2, val) {
				return fmt.Errorf("MOV_LIT_R2: error storing value %d into register: %v", val, err)
			}
			fmt.Printf("\tvalue %d is now in register %d\n", val, register.R2)
			return nil
		},
	},
	instruction.ADD_REG_REG: {
		description: "Add contents of two registers",
		execute: func(cpu *cpu) error {
			fmt.Println("ADD_REG_REG")
			pos, err := cpu.getRegister(register.Ip)
			if err != nil {
				return fmt.Errorf("ADD_REG_REG: error fetching Ip register: %v", err)
			}
			r1, err := cpu.memory.GetByte(address(pos))
			if err != nil {
				return fmt.Errorf("ADD_REG_REG: error fetching #1 register from memory: %v", err)
			}
			v1, err := cpu.getRegister(register.Register(r1))
			if err != nil {
				return fmt.Errorf("ADD_REG_REG: error fetching from register %d (#1): %v", r1, err)
			}
			pos++
			r2, err := cpu.memory.GetByte(address(pos))
			if err != nil {
				return fmt.Errorf("ADD_REG_REG: error fetching #2 register from memory: %v", err)
			}
			v2, err := cpu.getRegister(register.Register(r2))
			if err != nil {
				return fmt.Errorf("ADD_REG_REG: error fetching from register %d (#2): %v", r1, err)
			}
			pos++
			cpu.setRegister(register.Ip, pos)

			res := v1 + v2
			if err := cpu.setRegister(register.Ac, res); err != nil {
				return fmt.Errorf("ADD_REG_REG: error writing result %d to Ac: %v", res, err)
			}

			return nil
		},
	},
}
