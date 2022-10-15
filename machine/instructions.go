package machine

import (
	"fmt"
	"the-machine/machine/instruction"
	"the-machine/machine/register"
)

type instr struct {
	description string
	parameters  []instruction.Parameter
	executor    func(*cpu, []byte) instruction.Executor
}

func (x instr) Execute(cpu *cpu) error {
	params, err := x.getParams(cpu)
	if err != nil {
		return fmt.Errorf("unable to execute \"%s\": %v", x.description, err)
	}
	result, err := x.executor(cpu, params).Execute(params)
	if err != nil {
		return fmt.Errorf("error executing \"%s\": %v", x.description, err)
	}
	x.processResult(cpu, result)
	return nil
}

func (x instr) getParams(cpu *cpu) ([]byte, error) {
	length := 0
	for _, p := range x.parameters {
		length += int(p)
	}
	params := []byte{}

	pos, err := cpu.getRegister(register.Ip)
	if err != nil {
		return params, fmt.Errorf("%s: error fetching Ip register: %v", x.description, err)
	}
	for idx, p := range x.parameters {
		switch p {
		case instruction.PARAM8:
			val, err := cpu.memory.GetByte(address(pos))
			if err != nil {
				return params, fmt.Errorf("%s: error getting param %d: %v", x.description, idx, err)
			}
			pos++
			params = append(params, val)
		case instruction.PARAM16:
			hi, err := cpu.memory.GetByte(address(pos))
			if err != nil {
				return params, fmt.Errorf("%s: error getting param %d: %v", x.description, idx, err)
			}
			params = append(params, hi)
			pos++
			lo, err := cpu.memory.GetByte(address(pos))
			if err != nil {
				return params, fmt.Errorf("%s: error getting param %d: %v", x.description, idx, err)
			}
			pos++
			params = append(params, lo)
		default:
			return params, fmt.Errorf("unexpected parameter: %d", p)
		}
	}
	if err != cpu.setRegister(register.Ip, pos) {
		return params, fmt.Errorf("%s: error updating Ip register: %v", x.description, err)
	}

	return params, nil
}

func (x instr) processResult(cpu *cpu, res instruction.Result) {
	switch res.Action {
	case instruction.Nop:
		return
	case instruction.RecordRegister:
		cpu.setRegister(res.Target, res.Value)
	}
}

var Instructions = map[instruction.Instruction]instr{
	instruction.NOP: {
		description: "No-op",
		executor: func(cpu *cpu, _ []byte) instruction.Executor {
			return instruction.Passthrough{}
		},
	},
	instruction.MOV_LIT_AC: {
		description: "Move literal to register AC",
		parameters: []instruction.Parameter{
			instruction.PARAM16,
		},
		executor: func(cpu *cpu, _ []byte) instruction.Executor {
			return instruction.Lit2Reg{Target: register.Ac}
		},
	},
	instruction.MOV_LIT_R1: {
		description: "Move literal to register R1",
		parameters: []instruction.Parameter{
			instruction.PARAM16,
		},
		executor: func(cpu *cpu, _ []byte) instruction.Executor {
			return instruction.Lit2Reg{Target: register.R1}
		},
	},
	instruction.MOV_LIT_R2: {
		description: "Move literal to register R2",
		parameters: []instruction.Parameter{
			instruction.PARAM16,
		},
		executor: func(cpu *cpu, _ []byte) instruction.Executor {
			return instruction.Lit2Reg{Target: register.R2}
		},
	},
	instruction.MOV_LIT_R3: {
		description: "Move literal to register R3",
		parameters: []instruction.Parameter{
			instruction.PARAM16,
		},
		executor: func(cpu *cpu, _ []byte) instruction.Executor {
			return instruction.Lit2Reg{Target: register.R3}
		},
	},
	instruction.MOV_LIT_R4: {
		description: "Move literal to register R4",
		parameters: []instruction.Parameter{
			instruction.PARAM16,
		},
		executor: func(cpu *cpu, _ []byte) instruction.Executor {
			return instruction.Lit2Reg{Target: register.R2}
		},
	},
	instruction.ADD_REG_REG: {
		description: "Add contents of two registers",
		parameters: []instruction.Parameter{
			instruction.PARAM8,
			instruction.PARAM8,
		},
		executor: func(cpu *cpu, params []byte) instruction.Executor {
			if len(params) != 2 {
				return instruction.ExecError{
					Error: fmt.Errorf("ADD_REG_REG: invalid params: %v", params)}
			}
			v1, err := cpu.getRegister(register.Register(params[0]))
			if err != nil {
				return instruction.ExecError{
					Error: fmt.Errorf("ADD_REG_REG: error fetching from register %d (#1): %v", params[0], err)}
			}
			v2, err := cpu.getRegister(register.Register(params[1]))
			if err != nil {
				return instruction.ExecError{
					Error: fmt.Errorf("ADD_REG_REG: error fetching from register %d (#2): %v", params[1], err)}
			}
			return instruction.AddTwo{V1: v1, V2: v2}
		},
	},
	instruction.JNE: {
		description: "Jump if not equal",
		parameters: []instruction.Parameter{
			instruction.PARAM16,
			instruction.PARAM16,
		},
		executor: func(cpu *cpu, _ []byte) instruction.Executor {
			value, err := cpu.getRegister(register.Register(register.Ac))
			if err != nil {
				return instruction.ExecError{
					Error: fmt.Errorf("JNE: error fetching from Ac: %v", err)}
			}
			return instruction.Jump{Against: value}
		},
	},
	instruction.HALT: {
		description: "Halt execution",
		executor: func(cpu *cpu, _ []byte) instruction.Executor {
			end := cap(*cpu.memory)
			return instruction.Halt{End: uint16(end)}
		},
	},
}
