package instruction

import (
	"fmt"
	"the-machine/machine/cpu"
	"the-machine/machine/memory"
	"the-machine/machine/register"
)

type Executor interface {
	Execute(uint16, *cpu.Cpu, memory.MemoryAccess) error
}

type Passthrough struct{}

func (x Passthrough) Execute(_ uint16, _ *cpu.Cpu, _ memory.MemoryAccess) error {
	return nil
}

type Lit2Reg struct {
	Target register.Register
}

func (x Lit2Reg) Execute(value uint16, cpu *cpu.Cpu, mem memory.MemoryAccess) error {
	cpu.SetRegister(x.Target, value)
	return nil
}

type Reg2Reg struct{ unpacker }

func (x Reg2Reg) Execute(raw uint16, cpu *cpu.Cpu, mem memory.MemoryAccess) error {
	params := x.unpack(raw)

	source, err := register.FromByte(params[0])
	if err != nil {
		return fmt.Errorf("REG2REG: invalid source register (%#02x): %w", params[0], err)
	}
	value := cpu.GetRegister(source)

	destination, err := register.FromByte(params[1])
	if err != nil {
		return fmt.Errorf("REG2REG: invalid destination register (%#02x): %w", params[1], err)
	}

	cpu.SetRegister(destination, value)
	return nil
}

type Reg2Stack struct{}

func (x Reg2Stack) Execute(raw uint16, cpu *cpu.Cpu, mem memory.MemoryAccess) error {
	source, err := register.FromByte(byte(raw))
	if err != nil {
		return fmt.Errorf("REG2STACK: invalid source register (%#02x): %w", raw, err)
	}
	value := cpu.GetRegister(source)

	return cpu.Push(value)
}

type Lit2Stack struct{}

func (x Lit2Stack) Execute(value uint16, cpu *cpu.Cpu, mem memory.MemoryAccess) error {
	return cpu.Push(value)
}

type Stack2Reg struct{}

func (x Stack2Reg) Execute(raw uint16, cpu *cpu.Cpu, mem memory.MemoryAccess) error {
	destination, err := register.FromByte(byte(raw))
	if err != nil {
		return fmt.Errorf("STACK2REG: invalid source register (%#02x): %w", raw, err)
	}

	value, err := cpu.Pop()
	if err != nil {
		return fmt.Errorf("STACK2REG: stack underflow: %w", err)
	}

	cpu.SetRegister(destination, value)
	return nil
}

type Ac2Reg struct{}

func (x Ac2Reg) Execute(params uint16, cpu *cpu.Cpu, mem memory.MemoryAccess) error {
	value := cpu.GetRegister(register.Ac)

	destination, err := register.FromByte(byte(params))
	if err != nil {
		return fmt.Errorf("AC2REG: invalid destination register (%#02x): %w", params, err)
	}

	cpu.SetRegister(destination, value)
	return nil
}

type Reg2Mem struct{}

func (x Reg2Mem) Execute(params uint16, cpu *cpu.Cpu, mem memory.MemoryAccess) error {
	r1, err := register.FromByte(byte(params))
	if err != nil {
		return fmt.Errorf("REG2MEM: invalid register (%#02x): %w", params, err)
	}
	value := cpu.GetRegister(r1)

	address := memory.Address(cpu.GetRegister(register.Ac))
	return mem.SetUint16(address, value)
}

type Lit2Mem struct{}

func (x Lit2Mem) Execute(value uint16, cpu *cpu.Cpu, mem memory.MemoryAccess) error {
	address := memory.Address(cpu.GetRegister(register.Ac))
	return mem.SetUint16(address, value)
}

type Mem2Reg struct{ unpacker }

func (x Mem2Reg) Execute(raw uint16, cpu *cpu.Cpu, mem memory.MemoryAccess) error {
	params := x.unpack(raw)

	source, err := register.FromByte(params[0])
	if err != nil {
		return fmt.Errorf("MEM2REG: invalid address register (%#02x): %w", params[0], err)
	}
	address := cpu.GetRegister(source)

	destination, err := register.FromByte(params[1])
	if err != nil {
		return fmt.Errorf("MEM2REG: invalid destination register (%#02x): %w", params[1], err)
	}

	value, err := mem.GetUint16(memory.Address(address))
	if err != nil {
		return fmt.Errorf("MEM2REG: error accessing memory at %d (%#02x): %w", address, address, err)
	}

	cpu.SetRegister(destination, value)
	return nil
}

type OperateReg struct {
	unpacker
	Operation Op
}

