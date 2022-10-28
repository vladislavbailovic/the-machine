package internal

import "fmt"

type MachineErrorSource string

const (
	ErrorReg2Reg   MachineErrorSource = "Reg2Reg"
	ErrorReg2Stack MachineErrorSource = "Reg2Stack"
	ErrorStack2Reg MachineErrorSource = "Stack2Reg"
	ErrorAc2Reg    MachineErrorSource = "Ac2Reg"
	ErrorReg2Mem   MachineErrorSource = "Reg2Mem"
	ErrorMem2Reg   MachineErrorSource = "Mem2Reg"
	ErrorOpReg     MachineErrorSource = "OpReg"
	ErrorOpRegLit  MachineErrorSource = "OpRegLit"
	ErrorOpStack   MachineErrorSource = "OpStack"
	ErrorJmp       MachineErrorSource = "Jmp"
	ErrorCall      MachineErrorSource = "Call"
	ErrorRet       MachineErrorSource = "Ret"

	ErrorMemory      MachineErrorSource = "Memory"
	ErrorCpu         MachineErrorSource = "Cpu"
	ErrorInstruction MachineErrorSource = "Instruction"
	ErrorInterface   MachineErrorSource = "Interface"
	ErrorDebugger    MachineErrorSource = "Debugger"

	ErrorRuntime MachineErrorSource = "Runtime"
	ErrorLoading MachineErrorSource = "Loading"
)

type MachineError struct {
	errStr string
	parent error
	source MachineErrorSource
}

func Error(errStr string, parent error, source MachineErrorSource) error {
	return MachineError{errStr: errStr, parent: parent, source: source}
}

func (x MachineError) Error() string {
	return fmt.Sprintf("[%s] %s", x.source, x.errStr)
}

func (x MachineError) Unwrap() error {
	return x.parent
}
