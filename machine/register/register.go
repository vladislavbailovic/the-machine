package register

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
	return int(_registerSize)
}

func (r Register) AsByte() byte {
	return byte(r)
}
