package instruction

import (
	"encoding/binary"
	"fmt"
	"the-machine/machine/cpu"
	"the-machine/machine/memory"
	"the-machine/machine/register"
)

// TODO: move this somewhere
func unpack2(params []byte) []byte {
	raw := binary.LittleEndian.Uint16(params)

	b1 := byte(
		(raw & 0b0000_0000_1111_0000) >> 4,
	)
	b2 := byte(
		raw & 0b0000_0000_0000_1111,
	)
	// fmt.Printf("\t- raw: %016b (%d)\n", raw, raw)
	// fmt.Printf("\t-  b1: %016b (%d)\n", b1, b1)
	// fmt.Printf("\t-  b2: %016b (%d)\n", b2, b2)

	return []byte{b1, b2}
}

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

type Reg2Reg struct{}

func (x Reg2Reg) Execute(params []byte, cpu *cpu.Cpu, mem *memory.Memory) error {
	params = unpack2(params)

	source, err := register.FromByte(params[0])
	if err != nil {
		return fmt.Errorf("REG2REG: invalid source register (%#02x): %v", params[0], err)
	}
	value := cpu.GetRegister(source)

	destination, err := register.FromByte(params[1])
	if err != nil {
		return fmt.Errorf("REG2REG: invalid destination register (%#02x): %v", params[1], err)
	}

	cpu.SetRegister(destination, value)
	return nil
}

type Ac2Reg struct{}

func (x Ac2Reg) Execute(params []byte, cpu *cpu.Cpu, mem *memory.Memory) error {
	value := cpu.GetRegister(register.Ac)

	destination, err := register.FromByte(params[0])
	if err != nil {
		return fmt.Errorf("AC2REG: invalid destination register (%#02x): %v", params[0], err)
	}

	cpu.SetRegister(destination, value)
	return nil
}

type Reg2Mem struct{}

func (x Reg2Mem) Execute(params []byte, cpu *cpu.Cpu, mem *memory.Memory) error {
	r1, err := register.FromByte(params[0])
	if err != nil {
		return fmt.Errorf("REG2MEM: invalid register (%#02x): %v", params[0], err)
	}
	value := cpu.GetRegister(r1)

	address := memory.Address(cpu.GetRegister(register.Ac))
	return mem.SetUint16(address, value)
}

type Lit2Mem struct{}

func (x Lit2Mem) Execute(params []byte, cpu *cpu.Cpu, mem *memory.Memory) error {
	value := binary.LittleEndian.Uint16(params)
	address := memory.Address(cpu.GetRegister(register.Ac))
	return mem.SetUint16(address, value)
}

type OperateReg struct {
	Operation Op
}

func (x OperateReg) Execute(params []byte, cpu *cpu.Cpu, mem *memory.Memory) error {
	if len(params) != 2 {
		return fmt.Errorf("OP_REG %d: invalid params: %v", x.Operation, params)
	}
	params = unpack2(params)
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
	if len(params) != 2 {
		return fmt.Errorf("OP_REG_LIT %d: invalid params: %v", x.Operation, params)
	}
	params = unpack2(params)
	r1, err := register.FromByte(params[0])
	if err != nil {
		return fmt.Errorf("OP_REG_LIT %d: invalid register (%#02x): %v", x.Operation, params[0], err)
	}
	reg := cpu.GetRegister(r1)
	literal := uint16(params[1])
	// fmt.Printf("Got register: %d - %016b (from %016b)\n", params[0], params[0], params[0])
	// fmt.Printf("Got literal: %d - %016b (from %016b)\n", literal, literal, params[1])

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

	if len(params) != 2 {
		return fmt.Errorf("JMP[%d][%v]: invalid parameter: %v", x.Comparison, acu, params)
	}
	params = unpack2(params)

	cr, err := register.FromByte(params[0])
	if err != nil {
		return fmt.Errorf("JMP[%d][%v]: invalid comparison register (%#02x): %v", x.Comparison, cr, params[0], err)
	}
	fmt.Printf("\t- Comparison register: %v (%016b)\n", cr, params[0])
	compareWith := cpu.GetRegister(cr)

	ar, err := register.FromByte(params[1])
	if err != nil {
		return fmt.Errorf("JMP[%d][%v]: invalid address register (%#02x): %v", x.Comparison, acu, params[1], err)
	}
	fmt.Printf("\t- Address register: %v (%016b)\n", ar, params[1])
	address := cpu.GetRegister(ar)

	fmt.Printf("\tComparing %d from %v (%d) %d from acu\n", compareWith, cr, x.Comparison, acu)

	writeIp := false
	switch x.Comparison {
	case CompNe:
		if compareWith != acu {
			writeIp = true
		}
	case CompEq:
		if compareWith == acu {
			writeIp = true
		}
	case CompGt:
		if acu > compareWith {
			writeIp = true
		}
	case CompGe:
		if acu >= compareWith {
			writeIp = true
		}
	case CompLt:
		if acu < compareWith {
			writeIp = true
		}
	case CompLe:
		if acu <= compareWith {
			writeIp = true
		}
	default:
		return fmt.Errorf("JMP[%d][%v]: invalid comparison: %v", x.Comparison, acu, params)
	}

	if writeIp {
		fmt.Printf("----- Gonna JUMP to %d ----\n", address)
		cpu.SetRegister(register.Ip, address)
	}
	return nil
}

type Halt struct{ Passthrough }
