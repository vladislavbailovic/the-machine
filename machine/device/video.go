package device

import (
	"fmt"
	"io"
	"math"
	"os"
	"the-machine/machine/memory"
)

const (
	screenWidth  uint8 = 255
	screenHeight uint8 = 255
)

const (
	termEsc string = "\u001B"
)

type Video struct {
	stream io.Writer
}

func NewVideo() memory.MemoryAccess {
	return Video{stream: os.Stdout}
}

func (x Video) GetByte(at memory.Address) (byte, error)     { return 0, nil }
func (x Video) GetUint16(at memory.Address) (uint16, error) { return 0, nil }

func (x Video) SetUint16(at memory.Address, val uint16) error {
	return x.SetByte(at, byte(val))
}

func (x Video) SetByte(at memory.Address, val byte) error {
	coords, err := x.addressToCoords(at)
	if err != nil {
		return fmt.Errorf("unable to print output %c at %v: %v", val, at, err)
	}
	fmt.Fprintf(x.stream,
		fmt.Sprintf("%s[%d;%dH%c", termEsc, coords[1], coords[0], val))
	return nil
}

func (v Video) addressToCoords(at memory.Address) ([]uint8, error) {
	coords := make([]uint8, 2, 2)

	x := uint16(at) % uint16(screenWidth)
	if x > uint16(screenWidth) {
		return coords, fmt.Errorf("X outside bounds (%d): %d", screenWidth, x)
	}
	coords[0] = uint8(x)

	y := uint16(math.Floor(float64(at) / float64(screenWidth)))
	if y > uint16(screenWidth) {
		return coords, fmt.Errorf("X outside bounds (%d): %d", screenWidth, y)
	}
	coords[1] = uint8(y)

	return coords, nil
}
