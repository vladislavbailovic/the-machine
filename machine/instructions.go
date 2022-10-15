package machine

import (
	"fmt"
	"the-machine/machine/cpu"
	"the-machine/machine/instruction"
	"the-machine/machine/register"
)

var Instructions = map[instruction.Type]instruction.Instruction{
	instruction.NOP: {
		Description: "No-op",
		Executor: func(cpu *cpu.Cpu, _ []byte) instruction.Executor {
			return instruction.Passthrough{}
		},
	},
	instruction.MOV_LIT_AC: {
		Description: "Move literal to register AC",
		Parameters: []instruction.Parameter{
			instruction.PARAM16,
		},
		Executor: func(cpu *cpu.Cpu, _ []byte) instruction.Executor {
			return instruction.Lit2Reg{Target: register.Ac}
		},
	},
	instruction.MOV_LIT_R1: {
		Description: "Move literal to register R1",
		Parameters: []instruction.Parameter{
			instruction.PARAM16,
		},
		Executor: func(cpu *cpu.Cpu, _ []byte) instruction.Executor {
			return instruction.Lit2Reg{Target: register.R1}
		},
	},
	instruction.MOV_LIT_R2: {
		Description: "Move literal to register R2",
		Parameters: []instruction.Parameter{
			instruction.PARAM16,
		},
		Executor: func(cpu *cpu.Cpu, _ []byte) instruction.Executor {
			return instruction.Lit2Reg{Target: register.R2}
		},
	},
	instruction.MOV_LIT_R3: {
		Description: "Move literal to register R3",
		Parameters: []instruction.Parameter{
			instruction.PARAM16,
		},
		Executor: func(cpu *cpu.Cpu, _ []byte) instruction.Executor {
			return instruction.Lit2Reg{Target: register.R3}
		},
	},
	instruction.MOV_LIT_R4: {
		Description: "Move literal to register R4",
		Parameters: []instruction.Parameter{
			instruction.PARAM16,
		},
		Executor: func(cpu *cpu.Cpu, _ []byte) instruction.Executor {
			return instruction.Lit2Reg{Target: register.R2}
		},
	},
	instruction.ADD_REG_REG: {
		Description: "Add contents of two registers",
		Parameters: []instruction.Parameter{
			instruction.PARAM8,
			instruction.PARAM8,
		},
		Executor: func(c *cpu.Cpu, params []byte) instruction.Executor {
			if len(params) != 2 {
				return instruction.ExecError{
					Error: fmt.Errorf("ADD_REG_REG: invalid params: %v", params)}
			}
			v1, err := c.GetRegister(register.Register(params[0]))
			if err != nil {
				return instruction.ExecError{
					Error: fmt.Errorf("ADD_REG_REG: error fetching from register %d (#1): %v", params[0], err)}
			}
			v2, err := c.GetRegister(register.Register(params[1]))
			if err != nil {
				return instruction.ExecError{
					Error: fmt.Errorf("ADD_REG_REG: error fetching from register %d (#2): %v", params[1], err)}
			}
			return instruction.AddTwo{V1: v1, V2: v2}
		},
	},
	instruction.JNE: {
		Description: "Jump if not equal",
		Parameters: []instruction.Parameter{
			instruction.PARAM16,
			instruction.PARAM16,
		},
		Executor: func(cpu *cpu.Cpu, _ []byte) instruction.Executor {
			value, err := cpu.GetRegister(register.Register(register.Ac))
			if err != nil {
				return instruction.ExecError{
					Error: fmt.Errorf("JNE: error fetching from Ac: %v", err)}
			}
			return instruction.Jump{Against: value}
		},
	},
}
