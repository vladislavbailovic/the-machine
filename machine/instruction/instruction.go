package instruction

import (
	"fmt"
	"the-machine/machine/cpu"
	"the-machine/machine/internal"
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
		return internal.Error(fmt.Sprintf("error executing %v", x), err, internal.ErrorInstruction)
	}
	return nil
}

func (x Instruction) String() string {
	ex := x.Executor.String()
	if ex != "" {
		ex = " (" + ex + ") "
	} else {
		ex = ": "
	}
	return fmt.Sprintf("%s%s%d", x.Description, ex, x.Raw)
}
