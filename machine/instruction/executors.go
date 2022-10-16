package instruction

import (
	"encoding/binary"
	"fmt"
	"the-machine/machine/cpu"
	"the-machine/machine/register"
)

type Executor interface {
	Execute([]byte, *cpu.Cpu) error
}

type Passthrough struct{}

func (x Passthrough) Execute(p []byte, cpu *cpu.Cpu) error {
	return nil
}

type Lit2Reg struct {
	Target register.Register
}

func (x Lit2Reg) Execute(params []byte, cpu *cpu.Cpu) error {
	if len(params) != 2 {
		return fmt.Errorf("LIT2REG[%v]: invalid parameter: %v", x.Target, params)
	}
	value := binary.LittleEndian.Uint16(params)
	return cpu.SetRegister(x.Target, value)
}

type AddTwo struct{}

func (x AddTwo) Execute(params []byte, cpu *cpu.Cpu) error {
	if len(params) != 2 {
		return fmt.Errorf("ADD_REG_REG: invalid params: %v", params)
	}
	v1, err := cpu.GetRegister(register.Register(params[0]))
	if err != nil {
		return fmt.Errorf("ADD_REG_REG: error fetching from register %d (#1): %v", params[0], err)
	}
	v2, err := cpu.GetRegister(register.Register(params[1]))
	if err != nil {
		return fmt.Errorf("ADD_REG_REG: error fetching from register %d (#2): %v", params[1], err)
	}
	return cpu.SetRegister(register.Ac, v1+v2)
}

type SubTwo struct{}

func (x SubTwo) Execute(params []byte, cpu *cpu.Cpu) error {
	if len(params) != 2 {
		return fmt.Errorf("SUB_REG_REG: invalid params: %v", params)
	}
	v1, err := cpu.GetRegister(register.Register(params[0]))
	if err != nil {
		return fmt.Errorf("SUB_REG_REG: error fetching from register %d (#1): %v", params[0], err)
	}
	v2, err := cpu.GetRegister(register.Register(params[1]))
	if err != nil {
		return fmt.Errorf("SUB_REG_REG: error fetching from register %d (#2): %v", params[1], err)
	}
	return cpu.SetRegister(register.Ac, v1-v2)
}

type Jump struct{}

func (x Jump) Execute(params []byte, cpu *cpu.Cpu) error {
	against, err := cpu.GetRegister(register.Register(register.Ac))
	if err != nil {
		return fmt.Errorf("JNE: error fetching from Ac: %v", err)
	}

	if len(params) != 4 {
		return fmt.Errorf("JNE[%v]: invalid parameter: %v", against, params)
	}
	value := binary.LittleEndian.Uint16(params[0:2])
	address := binary.LittleEndian.Uint16(params[2:4])
	if value != against {
		return cpu.SetRegister(register.Ip, address)
	}
	return nil
}

type Halt struct{ Passthrough }
