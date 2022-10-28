package register

import (
	"fmt"
	"the-machine/machine/internal"
)

type Register struct {
	description string
	name        string
	pos         byte
}

func (x Register) AsByte() byte {
	return x.pos
}

func (x Register) AsUint16() uint16 {
	return uint16(x.pos)
}

func (x Register) Name() string {
	return x.name
}

var Ip = Register{
	description: "Instruction Pointer",
	name:        "Ip",
	pos:         15,
}

var Ac = Register{
	description: "Accumulator",
	name:        "Ac",
	pos:         14,
}

var Sp = Register{
	description: "Stack Pointer",
	name:        "Sp",
	pos:         13,
}

var Fp = Register{
	description: "Frame Pointer",
	name:        "Fp",
	pos:         12,
}

var R1 = Register{
	description: "Register #1",
	name:        "R1",
	pos:         0,
}

var R2 = Register{
	description: "Register #2",
	name:        "R2",
	pos:         1,
}

var R3 = Register{
	description: "Register #3",
	name:        "R3",
	pos:         2,
}

var R4 = Register{
	description: "Register #4",
	name:        "R4",
	pos:         3,
}

var R5 = Register{
	description: "Register #5",
	name:        "R5",
	pos:         4,
}

var R6 = Register{
	description: "Register #6",
	name:        "R6",
	pos:         5,
}

var R7 = Register{
	description: "Register #7",
	name:        "R7",
	pos:         6,
}

var R8 = Register{
	description: "Register #8",
	name:        "R8",
	pos:         7,
}

func FromByte(b byte) (Register, error) {
	switch b {
	case Ip.pos:
		return Ip, nil
	case Sp.pos:
		return Sp, nil
	case Fp.pos:
		return Fp, nil
	case Ac.pos:
		return Ac, nil
	case R1.pos:
		return R1, nil
	case R2.pos:
		return R2, nil
	case R3.pos:
		return R3, nil
	case R4.pos:
		return R4, nil
	case R5.pos:
		return R5, nil
	case R6.pos:
		return R6, nil
	case R7.pos:
		return R7, nil
	case R8.pos:
		return R8, nil
	}
	return Register{}, internal.Error(fmt.Sprintf("unknown register: %#02x", b), nil)
}
