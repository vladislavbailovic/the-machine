package instruction

import (
	"encoding/binary"
	"fmt"
	"the-machine/machine/cpu"
	"the-machine/machine/memory"
	"the-machine/machine/register"
)

type Executor interface {
	Execute([]byte, *cpu.Cpu, *memory.Memory) error
}

type Passthrough struct{}

func (x Passthrough) Execute(_ []byte, _ *cpu.Cpu, _ *memory.Memory) error {
	return nil
}

type Lit2Reg struct {
	Target register.Register
}

func (x Lit2Reg) Execute(params []byte, cpu *cpu.Cpu, mem *memory.Memory) error {
	if len(params) != 2 {
		return fmt.Errorf("LIT2REG[%v]: invalid parameter: %v", x.Target, params)
	}
	value := binary.LittleEndian.Uint16(params)
	return cpu.SetRegister(x.Target, value)
}

type Reg2Mem struct{}

func (x Reg2Mem) Execute(params []byte, cpu *cpu.Cpu, mem *memory.Memory) error {
	if len(params) != 3 {
		return fmt.Errorf("REG2MEM: invalid parameter: %v", params)
	}
	value, err := cpu.GetRegister(register.Register(params[0]))
	if err != nil {
		return fmt.Errorf("REG2MEM: error fetching from register %d (#2): %v", params[0], err)
	}
	address := memory.Address(binary.LittleEndian.Uint16(params[1:]))
	return mem.SetUint16(address, value)
}

type Lit2Mem struct{}

func (x Lit2Mem) Execute(params []byte, cpu *cpu.Cpu, mem *memory.Memory) error {
	if len(params) != 4 {
		return fmt.Errorf("LIT2MEM: invalid parameter: %v", params)
	}
	value := binary.LittleEndian.Uint16(params[:2])
	address := memory.Address(binary.LittleEndian.Uint16(params[2:]))
	return mem.SetUint16(address, value)
}

type OperateReg struct {
	Operation Op
}

func (x OperateReg) Execute(params []byte, cpu *cpu.Cpu, mem *memory.Memory) error {
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

func (x OperateRegLit) Execute(params []byte, cpu *cpu.Cpu, mem *memory.Memory) error {
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

type Jump struct {
	Comparison Comparison
}

func (x Jump) Execute(params []byte, cpu *cpu.Cpu, mem *memory.Memory) error {
	acu, err := cpu.GetRegister(register.Register(register.Ac))
	if err != nil {
		return fmt.Errorf("JMP[%d]: error fetching from Ac: %v", x.Comparison, err)
	}

	if len(params) != 4 {
		return fmt.Errorf("JMP[%d][%v]: invalid parameter: %v", x.Comparison, acu, params)
	}
	literal := binary.LittleEndian.Uint16(params[0:2])
	address := binary.LittleEndian.Uint16(params[2:4])

	writeIp := false
	switch x.Comparison {
	case CompNe:
		if literal != acu {
			writeIp = true
		}
	case CompEq:
		if literal == acu {
			writeIp = true
		}
	case CompGt:
		if acu > literal {
			writeIp = true
		}
	case CompGe:
		if acu >= literal {
			writeIp = true
		}
	case CompLt:
		if acu < literal {
			writeIp = true
		}
	case CompLe:
		if acu <= literal {
			writeIp = true
		}
	default:
		return fmt.Errorf("JMP[%d][%v]: invalid comparison: %v", x.Comparison, acu, params)
	}

	if writeIp {
		return cpu.SetRegister(register.Ip, address)
	}
	return nil
}

type Halt struct{ Passthrough }
