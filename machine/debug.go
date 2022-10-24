package machine

import (
	"fmt"
	"io"
	"strings"
	"the-machine/machine/memory"
	"the-machine/machine/register"
)

type MemoryType uint8

const (
	RAM MemoryType = 0
	ROM MemoryType = iota
)

type Resolution uint8

const (
	Byte Resolution = 0
	Uint Resolution = iota
)

type Representation uint8

const (
	Binary  Representation = 0
	Hex     Representation = iota
	Decimal Representation = iota
)

type RenderingDirection uint8

const (
	Horizontal RenderingDirection = 0
	Vertical   RenderingDirection = iota
)

type Debugger struct {
	vm        *Machine
	stream    io.Writer
	formatter Formatter
}

func NewDebugger(vm *Machine, f Formatter) Debugger {
	return Debugger{vm: vm, formatter: f}
}

func (x Debugger) CoreRegisters(number Representation) string {
	return x.registers(map[string]register.Register{
		"Ip": register.Ip,
		"Ac": register.Ac,
		"Sp": register.Sp,
		"Fp": register.Fp,
	}, number)
}

func (x Debugger) GeneralRegisters(number Representation) string {
	return x.registers(map[string]register.Register{
		"R1": register.R1,
		"R2": register.R2,
		"R3": register.R3,
		"R4": register.R4,
		"R5": register.R5,
		"R6": register.R6,
		"R7": register.R7,
		"R8": register.R8,
	}, number)
}

func (x Debugger) registers(registers map[string]register.Register, number Representation) string {
	_, valFormat := x.formatter.getFormat()
	positions := make([]string, len(registers))
	values := make([]string, len(registers))
	idx := 0
	for name, register := range registers {
		value := x.vm.cpu.GetRegister(register)
		values[idx] = fmt.Sprintf(valFormat, value)
		format := fmt.Sprintf("%%%ds", len(values[idx]))
		positions[idx] = fmt.Sprintf(format, name)
		idx++
	}
	return x.formatter.stitch(positions, values)
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

	return x.formatter.stitch(positions, values)
}

func (x Debugger) renderPosition(source memory.MemoryAccess, at memory.Address) (string, string) {
	posFormat, valFormat := x.formatter.getFormat()
	position := fmt.Sprintf(posFormat, at)
	var value string
	switch x.formatter.OutputAs {
	case Byte:
		if b, err := source.GetByte(at); err != nil {
			x.out(fmt.Sprintf("ERROR: unable to access byte at %v: %v", at, err))
			return position, fmt.Sprintf(strings.Repeat("", len(position)))
		} else {
			value = fmt.Sprintf(valFormat, b)
		}
	case Uint:
		if b, err := source.GetUint16(at); err != nil {
			x.out(fmt.Sprintf("ERROR: unable to access uint at %v: %v", at, err))
			return position, fmt.Sprintf(strings.Repeat("", len(position)))
		} else {
			value = fmt.Sprintf(valFormat, b)
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

type Formatter struct {
	Numbers   Representation
	OutputAs  Resolution
	Rendering RenderingDirection
}

func (x Formatter) getFormat() (string, string) {
	posFormat := "%4d"
	valFormat := "%#02x"
	switch x.Numbers {
	case Binary:
		switch x.OutputAs {
		case Byte:
			posFormat = "%10d"
			valFormat = "%#08b"
		case Uint:
			posFormat = "%18d"
			valFormat = "%#016b"
		}
	case Decimal:
		switch x.OutputAs {
		case Byte:
			posFormat = "%3d"
			valFormat = "%3d"
		case Uint:
			posFormat = "%5d"
			valFormat = "%05d"
		}

	}
	return posFormat, valFormat
}

func (x Formatter) stitch(first []string, rest ...[]string) string {
	switch x.Rendering {
	case Vertical:
		return x.stitchCols(first, rest...)
	case Horizontal:
		return x.stitchRows(first, rest...)
	default:
		return fmt.Sprintf("ERROR: unknown rendering direction: %d", x.Rendering)
	}
	return ""
}

func (x Formatter) stitchRows(first []string, rest ...[]string) string {
	out := make([]string, len(rest)+1)
	out[0] = strings.Join(first, " ")
	separator := strings.Repeat("-", len(out[0]))
	for idx, item := range rest {
		out[idx+1] = strings.Join(item, " ")
	}
	return strings.Join(out, fmt.Sprintf("\n%s\n", separator))
}

func (x Formatter) stitchCols(first []string, rest ...[]string) string {
	cols := make([]string, len(rest)+1)
	rows := make([]string, len(first))

	for rowIdx, item := range first {
		cols[0] = item
		ln := len(item)
		for colIdx, col := range rest {
			if rowIdx < len(col) {
				cols[colIdx+1] = col[rowIdx]
				ln = len(col[rowIdx])
			} else {
				cols[colIdx+1] = strings.Repeat(" ", ln)
			}
		}
		rows[rowIdx] = strings.Join(cols, " | ")
	}
	return strings.Join(rows, "\n")
}
