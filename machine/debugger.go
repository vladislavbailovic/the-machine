package machine

import (
	"fmt"
	"io"
	"strings"
	"the-machine/machine/debug"
	"the-machine/machine/memory"
	"the-machine/machine/register"
)

type MemoryType uint8

const (
	RAM MemoryType = 0
	ROM MemoryType = iota
)

type Debugger struct {
	vm        *Machine
	stream    io.Writer
	formatter debug.Formatter
}

func NewDebugger(vm *Machine, f debug.Formatter) *Debugger {
	return &Debugger{vm: vm, formatter: f}
}

func (x *Debugger) SetFormatter(f debug.Formatter) {
	x.formatter = f
}

func (x Debugger) CoreRegisters() string {
	return x.registers([]register.Register{
		register.Ip,
		register.Ac,
		register.Sp,
		register.Fp,
	})
}

func (x Debugger) GeneralRegisters() string {
	return x.registers([]register.Register{
		register.R1,
		register.R2,
		register.R3,
		register.R4,
		register.R5,
		register.R6,
		register.R7,
		register.R8,
	})
}

func (x Debugger) AllRegisters() string {
	return x.registers([]register.Register{
		register.Ip,
		register.Ac,
		register.Sp,
		register.Fp,
		register.R1,
		register.R2,
		register.R3,
		register.R4,
		register.R5,
		register.R6,
		register.R7,
		register.R8,
	})
}

func (x Debugger) registers(registers []register.Register) string {
	_, valFormat := x.formatter.GetFormat()
	positions := make([]string, len(registers))
	values := make([]string, len(registers))
	for idx, register := range registers {
		name := register.Name()
		value := x.vm.cpu.GetRegister(register)
		values[idx] = fmt.Sprintf(valFormat, value)
		format := fmt.Sprintf("%%%ds", len(values[idx]))
		positions[idx] = fmt.Sprintf(format, name)
	}
	return x.formatter.Stitch(positions, values)
}

func (x Debugger) Peek(startAt memory.Address, outputLen int, srcType MemoryType) string {
	var source memory.MemoryAccess
	switch srcType {
	case RAM:
		source = x.vm.ram
	case ROM:
		source = x.vm.rom
	default:
		x.out(fmt.Sprintf("ERROR: unknown source type: %v", srcType))
		return ""
	}

	positions := make([]string, outputLen, outputLen)
	values := make([]string, outputLen, outputLen)
	for i := 0; i < outputLen; i++ {
		pos := int(startAt) + i
		positions[i], values[i] = x.renderPosition(source, memory.Address(pos))
	}

	return x.formatter.Stitch(positions, values)
}

func (x Debugger) Disassemble(startAt memory.Address, outputLen int) string {
	x.formatter.OutputAs = debug.Uint // Required for disassembly
	source := x.vm.rom

	positions := make([]string, outputLen, outputLen)
	values := make([]string, outputLen, outputLen)
	instructions := make([]string, outputLen, outputLen)
	for i := 0; i < outputLen; i++ {
		pos := int(startAt) + i
		positions[i], values[i] = x.renderPosition(source, memory.Address(pos))

		instr := strings.Repeat(" ", len(positions[i]))
		if i%2 == 0 {
			if b, err := source.GetUint16(memory.Address(pos)); err != nil {
				x.out(fmt.Sprintf("ERROR: unable to access uint at %v: %v", pos, err))
			} else {
				decoded, err := x.vm.decode(b)
				if err != nil {
					x.out(fmt.Sprintf("ERROR: unable to disassemble at %d: %d: %v", pos, b, err))
				} else {
					instr = fmt.Sprintf("%v :: %v (%#010b)", decoded, decoded.Raw, decoded.Raw)
				}
			}
		}
		instructions[i] = instr
	}

	return x.formatter.Stitch(positions, values, instructions)
}

func (x Debugger) renderPosition(source memory.MemoryAccess, at memory.Address) (string, string) {
	posFormat, valFormat := x.formatter.GetFormat()
	position := fmt.Sprintf(posFormat, at)
	var value string
	if x.formatter.OutputAs == debug.Uint && uint16(at)%2 != 0 {
		value = strings.Repeat(" ", len(position))
	} else {
		switch x.formatter.OutputAs {
		case debug.Byte:
			if b, err := source.GetByte(at); err != nil {
				x.out(fmt.Sprintf("ERROR: unable to access byte at %v: %v", at, err))
				return position, fmt.Sprintf(strings.Repeat("", len(position)))
			} else {
				value = fmt.Sprintf(valFormat, b)
			}
		case debug.Uint:
			if b, err := source.GetUint16(at); err != nil {
				x.out(fmt.Sprintf("ERROR: unable to access uint at %v: %v", at, err))
				return position, fmt.Sprintf(strings.Repeat("", len(position)))
			} else {
				value = fmt.Sprintf(valFormat, b)
			}
		}
	}

	if len(value) > len(position) {
		diff := len(value) - len(position)
		position = strings.Repeat(" ", diff) + position
	}

	return position, value
}

func (x Debugger) out(msg string) {
	fmt.Println(msg)
}
