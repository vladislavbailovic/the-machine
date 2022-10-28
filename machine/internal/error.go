package internal

type MachineError struct {
	errStr string
	parent error
}

func Error(errStr string, parent error) error {
	return MachineError{errStr: errStr, parent: parent}
}

func (x MachineError) Error() string {
	return x.errStr
}

func (x MachineError) Unwrap() error {
	return x.parent
}