func (x OperateReg) Execute(raw uint16, cpu *cpu.Cpu, mem memory.MemoryAccess) error {
	params := x.unpack(raw)
	r1, err := register.FromByte(params[0])
	if err != nil {
		return fmt.Errorf("OP_REG %d: invalid register #1 (%#02x): %w", x.Operation, params[0], err)
	}
	v1 := cpu.GetRegister(r1)
	r2, err := register.FromByte(params[1])
	if err != nil {
		return fmt.Errorf("OP_REG %d: invalid register #2 (%#02x): %w", x.Operation, params[1], err)
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
	case OpAnd:
		cpu.SetRegister(register.Ac, v1&v2)
		return nil
	case OpOr:
		cpu.SetRegister(register.Ac, v1|v2)
		return nil
	case OpXor:
		cpu.SetRegister(register.Ac, v1^v2)
		return nil
	default:
		return fmt.Errorf("OP_REG %d: unknown operation", x.Operation)
	}
}

type OperateRegLit struct {
	unpacker
	Operation Op
}

func (x OperateRegLit) Execute(raw uint16, cpu *cpu.Cpu, mem memory.MemoryAccess) error {
	params := x.unpack(raw)
	r1, err := register.FromByte(params[0])
	if err != nil {
		return fmt.Errorf("OP_REG_LIT %d: invalid register (%#02x): %w", x.Operation, params[0], err)
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
	case OpShl:
		cpu.SetRegister(register.Ac, reg<<literal)
		return nil
	case OpShr:
		cpu.SetRegister(register.Ac, reg>>literal)
		return nil
	case OpAnd:
		cpu.SetRegister(register.Ac, reg&literal)
		return nil
	case OpOr:
		cpu.SetRegister(register.Ac, reg|literal)
		return nil
	case OpXor:
		cpu.SetRegister(register.Ac, reg^literal)
		return nil
	default:
		return fmt.Errorf("OP_REG_LIT %d: unknown operation", x.Operation)
	}
}

type OperateStack struct {
	Operation Op
}

func (x OperateStack) Execute(raw uint16, cpu *cpu.Cpu, mem memory.MemoryAccess) error {
	operand1, err := cpu.Pop()
	if err != nil {
		return fmt.Errorf("OP_STACK %d: stack underflow getting first operand: %w", x.Operation, err)
	}
	operand2, err := cpu.Pop()
	if err != nil {
		return fmt.Errorf("OP_STACK %d: stack underflow getting second operand: %w", x.Operation, err)
	}

	switch x.Operation {
	case OpAdd:
		cpu.Push(operand1 + operand2)
		return nil
	case OpSub:
		cpu.Push(operand1 - operand2)
		return nil
	case OpMul:
		cpu.Push(operand1 * operand2)
		return nil
	case OpDiv:
		cpu.Push(operand1 / operand2)
		return nil
	case OpMod:
		cpu.Push(operand1 % operand2)
		return nil
	default:
		return fmt.Errorf("OP_STACK %d: unknown operation", x.Operation)
	}
}

type Jump struct {
	unpacker
	Comparison Comparison
}

func (x Jump) Execute(raw uint16, cpu *cpu.Cpu, mem memory.MemoryAccess) error {
	acu := cpu.GetRegister(register.Register(register.Ac))

	params := x.unpack(raw)

	cr, err := register.FromByte(params[0])
	if err != nil {
		return fmt.Errorf("JMP[%d][%v]: invalid comparison register (%#02x): %w", x.Comparison, cr, params[0], err)
	}
	// fmt.Printf("\t- Comparison register: %v (%016b)\n", cr, params[0])
	compareWith := cpu.GetRegister(cr)

	ar, err := register.FromByte(params[1])
	if err != nil {
		return fmt.Errorf("JMP[%d][%v]: invalid address register (%#02x): %w", x.Comparison, acu, params[1], err)
	}
	// fmt.Printf("\t- Address register: %v (%016b)\n", ar, params[1])
	address := cpu.GetRegister(ar)

	// fmt.Printf("\tComparing %d from %v (%d) %d from acu\n", compareWith, cr, x.Comparison, acu)

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
		return fmt.Errorf("JMP[%d][%v]: invalid comparison: %w", x.Comparison, acu, params)
	}

	if writeIp {
		// fmt.Printf("----- Gonna JUMP to %d ----\n", address)
		cpu.SetRegister(register.Ip, address)
	}
	return nil
}

type Halt struct{ Passthrough }

type Call struct{}

func (x Call) Execute(raw uint16, cpu *cpu.Cpu, _ memory.MemoryAccess) error {
	reg, err := register.FromByte(byte(raw))
	if err != nil {
		return fmt.Errorf("CALL: unknown register %d: %w", raw, err)
	}
	address := cpu.GetRegister(reg)

	if err := cpu.StoreFrame(); err != nil {
		return fmt.Errorf("CALL: error storing frame before calling %d: %w", address, err)
	}

	cpu.SetRegister(register.Ip, address)

	return nil
}

type Return struct{}

func (x Return) Execute(_ uint16, cpu *cpu.Cpu, _ memory.MemoryAccess) error {
	if err := cpu.RestoreFrame(); err != nil {
		return fmt.Errorf("RET: error restoring frame: %w", err)
	}
	return nil
}
