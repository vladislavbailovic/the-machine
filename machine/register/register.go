package register

type Register byte

const (
	Ip            Register = 0
	Ac            Register = iota
	R1            Register = iota
	R2            Register = iota
	R3            Register = iota
	R4            Register = iota
	_registerSize          = iota
)

func Size() byte {
	return _registerSize
}

func (r Register) Address() uint16 {
	return uint16(r * 2)

}
