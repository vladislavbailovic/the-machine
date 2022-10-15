package machine

import (
	"encoding/binary"
	"fmt"
)

type address uint16
type memory []byte

func NewMemory(size int) *memory {
	mem := make(memory, size, size)
	return &mem
}

func (mem memory) GetByte(at address) (byte, error) {
	addr := int(at)
	if addr < len(mem) {
		return mem[addr], nil
	} else {
		return 0, fmt.Errorf("invalid memory access at %d", at)
	}
}

func (mem memory) GetUint16(at address) (uint16, error) {
	addr := int(at)
	if addr+1 <= len(mem) {
		hi, err := mem.GetByte(at)
		if err != nil {
			return 0, fmt.Errorf("uint16: error getting hi byte from %d: %v", at, err)
		}
		lo, err := mem.GetByte(at + 1)
		if err != nil {
			return 0, fmt.Errorf("uint16: error getting lo byte from %d: %v", at, err)
		}
		res := binary.LittleEndian.Uint16([]byte{hi, lo})
		return res, err
	} else {
		return 0, fmt.Errorf("invalid memory access at %d", at)
	}
}

func (mem *memory) SetByte(at address, value byte) error {
	if int(at) > cap(*mem) {
		return fmt.Errorf("invalid memory access at %d (of %d): trying to set byte 0x%02x", at, len(*mem), value)
	}
	(*mem)[at] = value
	return nil
}

func (mem *memory) SetUint16(at address, value uint16) error {
	if int(at)+1 > cap(*mem) {
		return fmt.Errorf("invalid memory access at %d (of %d): trying to set %d", at, len(*mem), value)
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
