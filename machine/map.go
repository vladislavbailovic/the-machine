package machine

import (
	"the-machine/machine/device"
	"the-machine/machine/memory"
)

type MemoryMap map[memory.MemoryType]memory.MemoryAccess

func NewMemoryMap(ramSize int, romSize int) MemoryMap {
	return map[memory.MemoryType]memory.MemoryAccess{
		memory.RAM:       memory.NewMemory(ramSize),
		memory.ROM:       memory.NewMemory(romSize),
		memory.DeviceVGA: device.NewVideo(),
		memory.DeviceIO:  device.NewIoMap(),
	}
}
