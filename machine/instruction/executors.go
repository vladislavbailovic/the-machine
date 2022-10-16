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

type OperateReg struct {
	Operation Op
}

func (x OperateReg) Execute(params []byte, cpu *cpu.Cpu) error {
	if len(params) != 2 {
		return fmt.Errorf("OP_REG %d: invalid params: %v", x.Operation, params)
	}
	v1, err := cpu.GetRegister(register.Register(params[0]))
	if err != nil {
		return fmt.Errorf("OP_REG %d: error fetching from register %d (#1): %v", x.Operation, params[0], err)
	}
	v2, err := cpu.GetRegister(register.Register(params[1]))
	if err != nil {
		return fmt.Errorf("OP_REG %d: error fetching from register %d (#2): %v", x.Operation, params[1], err)
	}

	switch x.Operation {
	case OpAdd:
		return cpu.SetRegister(register.Ac, v1+v2)
	case OpSub:
		return cpu.SetRegister(register.Ac, v1-v2)
	case OpMul:
		return cpu.SetRegister(register.Ac, v1*v2)
	case OpDiv:
		return cpu.SetRegister(register.Ac, v1/v2)
	case OpMod:
		return cpu.SetRegister(register.Ac, v1%v2)
	default:
		return fmt.Errorf("OP_REG %d: unknown operation", x.Operation)
	}
}

type OperateRegLit struct {
	Operation Op
}

func (x OperateRegLit) Execute(params []byte, cpu *cpu.Cpu) error {
	if len(params) != 3 {
		return fmt.Errorf("OP_REG_LIT %d: invalid params: %v", x.Operation, params)
	}
	reg, err := cpu.GetRegister(register.Register(params[0]))
	if err != nil {
		return fmt.Errorf("OP_REG_LIT %d: error fetching from register %d (#2): %v", x.Operation, params[0], err)
	}
	literal := binary.LittleEndian.Uint16(params[1:])

	switch x.Operation {
	case OpAdd:
		return cpu.SetRegister(register.Ac, reg+literal)
	case OpSub:
		return cpu.SetRegister(register.Ac, reg-literal)
	case OpMul:
		return cpu.SetRegister(register.Ac, reg*literal)
	case OpDiv:
		return cpu.SetRegister(register.Ac, reg/literal)
	case OpMod:
		return cpu.SetRegister(register.Ac, reg%literal)
	default:
		return fmt.Errorf("OP_REG_LIT %d: unknown operation", x.Operation)
	}
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
