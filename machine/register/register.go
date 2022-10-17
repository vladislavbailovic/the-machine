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
	pos:  0x01,
}

var Ac = Register{
	name: "Accumulator",
	pos:  0x02,
}

var R1 = Register{
	name: "Register #1",
	pos:  0x11,
}

var R2 = Register{
	name: "Register #2",
	pos:  0x12,
}

var R3 = Register{
	name: "Register #3",
	pos:  0x13,
}

var R4 = Register{
	name: "Register #4",
	pos:  0x14,
}

func FromByte(b byte) (Register, error) {
	switch b {
	case 0x01:
		return Ip, nil
	case 0x02:
		return Ac, nil
	case 0x11:
		return R1, nil
	case 0x12:
		return R2, nil
	case 0x13:
		return R3, nil
	case 0x14:
		return R4, nil
	}
	return Register{}, fmt.Errorf("unknown register: %#02x", b)
}
