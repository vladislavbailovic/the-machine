package instruction

import (
	"encoding/binary"
	"fmt"
	"the-machine/machine/register"
)

type Executor interface {
	Execute([]byte) (Result, error)
}

type Passthrough struct{}

func (x Passthrough) Execute(p []byte) (Result, error) {
	return Result{}, nil
}

type Lit2Reg struct {
	Target register.Register
}

func (x Lit2Reg) Execute(params []byte) (Result, error) {
	if len(params) != 2 {
		return Result{}, fmt.Errorf("lit2reg[%v]: invalid parameter: %v", x.Target, params)
	}
	value := binary.LittleEndian.Uint16(params)
	return Result{Action: RecordRegister, Target: x.Target, Value: value}, nil
}

type AddTwo struct {
	V1 uint16
	V2 uint16
}

func (x AddTwo) Execute(_ []byte) (Result, error) {
	return Result{Action: RecordRegister, Target: register.Ac, Value: x.V1 + x.V2}, nil
}

type ExecError struct {
	Error error
}

func (x ExecError) Execute(_ []byte) (Result, error) {
	return Result{}, x.Error
}

type Jump struct {
	Against uint16
}

func (x Jump) Execute(params []byte) (Result, error) {
	if len(params) != 4 {
		return Result{}, fmt.Errorf("jne[%v]: invalid parameter: %v", x.Against, params)
	}
	value := binary.LittleEndian.Uint16(params[0:2])
	address := binary.LittleEndian.Uint16(params[2:4])
	if value != x.Against {
		return Result{Action: RecordRegister, Target: register.Ip, Value: address}, nil
	}
	return Result{}, nil
}

type Halt struct {
	End uint16
}

func (x Halt) Execute(_ []byte) (Result, error) {
	return Result{Action: RecordRegister, Target: register.Ip, Value: x.End}, nil
}
