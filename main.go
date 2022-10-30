package main

import (
	"os"
	"the-machine/cmd"
	"the-machine/machine"
)

func main() {
	if len(os.Args) > 1 {
		fname := os.Args[1]
		// TODO: validate fname
		cmd.RunFile(fname)
	} else {
		main_InteractiveDebugger()
	}
}

func main_InteractiveDebugger() {
	vm := machine.NewMachine(0xffff)
	vm.Debug()
}
