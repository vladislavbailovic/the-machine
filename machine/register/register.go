package register

import "the-machine/machine/memory"

type Register byte

const (
	Ip            Register = 0
	Ac            Register = iota
	R1            Register = iota
	R2            Register = iota
	R3            Register = iota
	R4            Register = iota
	_registerSize Register = iota
)

func Size() int {
	return int(_registerSize.AsAddress()) * 2
}

func (r Register) AsAddress() memory.Address {
	return memory.Address(r * 2)
}

func (r Register) AsByte() byte {
	return byte(r)
}
