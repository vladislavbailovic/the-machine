package instruction

import (
	"fmt"
	"the-machine/machine/cpu"
	"the-machine/machine/memory"
)

type Instruction struct {
	Description string
	Raw         uint16
	Executor    Executor
}

func (x Instruction) Execute(cpu *cpu.Cpu, memory memory.MemoryAccess) error {
	err := x.Executor.Execute(x.Raw, cpu, memory)
	if err != nil {
		return fmt.Errorf("error executing \"%s\": %v", x.Description, err)
	}
	return nil
}
