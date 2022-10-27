package debug

import (
	"fmt"
	"strings"
	"the-machine/machine/cpu"
	"the-machine/machine/instruction"
	"the-machine/machine/memory"
	"the-machine/machine/register"
)

type Renderer struct {
	formatter Formatter
}

func NewRenderer(f Formatter) *Renderer {
	return &Renderer{formatter: f}
}

func (x *Renderer) SetFormatter(f Formatter) {
	x.formatter = f
}

func (x Renderer) GetFormatter() Formatter {
	return x.formatter
}

func (x Renderer) Registers(cpu *cpu.Cpu, registers []register.Register) string {
	_, valFormat := x.formatter.GetFormat()
	positions := make([]string, len(registers))
	values := make([]string, len(registers))
	for idx, register := range registers {
		name := register.Name()
		value := cpu.GetRegister(register)
		values[idx] = fmt.Sprintf(valFormat, value)
		format := fmt.Sprintf("%%%ds", len(values[idx]))
		positions[idx] = fmt.Sprintf(format, name)
	}
	return x.formatter.Stitch(positions, values)
}

func (x Renderer) Memory(source memory.MemoryAccess, startAt memory.Address, outputLen int) string {
	positions := make([]string, outputLen, outputLen)
	values := make([]string, outputLen, outputLen)
	for i := 0; i < outputLen; i++ {
		pos := int(startAt) + i
		positions[i], values[i] = x.memoryAt(source, memory.Address(pos))
	}

	return x.formatter.Stitch(positions, values)
}

func (x Renderer) memoryAt(source memory.MemoryAccess, at memory.Address) (string, string) {
	posFormat, valFormat := x.formatter.GetFormat()
	position := fmt.Sprintf(posFormat, at)
	var value string
	if x.formatter.OutputAs == Uint && uint16(at)%2 != 0 {
		value = strings.Repeat(" ", len(position))
	} else {
		switch x.formatter.OutputAs {
		case Byte:
			if b, err := source.GetByte(at); err != nil {
				x.Out(fmt.Sprintf("ERROR: unable to access byte at %v: %v", at, err))
				return position, fmt.Sprintf(strings.Repeat("", len(position)))
			} else {
				value = fmt.Sprintf(valFormat, b)
			}
		case Uint:
			if b, err := source.GetUint16(at); err != nil {
				x.Out(fmt.Sprintf("ERROR: unable to access uint at %v: %v", at, err))
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

func (x Renderer) decodeInstruction(instr uint16) (instruction.Instruction, error) {
	kind, raw := instruction.Decode(instr)
	if kind == instruction.HALT {
		return instruction.Descriptors[instruction.NOP], nil
	}

	decoded, ok := instruction.Descriptors[kind]
	if !ok {
		return instruction.Descriptors[instruction.NOP], fmt.Errorf("unknown instruction: %#02x", instr)
	}

	decoded.Raw = raw
	return decoded, nil
}

func (x Renderer) Disassembly(source memory.MemoryAccess, startAt memory.Address, outputLen int) string {
	x.formatter.OutputAs = Uint // Required for disassembly

	positions := make([]string, outputLen, outputLen)
	values := make([]string, outputLen, outputLen)
	instructions := make([]string, outputLen, outputLen)
	for i := 0; i < outputLen; i++ {
		pos := int(startAt) + i
		positions[i], values[i] = x.memoryAt(source, memory.Address(pos))

		instr := strings.Repeat(" ", len(positions[i]))
		if i%2 == 0 {
			if b, err := source.GetUint16(memory.Address(pos)); err != nil {
				x.Out(fmt.Sprintf("ERROR: unable to access uint at %v: %v", pos, err))
			} else {
				decoded, err := x.decodeInstruction(b)
				if err != nil {
					x.Out(fmt.Sprintf("ERROR: unable to disassemble at %d: %d: %v", pos, b, err))
				} else {
					instr = fmt.Sprintf("%v :: %v (%#010b)", decoded, decoded.Raw, decoded.Raw)
				}
			}
		}
		instructions[i] = instr
	}

	return x.formatter.Stitch(positions, values, instructions)
}

// TODO: figure this out
func (x Renderer) Out(msg string) {
	fmt.Println(msg)
}
