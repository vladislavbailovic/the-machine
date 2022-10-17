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
	cpu.SetRegister(x.Target, value)
	return nil
}

type Reg2Mem struct{}

func (x Reg2Mem) Execute(params []byte, cpu *cpu.Cpu, mem *memory.Memory) error {
	if len(params) != 3 {
		return fmt.Errorf("REG2MEM: invalid parameter: %v", params)
	}
	r1, err := register.FromByte(params[0])
	if err != nil {
		return fmt.Errorf("REG2MEM: invalid register (%#02x): %v", params[0], err)
	}
	value := cpu.GetRegister(r1)
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
	r1, err := register.FromByte(params[0])
	if err != nil {
		return fmt.Errorf("OP_REG %d: invalid register #1 (%#02x): %v", x.Operation, params[0], err)
	}
	v1 := cpu.GetRegister(r1)
	r2, err := register.FromByte(params[1])
	if err != nil {
		return fmt.Errorf("OP_REG %d: invalid register #2 (%#02x): %v", x.Operation, params[1], err)
	}
	v2 := cpu.GetRegister(r2)

	switch x.Operation {
	case OpAdd:
		cpu.SetRegister(register.Ac, v1+v2)
		return nil
	case OpSub:
		cpu.SetRegister(register.Ac, v1-v2)
		return nil
	case OpMul:
		cpu.SetRegister(register.Ac, v1*v2)
		return nil
	case OpDiv:
		cpu.SetRegister(register.Ac, v1/v2)
		return nil
	case OpMod:
		cpu.SetRegister(register.Ac, v1%v2)
		return nil
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
	r1, err := register.FromByte(params[0])
	if err != nil {
		return fmt.Errorf("OP_REG_LIT %d: invalid register (%#02x): %v", x.Operation, params[0], err)
	}
	reg := cpu.GetRegister(r1)
	literal := binary.LittleEndian.Uint16(params[1:])

	switch x.Operation {
	case OpAdd:
		cpu.SetRegister(register.Ac, reg+literal)
		return nil
	case OpSub:
		cpu.SetRegister(register.Ac, reg-literal)
		return nil
	case OpMul:
		cpu.SetRegister(register.Ac, reg*literal)
		return nil
	case OpDiv:
		cpu.SetRegister(register.Ac, reg/literal)
		return nil
	case OpMod:
		cpu.SetRegister(register.Ac, reg%literal)
		return nil
	default:
		return fmt.Errorf("OP_REG_LIT %d: unknown operation", x.Operation)
	}
}

type Jump struct {
	Comparison Comparison
}

func (x Jump) Execute(params []byte, cpu *cpu.Cpu, mem *memory.Memory) error {
	acu := cpu.GetRegister(register.Register(register.Ac))

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
		cpu.SetRegister(register.Ip, address)
	}
	return nil
}

type Halt struct{ Passthrough }
