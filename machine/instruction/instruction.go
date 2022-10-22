package instruction

import (
	"encoding/binary"
	"fmt"
	"the-machine/machine/cpu"
	"the-machine/machine/memory"
)

type Instruction struct {
	Description string
	Raw         uint16
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

// TODO: move param parsing to machine::decode
func (x Instruction) getParams(cpu *cpu.Cpu, mem *memory.Memory) ([]byte, error) {
	b := make([]byte, 2)
	binary.LittleEndian.PutUint16(b, x.Raw)
	return b, nil
}
