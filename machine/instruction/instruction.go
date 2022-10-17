package instruction

import (
	"fmt"
	"the-machine/machine/cpu"
	"the-machine/machine/memory"
	"the-machine/machine/register"
)

type Instruction struct {
	Description string
	Parameters  []Parameter
	Executor    Executor
}

func (x Instruction) Execute(cpu *cpu.Cpu, memory *memory.Memory) error {
	params, err := x.getParams(cpu, memory)
	if err != nil {
		return fmt.Errorf("unable to execute \"%s\": %v", x.Description, err)
	}
	err = x.Executor.Execute(params, cpu, memory)
	if err != nil {
		return fmt.Errorf("error executing \"%s\": %v", x.Description, err)
	}
	return nil
}

func (x Instruction) getParams(cpu *cpu.Cpu, mem *memory.Memory) ([]byte, error) {
	length := 0
	for _, p := range x.Parameters {
		length += p.GetBytesLength()
	}
	params := []byte{}

	pos := cpu.GetRegister(register.Ip)
	for idx, p := range x.Parameters {
		switch p {
		case ParamRegister:
			val, err := mem.GetByte(memory.Address(pos))
			if err != nil {
				return params, fmt.Errorf("%s: error getting param %d: %v", x.Description, idx, err)
			}
			pos++
			params = append(params, val)
		case ParamLiteral, ParamAddress:
			hi, err := mem.GetByte(memory.Address(pos))
			if err != nil {
				return params, fmt.Errorf("%s: error getting param %d: %v", x.Description, idx, err)
			}
			params = append(params, hi)
			pos++
			lo, err := mem.GetByte(memory.Address(pos))
			if err != nil {
				return params, fmt.Errorf("%s: error getting param %d: %v", x.Description, idx, err)
			}
			pos++
			params = append(params, lo)
		default:
			return params, fmt.Errorf("unexpected parameter: %d", p)
		}
	}
	cpu.SetRegister(register.Ip, pos)

	return params, nil
}
