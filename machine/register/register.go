package register

import "fmt"

type Register struct {
	name string
	pos  byte
}

func (r Register) AsByte() byte {
	return r.pos
}

var Ip = Register{
	name: "Instruction Pointer",
	pos:  15,
}

var Ac = Register{
	name: "Accumulator",
	pos:  14,
}

var Sp = Register{
	name: "Stack Pointer",
	pos:  13,
}

var Fp = Register{
	name: "Frame Pointer",
	pos:  12,
}

var R1 = Register{
	name: "Register #1",
	pos:  0,
}

var R2 = Register{
	name: "Register #2",
	pos:  1,
}

var R3 = Register{
	name: "Register #3",
	pos:  2,
}

var R4 = Register{
	name: "Register #4",
	pos:  3,
}

var R5 = Register{
	name: "Register #5",
	pos:  4,
}

var R6 = Register{
	name: "Register #6",
	pos:  5,
}

var R7 = Register{
	name: "Register #7",
	pos:  6,
}

var R8 = Register{
	name: "Register #8",
	pos:  7,
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
	return Register{}, fmt.Errorf("unknown register: %#02x", b)
}
