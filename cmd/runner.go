package cmd

import (
	"fmt"
	"the-machine/machine"
	"the-machine/machine/debug"
)

func Run(vm machine.Machine) (int, error) {
	step := 0
	for step < 0xffff {
		if err := vm.Tick(); err != nil {
			terr := fmt.Errorf("error at tick %d: %w", step, err)
			vm.DebugError(terr)
			return step, terr
		}
		step++
		if vm.IsDone() {
			break
		}
	}
	return step, nil
}

func RunFile(fname string) {
	vm := machine.NewMachine(2048)
	loader := debug.NewAsciiLoader(fname, debug.Decimal)
	if program, err := loader.Load(); err != nil {
		vm.DebugError(err)
		return
	} else {
		vm.LoadProgram(0, program)
	}
	if _, err := Run(vm); err != nil {
		vm.DebugError(err)
		return
	}
}
