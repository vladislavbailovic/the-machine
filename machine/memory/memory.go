package memory

import (
	"encoding/binary"
	"fmt"
	"the-machine/machine/internal"
)

type MemoryAccess interface {
	GetByte(Address) (byte, error)
	GetUint16(Address) (uint16, error)
	SetByte(Address, byte) error
	SetUint16(Address, uint16) error
}

type Address uint16
type Memory []byte

func NewMemory(size int) MemoryAccess {
	mem := make(Memory, size, size)
	return &mem
}

func (mem Memory) GetByte(at Address) (byte, error) {
	addr := int(at)
	if addr < len(mem) {
		return mem[addr], nil
	} else {
		return 0, internal.Error(fmt.Sprintf("invalid memory access at %d", at), nil, internal.ErrorMemory)
	}
}

func (mem Memory) GetUint16(at Address) (uint16, error) {
	addr := int(at)
	if addr+1 <= len(mem) {
		hi, err := mem.GetByte(at)
		if err != nil {
			return 0, internal.Error(fmt.Sprintf("uint16: error getting hi byte from %d", at), err, internal.ErrorMemory)
		}
		lo, err := mem.GetByte(at + 1)
		if err != nil {
			return 0, internal.Error(fmt.Sprintf("uint16: error getting lo byte from %d", at), err, internal.ErrorMemory)
		}
		res := binary.LittleEndian.Uint16([]byte{hi, lo})
		return res, err
	} else {
		return 0, internal.Error(fmt.Sprintf("invalid memory access at %d", at), nil, internal.ErrorMemory)
	}
}

func (mem *Memory) SetByte(at Address, value byte) error {
	if int(at) > cap(*mem) {
		return internal.Error(fmt.Sprintf("invalid memory access at %d (of %d): trying to set byte %#02x", at, len(*mem), value), nil, internal.ErrorMemory)
	}
	(*mem)[at] = value
	return nil
}

func (mem *Memory) SetUint16(at Address, value uint16) error {
	if int(at)+1 > cap(*mem) {
		return internal.Error(fmt.Sprintf("invalid memory access at %d (of %d): trying to set %d", at, len(*mem), value), nil, internal.ErrorMemory)
	}
	b := make([]byte, 2)
	binary.LittleEndian.PutUint16(b, value)
	if err := mem.SetByte(at, b[0]); err != nil { // hi
		return err
	}
	if err := mem.SetByte(at+1, b[1]); err != nil { // lo
		return err
	}
	return nil
}
