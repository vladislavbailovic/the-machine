package memory

import "fmt"

type MemoryType uint8

const (
	RAM       MemoryType = 0
	ROM       MemoryType = iota
	DeviceVGA MemoryType = iota
	DeviceIO  MemoryType = iota
)

func (x MemoryType) String() string {
	switch x {
	case RAM:
		return "RAM"
	case ROM:
		return "ROM"
	case DeviceVGA:
		return "VGA"
	case DeviceIO:
		return "IO"
	default:
		return fmt.Sprintf("unknown memory type: %d", x)
	}
}
