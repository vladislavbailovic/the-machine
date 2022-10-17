package instruction

import "fmt"

type Parameter uint8

const (
	ParamRegister Parameter = 0
	ParamLiteral  Parameter = iota
	ParamAddress  Parameter = iota
)

func (x Parameter) GetBytesLength() int {
	switch x {
	case ParamRegister:
		return 1
	case ParamLiteral, ParamAddress:
		return 2
	}
	panic(fmt.Sprintf("unknown parameter type: %v", x))
}
