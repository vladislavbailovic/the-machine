package instruction

import "the-machine/machine/register"

type Action byte

const (
	Nop            Action = 0
	RecordRegister Action = iota
)

type Result struct {
	Action Action
	Value  uint16
	Target register.Register
}
